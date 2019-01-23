package server

import (
	"log"
	"sync"
	"time"

	"fmt"

	"github.com/khlieng/dispatch/pkg/irc"
	"github.com/khlieng/dispatch/pkg/session"
	"github.com/khlieng/dispatch/storage"
)

const (
	// AnonymousUserExpiration is the time to wait before removing an anonymous
	// user that has no irc or websocket connections
	AnonymousUserExpiration = 1 * time.Minute
)

// State is the live state of a single user
type State struct {
	stateData

	irc             map[string]*irc.Client
	connectionState map[string]irc.ConnectionState
	ircLock         sync.Mutex

	ws        map[string]*wsConn
	wsLock    sync.Mutex
	broadcast chan WSResponse

	srv        *Dispatch
	user       *storage.User
	expiration *time.Timer
	reset      chan time.Duration
}

func NewState(user *storage.User, srv *Dispatch) *State {
	return &State{
		stateData:       stateData{m: map[string]interface{}{}},
		irc:             make(map[string]*irc.Client),
		connectionState: make(map[string]irc.ConnectionState),
		ws:              make(map[string]*wsConn),
		broadcast:       make(chan WSResponse, 32),
		srv:             srv,
		user:            user,
		expiration:      time.NewTimer(AnonymousUserExpiration),
		reset:           make(chan time.Duration, 1),
	}
}

func (s *State) getIRC(server string) (*irc.Client, bool) {
	s.ircLock.Lock()
	i, ok := s.irc[server]
	s.ircLock.Unlock()

	return i, ok
}

func (s *State) setIRC(server string, i *irc.Client) {
	s.ircLock.Lock()
	s.irc[server] = i
	s.connectionState[server] = irc.ConnectionState{
		Connected: false,
	}
	s.ircLock.Unlock()

	s.reset <- 0
}

func (s *State) deleteIRC(server string) {
	s.ircLock.Lock()
	delete(s.irc, server)
	delete(s.connectionState, server)
	s.ircLock.Unlock()

	s.resetExpirationIfEmpty()
}

func (s *State) numIRC() int {
	s.ircLock.Lock()
	n := len(s.irc)
	s.ircLock.Unlock()

	return n
}

func (s *State) getConnectionStates() map[string]irc.ConnectionState {
	s.ircLock.Lock()
	state := make(map[string]irc.ConnectionState, len(s.connectionState))

	for k, v := range s.connectionState {
		state[k] = v
	}
	s.ircLock.Unlock()

	return state
}

func (s *State) setConnectionState(server string, state irc.ConnectionState) {
	s.ircLock.Lock()
	s.connectionState[server] = state
	s.ircLock.Unlock()
}

func (s *State) setWS(addr string, w *wsConn) {
	s.wsLock.Lock()
	s.ws[addr] = w
	s.wsLock.Unlock()

	s.reset <- 0
}

func (s *State) deleteWS(addr string) {
	s.wsLock.Lock()
	delete(s.ws, addr)
	s.wsLock.Unlock()

	s.resetExpirationIfEmpty()
}

func (s *State) numWS() int {
	s.ircLock.Lock()
	n := len(s.ws)
	s.ircLock.Unlock()

	return n
}

func (s *State) sendJSON(t string, v interface{}) {
	s.broadcast <- WSResponse{t, v}
}

func (s *State) sendError(err error, server string) {
	s.sendJSON("error", Error{
		Server:  server,
		Message: err.Error(),
	})
}

