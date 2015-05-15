package server

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/khlieng/name_pending/Godeps/_workspace/src/github.com/gorilla/websocket"
	"github.com/khlieng/name_pending/Godeps/_workspace/src/github.com/julienschmidt/httprouter"

	"github.com/khlieng/name_pending/storage"
)

var (
	channelStore *storage.ChannelStore
	sessions     map[string]*Session
	sessionLock  sync.Mutex

	files = []File{
		File{"bundle.js", "text/javascript"},
		File{"bundle.css", "text/css"},
		File{"font/fontello.eot", "application/vnd.ms-fontobject"},
		File{"font/fontello.svg", "image/svg+xml"},
		File{"font/fontello.ttf", "application/x-font-ttf"},
		File{"font/fontello.woff", "application/font-woff"},
	}

	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

type File struct {
	Path        string
	ContentType string
}

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

func serveFiles(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		serveFile("index.html.gz", "text/html", w, r)
		return
	}

	for _, file := range files {
		if strings.HasSuffix(r.URL.Path, file.Path) {
			serveFile(file.Path+".gz", file.ContentType, w, r)
			return
		}
	}

	serveFile("index.html.gz", "text/html", w, r)
}

func serveFile(path, contentType string, w http.ResponseWriter, r *http.Request) {
	data, _ := Asset(path)

	w.Header().Set("Content-Type", contentType)

	if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Set("Content-Length", strconv.Itoa(len(data)))
		w.Write(data)
	} else {
		gzr, _ := gzip.NewReader(bytes.NewReader(data))
		buf, _ := ioutil.ReadAll(gzr)
		w.Header().Set("Content-Length", strconv.Itoa(len(buf)))
		w.Write(buf)
	}
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
