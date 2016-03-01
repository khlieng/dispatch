package server

import (
	"crypto/tls"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/spf13/viper"
)

func listenAndServeTLS(srv *http.Server) error {
	if srv.TLSConfig.NextProtos == nil {
		srv.TLSConfig.NextProtos = []string{"http/1.1"}
	}

	ln, err := net.Listen("tcp", srv.Addr)
	if err != nil {
		return err
	}

	tlsListener := tls.NewListener(tcpKeepAliveListener{ln.(*net.TCPListener)}, srv.TLSConfig)
	return srv.Serve(tlsListener)
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
