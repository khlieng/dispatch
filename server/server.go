package server

import (
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/khlieng/dispatch/config"
	"github.com/khlieng/dispatch/pkg/letsencrypt"
	"github.com/khlieng/dispatch/pkg/session"
	"github.com/khlieng/dispatch/storage"
)

var channelStore = storage.NewChannelStore()

type Dispatch struct {
	Store        storage.Store
	SessionStore storage.SessionStore

	GetMessageStore          func(*storage.User) (storage.MessageStore, error)
	GetMessageSearchProvider func(*storage.User) (storage.MessageSearchProvider, error)

	cfg      *config.Config
	upgrader websocket.Upgrader
	states   *stateStore
	lock     sync.Mutex
}

func New(cfg *config.Config) *Dispatch {
	return &Dispatch{
		cfg: cfg,
	}
}

func (d *Dispatch) Config() *config.Config {
	d.lock.Lock()
	cfg := d.cfg
	d.lock.Unlock()
	return cfg
}

func (d *Dispatch) SetConfig(cfg *config.Config) {
	d.lock.Lock()
	d.cfg = cfg
	d.lock.Unlock()
}

func (d *Dispatch) Run() {
	d.upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	if d.Config().Dev {
		d.upgrader.CheckOrigin = func(r *http.Request) bool {
			return true
		}
	}

	session.CookieName = "dispatch"

	d.states = newStateStore(d.SessionStore)
	go d.states.run()

	d.loadUsers()
	d.initFileServer()
	d.startHTTP()
}

func (d *Dispatch) loadUsers() {
	users, err := storage.LoadUsers(d.Store)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("[Init] %d users", len(users))

	for _, user := range users {
		go d.loadUser(user)
	}
}

func (d *Dispatch) loadUser(user *storage.User) {
	messageStore, err := d.GetMessageStore(user)
	if err != nil {
		log.Fatal(err)
	}
	user.SetMessageStore(messageStore)

	search, err := d.GetMessageSearchProvider(user)
	if err != nil {
		log.Fatal(err)
	}
	user.SetMessageSearchProvider(search)

	state := NewState(user, d)
	d.states.set(state)
	go state.run()

	channels, err := user.GetChannels()
	if err != nil {
		log.Fatal(err)
	}

	servers, err := user.GetServers()
	if err != nil {
		log.Fatal(err)
	}

	for _, server := range servers {
		i := connectIRC(server, state, user.GetLastIP())

		var joining []string
		for _, channel := range channels {
			if channel.Server == server.Host {
				joining = append(joining, channel.Name)
			}
		}
		i.Join(joining...)
	}
}

func (d *Dispatch) startHTTP() {
	cfg := d.Config()
	addr := cfg.Address
	port := cfg.Port

	if cfg.HTTPS.Enabled {
		portHTTPS := cfg.HTTPS.Port
		redirect := cfg.HTTPS.Redirect

		if redirect {
			log.Println("[HTTP] Listening on port", port, "(HTTPS Redirect)")
			go http.ListenAndServe(net.JoinHostPort(addr, port), d.createHTTPSRedirect(portHTTPS))
		}

		server := &http.Server{
			Addr:    net.JoinHostPort(addr, portHTTPS),
			Handler: d,
		}

		if d.certExists() {
			log.Println("[HTTPS] Listening on port", portHTTPS)
			server.ListenAndServeTLS(cfg.HTTPS.Cert, cfg.HTTPS.Key)
		} else if domain := cfg.LetsEncrypt.Domain; domain != "" {
			dir := storage.Path.LetsEncrypt()
			email := cfg.LetsEncrypt.Email
			lePort := cfg.LetsEncrypt.Port

			if cfg.LetsEncrypt.Proxy && lePort != "" && (port != "80" || !redirect) {
				log.Println("[HTTP] Listening on port 80 (Let's Encrypt Proxy))")
				go http.ListenAndServe(net.JoinHostPort(addr, "80"), http.HandlerFunc(d.letsEncryptProxy))
			}

			le, err := letsencrypt.Run(dir, domain, email, ":"+lePort)
			if err != nil {
				log.Fatal(err)
			}

			server.TLSConfig = &tls.Config{
				GetCertificate: le.GetCertificate,
			}

			log.Println("[HTTPS] Listening on port", portHTTPS)
			log.Fatal(server.ListenAndServeTLS("", ""))
		} else {
			log.Fatal("Could not locate SSL certificate or private key")
		}
	} else {
		if cfg.Dev {
			// The node dev server will proxy index page requests and
			// websocket connections to this port
			port = "1337"
		}
		log.Println("[HTTP] Listening on port", port)
		log.Fatal(http.ListenAndServe(net.JoinHostPort(addr, port), d))
	}
}

func (d *Dispatch) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		fail(w, http.StatusNotFound)
		return
	}

	if r.URL.Path == "/init" {
		referer, err := url.Parse(r.Header.Get("Referer"))
		if err != nil {
			fail(w, http.StatusInternalServerError)
			return
		}

		state := d.handleAuth(w, r, true, true)
		data := d.getIndexData(r, referer.EscapedPath(), state)

		writeJSON(w, r, data)
	} else if strings.HasPrefix(r.URL.Path, "/ws") {
		if !websocket.IsWebSocketUpgrade(r) {
			fail(w, http.StatusBadRequest)
			return
		}

		state := d.handleAuth(w, r, false, false)
		if state == nil {
			log.Println("[Auth] No state")
			fail(w, http.StatusInternalServerError)
			return
		}

		d.upgradeWS(w, r, state)
	} else {
		d.serveFiles(w, r)
	}
}

func (d *Dispatch) upgradeWS(w http.ResponseWriter, r *http.Request, state *State) {
	conn, err := d.upgrader.Upgrade(w, r, w.Header())
	if err != nil {
		log.Println(err)
		return
	}

	newWSHandler(conn, state, r).run()
}

func (d *Dispatch) createHTTPSRedirect(portHTTPS string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/.well-known/acme-challenge") {
			d.letsEncryptProxy(w, r)
			return
		}

		host, _, err := net.SplitHostPort(r.Host)
		if err != nil {
			host = r.Host
		}

		u := url.URL{
			Scheme: "https",
			Host:   net.JoinHostPort(host, portHTTPS),
			Path:   r.RequestURI,
		}

		w.Header().Set("Location", u.String())
		w.WriteHeader(http.StatusMovedPermanently)
	})
}

func (d *Dispatch) letsEncryptProxy(w http.ResponseWriter, r *http.Request) {
	host, _, err := net.SplitHostPort(r.Host)
	if err != nil {
		host = r.Host
	}

	upstream := &url.URL{
		Scheme: "http",
		Host:   net.JoinHostPort(host, d.Config().LetsEncrypt.Port),
	}

	httputil.NewSingleHostReverseProxy(upstream).ServeHTTP(w, r)
}

func fail(w http.ResponseWriter, code int) {
	http.Error(w, http.StatusText(code), code)
}
