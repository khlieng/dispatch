package ident

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
)

var (
	DefaultAddr = ":113"
)

type Server struct {
	Addr string

	idents   map[string]string
	listener net.Listener
	lock     sync.Mutex
}

func NewServer() *Server {
	return &Server{
		idents: map[string]string{},
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

func (s *Server) Add(local, remote, ident string) {
	s.lock.Lock()
	s.idents[local+","+remote] = ident
	s.lock.Unlock()
}

func (s *Server) Remove(local, remote string) {
	s.lock.Lock()
	delete(s.idents, local+","+remote)
	s.lock.Unlock()
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()

	scan := bufio.NewScanner(conn)
	if !scan.Scan() {
		return
	}

	line := scan.Text()
	ports := strings.ReplaceAll(line, " ", "")

	s.lock.Lock()
	ident, ok := s.idents[ports]
	s.lock.Unlock()

	if ok {
		conn.Write([]byte(fmt.Sprintf("%s : USERID : Dispatch : %s\r\n", line, ident)))
	}
}
