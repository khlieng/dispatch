package server

import (
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/khlieng/dispatch/pkg/letsencrypt"
	"github.com/khlieng/dispatch/pkg/session"
	"github.com/khlieng/dispatch/storage"
	"github.com/mailru/easyjson"
	"github.com/spf13/viper"
)

var channelStore = storage.NewChannelStore()

type Dispatch struct {
	Store        storage.Store
	SessionStore storage.SessionStore

	GetMessageStore          func(*storage.User) (storage.MessageStore, error)
	GetMessageSearchProvider func(*storage.User) (storage.MessageSearchProvider, error)

	upgrader websocket.Upgrader
	states   *stateStore
}

func (d *Dispatch) Run() {
	d.upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	if viper.GetBool("dev") {
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
	addr := viper.GetString("address")
	port := viper.GetString("port")

	if viper.GetBool("https.enabled") {
		portHTTPS := viper.GetString("https.port")
		redirect := viper.GetBool("https.redirect")

		if redirect {
			log.Println("[HTTP] Listening on port", port, "(HTTPS Redirect)")
			go http.ListenAndServe(net.JoinHostPort(addr, port), createHTTPSRedirect(portHTTPS))
		}

		server := &http.Server{
			Addr:    net.JoinHostPort(addr, portHTTPS),
			Handler: http.HandlerFunc(d.serve),
		}

		if certExists() {
			log.Println("[HTTPS] Listening on port", portHTTPS)
			server.ListenAndServeTLS(viper.GetString("https.cert"), viper.GetString("https.key"))
		} else if domain := viper.GetString("letsencrypt.domain"); domain != "" {
			dir := storage.Path.LetsEncrypt()
			email := viper.GetString("letsencrypt.email")
			lePort := viper.GetString("letsencrypt.port")

			if viper.GetBool("letsencrypt.proxy") && lePort != "" && (port != "80" || !redirect) {
				log.Println("[HTTP] Listening on port 80 (Let's Encrypt Proxy))")
				go http.ListenAndServe(net.JoinHostPort(addr, "80"), http.HandlerFunc(letsEncryptProxy))
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
		if viper.GetBool("dev") {
			// The node dev server will proxy index page requests and
			// websocket connections to this port
			port = "1337"
		}
		log.Println("[HTTP] Listening on port", port)
		log.Fatal(http.ListenAndServe(net.JoinHostPort(addr, port), http.HandlerFunc(d.serve)))
	}
}

func (d *Dispatch) serve(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		fail(w, http.StatusNotFound)
		return
	}

	if strings.HasPrefix(r.URL.Path, "/ws") {
		if !websocket.IsWebSocketUpgrade(r) {
			fail(w, http.StatusBadRequest)
			return
		}

		state := d.handleAuth(w, r, true)
		if state == nil {
			log.Println("[Auth] No state")
			fail(w, http.StatusInternalServerError)
			return
		}

		d.upgradeWS(w, r, state)
	} else if strings.HasPrefix(r.URL.Path, "/data") {
		state := d.handleAuth(w, r, true)
		if state == nil {
			log.Println("[Auth] No state")
			fail(w, http.StatusInternalServerError)
			return
		}

		easyjson.MarshalToHTTPResponseWriter(getIndexData(r, state), w)
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

func createHTTPSRedirect(portHTTPS string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/.well-known/acme-challenge") {
			letsEncryptProxy(w, r)
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

func letsEncryptProxy(w http.ResponseWriter, r *http.Request) {
	host, _, err := net.SplitHostPort(r.Host)
	if err != nil {
		host = r.Host
	}

	upstream := &url.URL{
		Scheme: "http",
		Host:   net.JoinHostPort(host, viper.GetString("letsencrypt.port")),
	}

	httputil.NewSingleHostReverseProxy(upstream).ServeHTTP(w, r)
}

func fail(w http.ResponseWriter, code int) {
	http.Error(w, http.StatusText(code), code)
}
