package server

import (
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/khlieng/name_pending/Godeps/_workspace/src/github.com/gorilla/websocket"

	"github.com/khlieng/name_pending/irc"
	"github.com/khlieng/name_pending/storage"
)

var (
	channelStore *storage.ChannelStore
	sessions     map[string]*Session
	sessionLock  sync.Mutex

	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func Run(port int) {
	defer storage.Close()

	channelStore = storage.NewChannelStore()
	sessions = make(map[string]*Session)

	reconnect()

	log.Println("Listening on port", port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), handler{}))
}

type handler struct{}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		return
	}

	if r.URL.Path == "/ws" {
		upgradeWS(w, r)
	} else {
		serveFiles(w, r)
	}
}

func upgradeWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	handleWS(conn)
}

func reconnect() {
	for _, user := range storage.LoadUsers() {
		session := NewSession()
		session.user = user
		sessions[user.UUID] = session
		go session.write()

		channels := user.GetChannels()

		for _, server := range user.GetServers() {
			i := irc.NewClient(server.Nick, server.Username)
			i.TLS = server.TLS
			i.Password = server.Password
			i.Realname = server.Realname

			i.Connect(server.Address)
			session.setIRC(i.Host, i)
			go handleIRC(i, session)

			var joining []string
			for _, channel := range channels {
				if channel.Server == server.Address {
					joining = append(joining, channel.Name)
				}
			}
			i.Join(joining...)
		}
	}
}
