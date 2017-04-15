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
	"github.com/spf13/viper"

	"github.com/khlieng/dispatch/letsencrypt"
	"github.com/khlieng/dispatch/storage"
)

var (
	sessions     *sessionStore
	channelStore *storage.ChannelStore

	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

func Run() {
	sessions = newSessionStore()
	channelStore = storage.NewChannelStore()

	if viper.GetBool("dev") {
		upgrader.CheckOrigin = func(r *http.Request) bool {
			return true
		}
	}

	reconnectIRC()
	initAuth()
	initFileServer()
	startHTTP()
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

			letsEncrypt, err := letsencrypt.Run(dir, domain, email, ":"+lePort)
			if err != nil {
				log.Fatal(err)
			}

			server.TLSConfig = &tls.Config{
				GetCertificate: letsEncrypt.GetCertificate,
			}

			log.Println("[HTTPS] Listening on port", portHTTPS)
			log.Fatal(server.ListenAndServeTLS("", ""))
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
		w.WriteHeader(404)
		return
	}

	if r.URL.Path == "/ws" {
		session := handleAuth(w, r)
		if session == nil {
			log.Println("[Auth] No session")
			w.WriteHeader(500)
			return
		}

		upgradeWS(w, r, session)
	} else {
		serveFiles(w, r)
	}
}

func upgradeWS(w http.ResponseWriter, r *http.Request, session *Session) {
	conn, err := upgrader.Upgrade(w, r, w.Header())
	if err != nil {
		log.Println(err)
		return
	}

	newWSHandler(conn, session).run()
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
