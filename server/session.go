package server

import (
	"crypto/rand"
	"encoding/base64"
	"sync"
	"time"

	"fmt"

	"github.com/khlieng/dispatch/irc"
	"github.com/khlieng/dispatch/storage"
)

const (
	AnonymousSessionExpiration = 1 * time.Minute
)

type Session struct {
	irc             map[string]*irc.Client
	connectionState map[string]irc.ConnectionState
	ircLock         sync.Mutex

	ws        map[string]*wsConn
	wsLock    sync.Mutex
	broadcast chan WSResponse

	id         string
	user       *storage.User
	expiration *time.Timer
	reset      chan time.Duration
}

func NewSession(user *storage.User) (*Session, error) {
	id, err := newSessionID()
	if err != nil {
		return nil, err
	}
	return &Session{
		irc:             make(map[string]*irc.Client),
		connectionState: make(map[string]irc.ConnectionState),
		ws:              make(map[string]*wsConn),
		broadcast:       make(chan WSResponse, 32),
		id:              id,
		user:            user,
		expiration:      time.NewTimer(AnonymousSessionExpiration),
		reset:           make(chan time.Duration, 1),
	}, nil
}

func newSessionID() (string, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	return base64.RawURLEncoding.EncodeToString(key), err
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
	s.connectionState[server] = irc.ConnectionState{
		Connected: false,
	}
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

func (s *Session) getConnectionStates() map[string]irc.ConnectionState {
	s.ircLock.Lock()
	state := make(map[string]irc.ConnectionState, len(s.connectionState))

	for k, v := range s.connectionState {
		state[k] = v
	}
	s.ircLock.Unlock()

	return state
}

func (s *Session) setConnectionState(server string, state irc.ConnectionState) {
	s.ircLock.Lock()
	s.connectionState[server] = state
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

func (s *Session) sendLastMessages(server, channel string, count int) {
	messages, hasMore, err := s.user.GetLastMessages(server, channel, count)
	if err == nil && len(messages) > 0 {
		res := Messages{
			Server:   server,
			To:       channel,
			Messages: messages,
		}

		if hasMore {
			res.Next = messages[0].ID
		}

		s.sendJSON("messages", res)
	}
}

func (s *Session) sendMessages(server, channel string, count int, fromID string) {
	messages, hasMore, err := s.user.GetMessages(server, channel, count, fromID)
	if err == nil && len(messages) > 0 {
		res := Messages{
			Server:   server,
			To:       channel,
			Messages: messages,
			Prepend:  true,
		}

		if hasMore {
			res.Next = messages[0].ID
		}

		s.sendJSON("messages", res)
	}
}

func (s *Session) print(a ...interface{}) {
	s.sendJSON("print", Message{
		Content: fmt.Sprintln(a...),
	})
}

func (s *Session) printError(a ...interface{}) {
	s.sendJSON("print", Message{
		Content: fmt.Sprintln(a...),
		Type:    "error",
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
			sessions.delete(s.id)
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
	sessions map[string]*Session
	lock     sync.Mutex
}

func newSessionStore() *sessionStore {
	return &sessionStore{
		sessions: make(map[string]*Session),
	}
}

func (s *sessionStore) get(id string) *Session {
	s.lock.Lock()
	session := s.sessions[id]
	s.lock.Unlock()
	return session
}

func (s *sessionStore) set(session *Session) {
	s.lock.Lock()
	s.sessions[session.id] = session
	s.lock.Unlock()
}

func (s *sessionStore) delete(id string) {
	s.lock.Lock()
	delete(s.sessions, id)
	s.lock.Unlock()
}
