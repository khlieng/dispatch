package server

import (
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/khlieng/name_pending/Godeps/_workspace/src/github.com/gorilla/websocket"
	"github.com/khlieng/name_pending/Godeps/_workspace/src/github.com/julienschmidt/httprouter"

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

	router := httprouter.New()

	router.HandlerFunc("GET", "/ws", upgradeWS)
	router.NotFound = serveFiles

	log.Println("Listening on port", port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), router))
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
			irc := NewIRC(server.Nick, server.Username)
			irc.TLS = server.TLS
			irc.Password = server.Password
			irc.Realname = server.Realname

			go func() {
				err := irc.Connect(server.Address)
				if err != nil {
					log.Println(err)
				} else {
					session.setIRC(irc.Host, irc)

					go handleMessages(irc, session)

					var joining []string
					for _, channel := range channels {
						if channel.Server == server.Address {
							joining = append(joining, channel.Name)
						}
					}
					irc.Join(joining...)
				}
			}()
		}
	}
}