func (s *State) sendLastMessages(server, channel string, count int) {
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

func (s *State) sendMessages(server, channel string, count int, fromID string) {
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

func (s *State) print(a ...interface{}) {
	s.sendJSON("print", Message{
		Content: fmt.Sprintln(a...),
	})
}

func (s *State) printError(a ...interface{}) {
	s.sendJSON("print", Message{
		Content: fmt.Sprintln(a...),
		Type:    "error",
	})
}

func (s *State) resetExpirationIfEmpty() {
	if s.numIRC() == 0 && s.numWS() == 0 {
		s.reset <- AnonymousUserExpiration
	}
}

func (s *State) kill() {
	s.wsLock.Lock()
	for _, ws := range s.ws {
		ws.conn.Close()
	}
	s.wsLock.Unlock()
	s.ircLock.Lock()
	for _, i := range s.irc {
		i.Quit()
	}
	s.ircLock.Unlock()
}

func (s *State) run() {
	for {
		select {
		case res := <-s.broadcast:
			s.wsLock.Lock()
			for _, ws := range s.ws {
				ws.out <- res
			}
			s.wsLock.Unlock()

		case <-s.expiration.C:
			s.srv.states.delete(s.user.ID)
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

type stateData struct {
	m    map[string]interface{}
	lock sync.Mutex
}

func (s stateData) Get(key string) interface{} {
	s.lock.Lock()
	v := s.m[key]
	s.lock.Unlock()
	return v
}

func (s stateData) Set(key string, value interface{}) {
	s.lock.Lock()
	s.m[key] = value
	s.lock.Unlock()
}

func (s stateData) String(key string) string {
	return s.Get(key).(string)
}

func (s stateData) Int(key string) int {
	return s.Get(key).(int)
}

func (s stateData) Bool(key string) bool {
	return s.Get(key).(bool)
}

type stateStore struct {
	states       map[uint64]*State
	sessions     map[string]*session.Session
	sessionStore storage.SessionStore
	lock         sync.Mutex
}

func newStateStore(sessionStore storage.SessionStore) *stateStore {
	store := &stateStore{
		states:       make(map[uint64]*State),
		sessions:     make(map[string]*session.Session),
		sessionStore: sessionStore,
	}

	sessions, err := sessionStore.GetSessions()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("[Init] %d sessions", len(sessions))

	for _, session := range sessions {
		if !session.Expired() {
			store.sessions[session.Key()] = session
		} else {
			go sessionStore.DeleteSession(session.Key())
		}
	}

	return store
}

func (s *stateStore) run() {
	pruneSessions := time.Tick(time.Minute * 5)
	for {
		select {
		case <-pruneSessions:
			s.lock.Lock()
			for key, session := range s.sessions {
				if session.Expired() {
					s.internalDeleteSession(key)
				}
			}
			s.lock.Unlock()
		}
	}
}

func (s *stateStore) get(id uint64) *State {
	s.lock.Lock()
	state := s.states[id]
	s.lock.Unlock()
	return state
}

func (s *stateStore) set(state *State) {
	s.lock.Lock()
	s.states[state.user.ID] = state
	s.lock.Unlock()
}

func (s *stateStore) delete(id uint64) {
	s.lock.Lock()
	delete(s.states, id)
	for key, session := range s.sessions {
		if session.UserID == id {
			delete(s.sessions, key)
			go s.sessionStore.DeleteSession(key)
		}
	}
	s.lock.Unlock()
}

func (s *stateStore) getSession(key string) *session.Session {
	s.lock.Lock()
	session := s.sessions[key]
	s.lock.Unlock()
	return session
}

func (s *stateStore) setSession(session *session.Session) {
	s.lock.Lock()
	s.sessions[session.Key()] = session
	s.lock.Unlock()
	go s.sessionStore.SaveSession(session)
}

func (s *stateStore) deleteSession(key string) {
	s.lock.Lock()
	s.internalDeleteSession(key)
	s.lock.Unlock()
}

func (s *stateStore) internalDeleteSession(key string) {
	id := s.sessions[key].UserID
	delete(s.sessions, key)

	n := 0
	for _, session := range s.sessions {
		if session.UserID == id {
			n++
		}
	}

	state := s.states[id]
	if n == 0 {
		delete(s.states, id)
	}

	go func() {
		if n == 0 {
			// This anonymous user is not reachable anymore since all sessions have
			// expired, so we clean it up
			state.kill()
			state.user.Remove()
		}

		s.sessionStore.DeleteSession(key)
	}()
}
