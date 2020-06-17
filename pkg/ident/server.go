package ident

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

var (
	// DefaultAddr is the address a Server listens on when no Addr is specified
	DefaultAddr = ":113"
	// DefaultTimeout is the the time a Server will wait before failing
	// reads and writes if no Timeout is specified
	DefaultTimeout = 5 * time.Second
)

// Server implements the server-side of the Ident protocol
type Server struct {
	// Addr is the host:port address to listen on
	Addr string
	// Timeout is the time to wait before failing reads and writes
	Timeout time.Duration

	entries  map[string]entry
	listener net.Listener
	lock     sync.Mutex
}

type entry struct {
	remoteHost string
	ident      string
}

func NewServer() *Server {
	return &Server{
		entries: map[string]entry{},
	}
}

func (s *Server) Listen() error {
	var err error

	addr := s.Addr
	if addr == "" {
		addr = DefaultAddr
	}

	s.listener, err = net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	defer s.listener.Close()

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			return err
		}

		go s.handle(conn)
	}
}

func (s *Server) Stop() error {
	return s.listener.Close()
}

func (s *Server) Add(local, remote net.Addr, ident string) {
	if local == nil || remote == nil {
		return
	}

	_, localPort, err := net.SplitHostPort(local.String())
	if err != nil {
		return
	}

	remoteHost, remotePort, err := net.SplitHostPort(remote.String())
	if err != nil {
		return
	}

	s.lock.Lock()
	s.entries[localPort+","+remotePort] = entry{
		remoteHost: remoteHost,
		ident:      ident,
	}
	s.lock.Unlock()
}

func (s *Server) Remove(local, remote net.Addr) {
	if local == nil || remote == nil {
		return
	}

	_, localPort, err := net.SplitHostPort(local.String())
	if err != nil {
		return
	}

	_, remotePort, err := net.SplitHostPort(remote.String())
	if err != nil {
		return
	}

	s.lock.Lock()
	delete(s.entries, localPort+","+remotePort)
	s.lock.Unlock()
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()

	timeout := s.Timeout
	if timeout == 0 {
		timeout = DefaultTimeout
	}

	scan := bufio.NewScanner(conn)
	scan.Buffer(make([]byte, 32), 32)

	conn.SetReadDeadline(time.Now().Add(timeout))
	if !scan.Scan() {
		return
	}
	query := scan.Text()

	s.lock.Lock()
	entry, ok := s.entries[strings.ReplaceAll(query, " ", "")]
	s.lock.Unlock()

	if ok {
		remoteHost, _, err := net.SplitHostPort(conn.RemoteAddr().String())
		if err != nil || remoteHost != entry.remoteHost {
			return
		}

		conn.SetWriteDeadline(time.Now().Add(timeout))
		conn.Write([]byte(fmt.Sprintf("%s : USERID : Dispatch : %s\r\n", query, entry.ident)))
	}
}
