package irc

import (
	"strings"
	"sync"
)

type state struct {
	client *Client

	users map[string][]*User
	topic map[string]string

	userBuffers map[string][]string

	motd []string

	lock sync.Mutex
}

const userModePrefixes = "~&@%+"
const userModeChars = "qaohv"

type User struct {
	nick   string
	modes  string
	prefix string
}

func NewUser(nick string) *User {
	user := &User{nick: nick}

	if i := strings.IndexAny(nick, userModePrefixes); i == 0 {
		i = strings.Index(userModePrefixes, string(nick[0]))
		user.modes = string(userModeChars[i])
		user.nick = nick[1:]
		user.updatePrefix()
	}

	return user
}

func (u *User) String() string {
	return u.prefix + u.nick
}

func (u *User) AddModes(modes string) {
	for _, mode := range modes {
		if strings.Contains(u.modes, string(mode)) {
			continue
		}
		u.modes += string(mode)
	}
	u.updatePrefix()
}

func (u *User) RemoveModes(modes string) {
	for _, mode := range modes {
		u.modes = strings.Replace(u.modes, string(mode), "", 1)
	}
	u.updatePrefix()
}

func (u *User) updatePrefix() {
	for i, mode := range userModeChars {
		if strings.Contains(u.modes, string(mode)) {
			u.prefix = string(userModePrefixes[i])
			return
		}
	}
	u.prefix = ""
}

func newState(client *Client) *state {
	return &state{
		client:      client,
		users:       make(map[string][]*User),
		topic:       make(map[string]string),
		userBuffers: make(map[string][]string),
	}
}

func (s *state) reset() {
	s.lock.Lock()
	s.users = make(map[string][]*User)
	s.topic = make(map[string]string)
	s.userBuffers = make(map[string][]string)
	s.motd = []string{}
	s.lock.Unlock()
}

func (s *state) removeChannel(channel string) {
	s.lock.Lock()
	delete(s.users, channel)
	delete(s.topic, channel)
	s.lock.Unlock()
}

func (s *state) getUsers(channel string) []string {
	s.lock.Lock()

	users := make([]string, len(s.users[channel]))
	for i, user := range s.users[channel] {
		users[i] = user.String()
	}

	s.lock.Unlock()

	return users
}

func (s *state) setUsers(users []string, channel string) {
	s.lock.Lock()

	s.users[channel] = make([]*User, len(users))
	for i, nick := range users {
		s.users[channel][i] = NewUser(nick)
	}

	s.lock.Unlock()
}

func (s *state) addUser(user, channel string) {
	s.lock.Lock()

	if users, ok := s.users[channel]; ok {
		for _, u := range users {
			if u.nick == user {
				s.lock.Unlock()
				return
			}
		}

		s.users[channel] = append(users, NewUser(user))
	} else {
		s.users[channel] = []*User{NewUser(user)}
	}

	s.lock.Unlock()
}

func (s *state) removeUser(user, channel string) {
	s.lock.Lock()
	s.internalRemoveUser(user, channel)
	s.lock.Unlock()
}

func (s *state) removeUserAll(user string) []string {
	channels := []string{}
	s.lock.Lock()

	for channel := range s.users {
		if s.internalRemoveUser(user, channel) {
			channels = append(channels, channel)
		}
	}

	s.lock.Unlock()
	return channels
}

func (s *state) renameUser(oldNick, newNick string) []string {
	s.lock.Lock()
	channels := s.renameAll(oldNick, newNick)
	s.lock.Unlock()
	return channels
}

func (s *state) setMode(channel, user, add, remove string) {
	s.lock.Lock()

	for _, u := range s.users[channel] {
		if u.nick == user {
			u.AddModes(add)
			u.RemoveModes(remove)

			s.lock.Unlock()
			return
		}
	}

	s.lock.Unlock()
}

func (s *state) getTopic(channel string) string {
	s.lock.Lock()
	topic := s.topic[channel]
	s.lock.Unlock()
	return topic
}

func (s *state) setTopic(topic, channel string) {
	s.lock.Lock()
	s.topic[channel] = topic
	s.lock.Unlock()
}

func (s *state) getMOTD() []string {
	s.lock.Lock()
	motd := s.motd
	s.lock.Unlock()
	return motd
}

func (s *state) rename(channel, oldNick, newNick string) bool {
	for _, user := range s.users[channel] {
		if user.nick == oldNick {
			user.nick = newNick
			return true
		}
	}
	return false
}

func (s *state) renameAll(oldNick, newNick string) []string {
	channels := []string{}

	for channel := range s.users {
		if s.rename(channel, oldNick, newNick) {
			channels = append(channels, channel)
		}
	}

	return channels
}

func (s *state) internalRemoveUser(user, channel string) bool {
	for i, u := range s.users[channel] {
		if u.nick == user {
			users := s.users[channel]
			s.users[channel] = append(users[:i], users[i+1:]...)
			return true
		}
	}
	return false
}
