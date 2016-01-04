package server

import (
	"crypto/tls"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/khlieng/dispatch/Godeps/_workspace/src/github.com/spf13/viper"
)

type restartableHTTPS struct {
	listener net.Listener
	handler  http.Handler
	addr     string
	cert     string
	key      string
}

func (r *restartableHTTPS) start() error {
	var err error

	config := &tls.Config{
		NextProtos:   []string{"http/1.1"},
		Certificates: make([]tls.Certificate, 1),
	}

	config.Certificates[0], err = tls.LoadX509KeyPair(r.cert, r.key)
	if err != nil {
		return err
	}

	ln, err := net.Listen("tcp", r.addr)
	if err != nil {
		return err
	}

	r.listener = tls.NewListener(tcpKeepAliveListener{ln.(*net.TCPListener)}, config)
	return http.Serve(r.listener, r.handler)
}

func (r *restartableHTTPS) stop() {
	r.listener.Close()
}

func (r *restartableHTTPS) restart() {
	r.stop()
	go r.start()
}

type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln tcpKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}

func certExists() bool {
	cert := viper.GetString("https.cert")
	key := viper.GetString("https.key")

	if cert == "" || key == "" {
		return false
	}

	if _, err := os.Stat(cert); err != nil {
		return false
	}
	if _, err := os.Stat(key); err != nil {
		return false
	}

	return true
}
