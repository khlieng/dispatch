package server

import (
	"sync"

	"github.com/khlieng/dispatch/irc"
	"github.com/khlieng/dispatch/storage"
)

type Session struct {
	irc     map[string]*irc.Client
	ircLock sync.Mutex

	ws     map[string]*wsConn
	wsLock sync.Mutex
	out    chan WSResponse

	user *storage.User
}

func NewSession() *Session {
	return &Session{
		irc: make(map[string]*irc.Client),
		ws:  make(map[string]*wsConn),
		out: make(chan WSResponse, 32),
	}
}

func (s *Session) getIRC(server string) (*irc.Client, bool) {
	s.ircLock.Lock()
	i, ok := s.irc[server]
	s.ircLock.Unlock()

	return i, ok
}

func (s *Session) setIRC(server string, i *irc.Client) {
	s.ircLock.Lock()
	s.irc[server] = i
	s.ircLock.Unlock()
}

func (s *Session) deleteIRC(server string) {
	s.ircLock.Lock()
	delete(s.irc, server)
	s.ircLock.Unlock()
}

func (s *Session) numIRC() int {
	s.ircLock.Lock()
	n := len(s.irc)
	s.ircLock.Unlock()

	return n
}

func (s *Session) setWS(addr string, w *wsConn) {
	s.wsLock.Lock()
	s.ws[addr] = w
	s.wsLock.Unlock()
}

func (s *Session) deleteWS(addr string) {
	s.wsLock.Lock()
	delete(s.ws, addr)
	s.wsLock.Unlock()
}

func (s *Session) sendJSON(t string, v interface{}) {
	s.out <- WSResponse{t, v}
}

func (s *Session) sendError(err error, server string) {
	s.sendJSON("error", Error{
		Server:  server,
		Message: err.Error(),
	})
}

func (s *Session) write() {
	for res := range s.out {
		s.wsLock.Lock()
		for _, ws := range s.ws {
			ws.out <- res
		}
		s.wsLock.Unlock()
	}
}
