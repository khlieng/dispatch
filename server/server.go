package server

import (
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/khlieng/dispatch/config"
	"github.com/khlieng/dispatch/pkg/netutil"
	"github.com/khlieng/dispatch/pkg/session"
	"github.com/khlieng/dispatch/storage"
	"github.com/mholt/certmagic"
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

	port := cfg.Port
	if cfg.Dev {
		// The node dev server will proxy index page requests and
		// websocket connections to this port
		port = "1337"
	}

	httpSrv := &http.Server{
		Addr: net.JoinHostPort(cfg.Address, port),
	}

	if cfg.HTTPS.Enabled {
		httpSrv.ReadTimeout = 5 * time.Second
		httpSrv.WriteTimeout = 5 * time.Second

		httpsSrv := &http.Server{
			Addr:              net.JoinHostPort(cfg.Address, cfg.HTTPS.Port),
			ReadHeaderTimeout: 5 * time.Second,
			WriteTimeout:      10 * time.Second,
			IdleTimeout:       120 * time.Second,
			Handler:           d,
		}

		redirect := createHTTPSRedirect(cfg.HTTPS.Port, d)

		if d.certExists() {
			httpSrv.Handler = redirect
			log.Println("[HTTP] Listening on port", port, "(HTTPS Redirect)")
			go httpSrv.ListenAndServe()

			log.Println("[HTTPS] Listening on port", cfg.HTTPS.Port)
			log.Fatal(httpsSrv.ListenAndServeTLS(cfg.HTTPS.Cert, cfg.HTTPS.Key))
		} else {
			cache := certmagic.NewCache(&certmagic.FileStorage{
				Path: storage.Path.LetsEncrypt(),
			})

			magic := certmagic.NewWithCache(cache, certmagic.Config{
				Agreed:     true,
				Email:      cfg.LetsEncrypt.Email,
				MustStaple: true,
			})

			domains := []string{cfg.LetsEncrypt.Domain}
			if cfg.LetsEncrypt.Domain == "" {
				domains = []string{}
				magic.OnDemand = &certmagic.OnDemandConfig{MaxObtain: 3}
			}

			err := magic.Manage(domains)
			if err != nil {
				log.Fatal(err)
			}

			tlsConfig := magic.TLSConfig()
			tlsConfig.MinVersion = tls.VersionTLS12
			tlsConfig.CipherSuites = getCipherSuites()
			tlsConfig.CurvePreferences = []tls.CurveID{
				tls.X25519,
				tls.CurveP256,
			}
			tlsConfig.PreferServerCipherSuites = true
			httpsSrv.TLSConfig = tlsConfig

			httpSrv.Handler = magic.HTTPChallengeHandler(redirect)
			log.Println("[HTTP] Listening on port", port, "(HTTPS Redirect)")
			go httpSrv.ListenAndServe()

			log.Println("[HTTPS] Listening on port", cfg.HTTPS.Port)
			log.Fatal(httpsSrv.ListenAndServeTLS("", ""))
		}
	} else {
		httpSrv.ReadHeaderTimeout = 5 * time.Second
		httpSrv.WriteTimeout = 10 * time.Second
		httpSrv.IdleTimeout = 120 * time.Second
		httpSrv.Handler = d

		log.Println("[HTTP] Listening on port", port)
		log.Fatal(httpSrv.ListenAndServe())
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

func createHTTPSRedirect(portHTTPS string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		host, _, err := net.SplitHostPort(r.Host)
		if err != nil {
			host = r.Host
		}

		if netutil.IsPrivate(host) {
			fallback.ServeHTTP(w, r)
			return
		}

		u := url.URL{
			Scheme: "https",
			Host:   net.JoinHostPort(host, portHTTPS),
			Path:   r.RequestURI,
		}

		w.Header().Set("Connection", "close")
		w.Header().Set("Location", u.String())
		w.WriteHeader(http.StatusMovedPermanently)
	}
}

func fail(w http.ResponseWriter, code int) {
	http.Error(w, http.StatusText(code), code)
}
