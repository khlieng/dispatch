package server

import (
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"

	"github.com/khlieng/dispatch/Godeps/_workspace/src/github.com/gorilla/websocket"
	"github.com/khlieng/dispatch/Godeps/_workspace/src/github.com/spf13/viper"
	"github.com/khlieng/dispatch/Godeps/_workspace/src/golang.org/x/net/http2"

	"github.com/khlieng/dispatch/letsencrypt"
	"github.com/khlieng/dispatch/storage"
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

func Run() {
	defer storage.Close()

	channelStore = storage.NewChannelStore()
	sessions = make(map[string]*Session)

	reconnectIRC()
	startHTTP()

	select {}
}

func startHTTP() {
	port := viper.GetString("port")

	if viper.GetBool("https.enabled") {
		portHTTPS := viper.GetString("https.port")
		redirect := viper.GetBool("https.redirect")

		if redirect {
			log.Println("[HTTP] Listening on port", port, "(HTTPS Redirect)")
			go http.ListenAndServe(":"+port, createHTTPSRedirect(portHTTPS))
		}

		server := &http.Server{
			Addr:    ":" + portHTTPS,
			Handler: http.HandlerFunc(serve),
		}

		http2.ConfigureServer(server, nil)

		if certExists() {
			log.Println("[HTTPS] Listening on port", portHTTPS)
			server.ListenAndServeTLS(viper.GetString("https.cert"), viper.GetString("https.key"))
		} else if domain := viper.GetString("letsencrypt.domain"); domain != "" {
			dir := storage.Path.LetsEncrypt()
			email := viper.GetString("letsencrypt.email")
			lePort := viper.GetString("letsencrypt.port")

			if viper.GetBool("letsencrypt.proxy") && lePort != "" && (port != "80" || !redirect) {
				log.Println("[HTTP] Listening on port 80 (Let's Encrypt Proxy))")
				go http.ListenAndServe(":80", http.HandlerFunc(letsEncryptProxy))
			}

			letsEncrypt, err := letsencrypt.Run(dir, domain, email, lePort)
			if err != nil {
				log.Fatal(err)
			}

			server.TLSConfig.GetCertificate = letsEncrypt.GetCertificate

			log.Println("[HTTPS] Listening on port", portHTTPS)
			log.Fatal(listenAndServeTLS(server))
		} else {
			log.Fatal("Could not locate SSL certificate or private key")
		}
	} else {
		log.Println("[HTTP] Listening on port", port)
		log.Fatal(http.ListenAndServe(":"+port, http.HandlerFunc(serve)))
	}
}

func serve(w http.ResponseWriter, r *http.Request) {
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

	uuid := r.URL.Query().Get("uuid")
	if uuid != "" {
		newWSHandler(conn, uuid).run()
	}
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
