package main

import (
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/net/websocket"

	"github.com/khlieng/name_pending/storage"
)

var (
	channelStore *storage.ChannelStore
	sessions     map[string]*Session
	sessionLock  sync.Mutex
	fs           http.Handler
)

func serveFiles(w http.ResponseWriter, r *http.Request) {
	var ext string

	if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		w.Header().Set("Content-Encoding", "gzip")
		ext = ".gz"
	}

	if strings.HasSuffix(r.URL.Path, "bundle.js") {
		w.Header().Set("Content-Type", "text/javascript")
		r.URL.Path = "/bundle.js" + ext
	} else if strings.HasSuffix(r.URL.Path, "style.css") {
		w.Header().Set("Content-Type", "text/css")
		r.URL.Path = "/style.css" + ext
	} else {
		w.Header().Set("Content-Type", "text/html")
		r.URL.Path = "/index.html" + ext
	}

	fs.ServeHTTP(w, r)
}

func main() {
	defer storage.Cleanup()

	channelStore = storage.NewChannelStore()
	sessions = make(map[string]*Session)
	fs = http.FileServer(http.Dir("client/dist"))

	/*for _, user := range storage.LoadUsers() {
		channels := user.GetChannels()

		for _, server := range user.GetServers() {
			session := NewSession()
			session.user = user
			sessions[user.UUID] = session

			irc := NewIRC(server.Nick, server.Username)
			irc.TLS = server.TLS
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

	router := httprouter.New()

	router.Handler("GET", "/ws", websocket.Handler(handleWS))
	router.NotFound = serveFiles

	log.Println("Listening on port 1337")
	http.ListenAndServe(":1337", router)
}
