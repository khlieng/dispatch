package server

import (
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/khlieng/dispatch/config"
	"github.com/khlieng/dispatch/pkg/https"
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

	port := cfg.Port
	if cfg.Dev {
		// The node dev server will proxy index page requests and
		// websocket connections to this port
		port = "1337"
	}

	if cfg.HTTPS.Enabled {
		log.Println("[HTTP] Listening on port", port, "(HTTPS Redirect)")
		log.Println("[HTTPS] Listening on port", cfg.HTTPS.Port)
	} else {
		log.Println("[HTTP] Listening on port", port)
	}

	log.Fatal(https.Serve(d, https.Config{
		Addr:      cfg.Address,
		PortHTTP:  port,
		PortHTTPS: cfg.HTTPS.Port,
		HTTPOnly:  !cfg.HTTPS.Enabled,

		StoragePath: storage.Path.LetsEncrypt(),
		Domain:      cfg.LetsEncrypt.Domain,
		Email:       cfg.LetsEncrypt.Email,

		Cert: cfg.HTTPS.Cert,
		Key:  cfg.HTTPS.Key,
	}))
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

func fail(w http.ResponseWriter, code int) {
	http.Error(w, http.StatusText(code), code)
}
