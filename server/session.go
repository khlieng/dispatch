package server

import (
	"sync"
	"time"

	"fmt"

	"github.com/khlieng/dispatch/irc"
	"github.com/khlieng/dispatch/storage"
)

const (
	AnonymousSessionExpiration = 24 * time.Hour
)

type Session struct {
	irc             map[string]*irc.Client
	connectionState map[string]bool
	ircLock         sync.Mutex

	ws        map[string]*wsConn
	wsLock    sync.Mutex
	broadcast chan WSResponse

	user       *storage.User
	expiration *time.Timer
	reset      chan time.Duration
}

func NewSession(user *storage.User) *Session {
	return &Session{
		irc:             make(map[string]*irc.Client),
		connectionState: make(map[string]bool),
		ws:              make(map[string]*wsConn),
		broadcast:       make(chan WSResponse, 32),
		user:            user,
		expiration:      time.NewTimer(AnonymousSessionExpiration),
		reset:           make(chan time.Duration, 1),
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
	s.connectionState[server] = false
	s.ircLock.Unlock()

	s.reset <- 0
}

func (s *Session) deleteIRC(server string) {
	s.ircLock.Lock()
	delete(s.irc, server)
	delete(s.connectionState, server)
	s.ircLock.Unlock()

	s.resetExpirationIfEmpty()
}

func (s *Session) numIRC() int {
	s.ircLock.Lock()
	n := len(s.irc)
	s.ircLock.Unlock()

	return n
}

func (s *Session) getConnectionStates() map[string]bool {
	s.ircLock.Lock()
	state := make(map[string]bool, len(s.connectionState))

	for k, v := range s.connectionState {
		state[k] = v
	}
	s.ircLock.Unlock()

	return state
}

func (s *Session) setConnectionState(server string, connected bool) {
	s.ircLock.Lock()
	s.connectionState[server] = connected
	s.ircLock.Unlock()
}

func (s *Session) setWS(addr string, w *wsConn) {
	s.wsLock.Lock()
	s.ws[addr] = w
	s.wsLock.Unlock()

	s.reset <- 0
}

func (s *Session) deleteWS(addr string) {
	s.wsLock.Lock()
	delete(s.ws, addr)
	s.wsLock.Unlock()

	s.resetExpirationIfEmpty()
}

func (s *Session) numWS() int {
	s.ircLock.Lock()
	n := len(s.ws)
	s.ircLock.Unlock()

	return n
}

func (s *Session) sendJSON(t string, v interface{}) {
	s.broadcast <- WSResponse{t, v}
}

func (s *Session) sendError(err error, server string) {
	s.sendJSON("error", Error{
		Server:  server,
		Message: err.Error(),
	})
}

func (s *Session) print(server string, a ...interface{}) {
	s.sendJSON("print", Message{
		Server:  server,
		Content: fmt.Sprintln(a...),
	})
}

func (s *Session) resetExpirationIfEmpty() {
	if s.numIRC() == 0 && s.numWS() == 0 {
		s.reset <- AnonymousSessionExpiration
	}
}

func (s *Session) run() {
	for {
		select {
		case res := <-s.broadcast:
			s.wsLock.Lock()
			for _, ws := range s.ws {
				ws.out <- res
			}
			s.wsLock.Unlock()

		case <-s.expiration.C:
			sessions.delete(s.user.ID)
			s.user.Remove()
			return

		case duration := <-s.reset:
			if duration == 0 {
				s.expiration.Stop()
			} else {
				s.expiration.Reset(duration)
			}
		}
	}
}

type sessionStore struct {
	sessions map[uint64]*Session
	lock     sync.Mutex
}

func newSessionStore() *sessionStore {
	return &sessionStore{
		sessions: make(map[uint64]*Session),
	}
}

func (s *sessionStore) get(userid uint64) *Session {
	s.lock.Lock()
	session := s.sessions[userid]
	s.lock.Unlock()
	return session
}

func (s *sessionStore) set(userid uint64, session *Session) {
	s.lock.Lock()
	s.sessions[userid] = session
	s.lock.Unlock()
}

func (s *sessionStore) delete(userid uint64) {
	s.lock.Lock()
	delete(s.sessions, userid)
	s.lock.Unlock()
}
