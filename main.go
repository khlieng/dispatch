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
	files        []File
)

type File struct {
	Path        string
	ContentType string
}

func reconnect() {
	for _, user := range storage.LoadUsers() {
		channels := user.GetChannels()

		for _, server := range user.GetServers() {
			session := NewSession()
			session.user = user
			sessions[user.UUID] = session

			irc := NewIRC(server.Nick, server.Username)
			irc.TLS = server.TLS
			irc.Password = server.Password
			irc.Realname = server.Realname

			err := irc.Connect(server.Address)
			if err != nil {
				log.Println(err)
			} else {
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
		}
	}
}

func serveFiles(w http.ResponseWriter, r *http.Request) {
	var ext string

	if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		w.Header().Set("Content-Encoding", "gzip")
		ext = ".gz"
	}

	if r.URL.Path == "/" {
		w.Header().Set("Content-Type", "text/html")
		r.URL.Path = "/index.html" + ext
		fs.ServeHTTP(w, r)
		return
	}

	for _, file := range files {
		if strings.HasSuffix(r.URL.Path, file.Path) {
			w.Header().Set("Content-Type", file.ContentType)
			r.URL.Path = file.Path + ext
			fs.ServeHTTP(w, r)
			return
		}
	}

	w.Header().Set("Content-Type", "text/html")
	r.URL.Path = "/index.html" + ext

	fs.ServeHTTP(w, r)
}

func main() {
	defer storage.Cleanup()

	channelStore = storage.NewChannelStore()
	sessions = make(map[string]*Session)
	fs = http.FileServer(http.Dir("client/dist"))

	files = []File{
		File{"/bundle.js", "text/javascript"},
		File{"/css/style.css", "text/css"},
		File{"/css/fontello.css", "text/css"},
		File{"/font/fontello.eot", "application/vnd.ms-fontobject"},
		File{"/font/fontello.svg", "image/svg+xml"},
		File{"/font/fontello.ttf", "application/x-font-ttf"},
		File{"/font/fontello.woff", "application/font-woff"},
	}

	//reconnect()

	router := httprouter.New()

	router.Handler("GET", "/ws", websocket.Handler(handleWS))
	router.NotFound = serveFiles

	log.Println("Listening on port 1337")
	http.ListenAndServe(":1337", router)
}
