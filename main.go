package main

import (
	"log"
	"net/http"
	"sync"

	"golang.org/x/net/websocket"

	"github.com/khlieng/irc/storage"
)

var (
	channelStore *storage.ChannelStore
	sessions     map[string]*Session
	sessionLock  sync.Mutex
)

func main() {
	defer storage.Cleanup()

	channelStore = storage.NewChannelStore()
	sessions = make(map[string]*Session)

	/*for _, user := range storage.LoadUsers() {
		channels := user.GetChannels()

		for _, server := range user.GetServers() {
			session := NewSession()
			session.user = user
			sessions[user.UUID] = session

			irc := NewIRC(server.Nick, server.Username)
			irc.TLS = true
			irc.Connect(server.Address)

			session.setIRC(irc.Host, irc)

			go session.write()
			go handleMessages(irc, session)

			var joining []string
			for _, channel := range channels {
				if channel.Server == server.Address {
					joining = append(joining, channel.Name)
				}
			}
			irc.Join(joining...)
		}
	}*/

	http.Handle("/", http.FileServer(http.Dir("client/dist")))
	http.Handle("/ws", websocket.Handler(handleWS))

	log.Println("Listening on port 1337")
	http.ListenAndServe(":1337", nil)
}
