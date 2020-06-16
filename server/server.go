package server

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/khlieng/dispatch/config"
	"github.com/khlieng/dispatch/pkg/https"
	"github.com/khlieng/dispatch/pkg/ident"
	"github.com/khlieng/dispatch/pkg/session"
	"github.com/khlieng/dispatch/storage"
)

var channelIndexes = storage.NewChannelIndexManager()

type Dispatch struct {
	Store        storage.Store
	SessionStore storage.SessionStore

	cfg      *config.Config
	upgrader websocket.Upgrader
	states   *stateStore
	identd   *ident.Server
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
	cfg := d.Config()

	d.upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	if cfg.Dev {
		d.upgrader.CheckOrigin = func(r *http.Request) bool {
			return true
		}
	}

	if cfg.Identd {
		d.identd = ident.NewServer()
		go d.identd.Listen()
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
	state := NewState(user, d)
	d.states.set(state)
	go state.run()

	networks, err := user.Networks()
	if err != nil {
		log.Fatal(err)
	}

	channels, err := user.Channels()
	if err != nil {
		log.Fatal(err)
	}

	for _, network := range networks {
		i := connectIRC(network, state, user.GetLastIP())

		var joining []string
		for _, channel := range channels {
			if channel.Network == network.Host {
				network.AddChannel(network.NewChannel(channel.Name))
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
		// The node dev network will proxy index page requests and
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
		state := d.handleAuth(w, r, true, true)
		data := d.getIndexData(r, state)

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
	} else if strings.HasPrefix(r.URL.Path, "/downloads") {
		state := d.handleAuth(w, r, false, false)
		if state == nil {
			log.Println("[Auth] No state")
			fail(w, http.StatusInternalServerError)
			return
		}

		params := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

		if len(params) == 3 {
			userID, err := strconv.ParseUint(params[1], 10, 64)
			if err != nil {
				fail(w, http.StatusBadRequest)
			}

			if userID != state.user.ID {
				fail(w, http.StatusUnauthorized)
			}

			filename := params[2]
			w.Header().Set("Content-Disposition", "attachment; filename="+filename)

			if pack, ok := state.pendingDCC(filename); ok {
				state.deletePendingDCC(filename)

				w.Header().Set("Content-Length", strconv.FormatUint(pack.Length, 10))
				pack.Download(w, nil)
			} else {
				file := storage.Path.DownloadedFile(state.user.Username, filename)
				http.ServeFile(w, r, file)

				if d.Config().DCC.Autoget.Delete {
					os.Remove(file)
				}
			}
		} else {
			fail(w, http.StatusNotFound)
		}
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
