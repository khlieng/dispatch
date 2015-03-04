package main

import (
	"encoding/json"
	"sync"

	"github.com/khlieng/name_pending/Godeps/_workspace/src/golang.org/x/net/websocket"
	"github.com/khlieng/name_pending/storage"
)

type Session struct {
	irc     map[string]*IRC
	ircLock sync.Mutex

	ws     map[string]*WebSocket
	wsLock sync.Mutex
	out    chan []byte

	user storage.User
}

func NewSession() *Session {
	return &Session{
		irc: make(map[string]*IRC),
		ws:  make(map[string]*WebSocket),
		out: make(chan []byte, 32),
	}
}

func (s *Session) getIRC(server string) (*IRC, bool) {
	s.ircLock.Lock()
	irc, ok := s.irc[server]
	s.ircLock.Unlock()

	return irc, ok
}

func (s *Session) setIRC(server string, irc *IRC) {
	s.ircLock.Lock()
	s.irc[server] = irc
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

func (s *Session) setWS(addr string, ws *websocket.Conn) {
	socket := NewWebSocket(ws)
	go socket.write()

	s.wsLock.Lock()
	s.ws[addr] = socket
	s.wsLock.Unlock()
}

func (s *Session) deleteWS(addr string) {
	s.wsLock.Lock()
	delete(s.ws, addr)
	s.wsLock.Unlock()
}

func (s *Session) sendJSON(t string, v interface{}) {
	data, _ := json.Marshal(v)
	raw := json.RawMessage(data)
	res, _ := json.Marshal(WSResponse{Type: t, Response: &raw})

	s.out <- res
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
			ws.Out <- res
		}
		s.wsLock.Unlock()
	}
}
