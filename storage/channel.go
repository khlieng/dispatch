package storage

import (
	"strings"
	"sync"
)

type ChannelStore struct {
	users    map[string]map[string][]*ChannelStoreUser
	userLock sync.Mutex

	topic     map[string]map[string]string
	topicLock sync.Mutex
}

const userModePrefixes = "~&@%+"
const userModeChars = "qaohv"

type ChannelStoreUser struct {
	nick   string
	modes  string
	prefix string
}

func NewChannelStoreUser(nick string) *ChannelStoreUser {
	user := &ChannelStoreUser{nick: nick}

	if i := strings.IndexAny(nick, userModePrefixes); i == 0 {
		i = strings.Index(userModePrefixes, string(nick[0]))
		user.modes = string(userModeChars[i])
		user.nick = nick[1:]
		user.updatePrefix()
	}

	return user
}

func (c *ChannelStoreUser) String() string {
	return c.prefix + c.nick
}

func (c *ChannelStoreUser) addModes(modes string) {
	for _, mode := range modes {
		if strings.Contains(c.modes, string(mode)) {
			continue
		}
		c.modes += string(mode)
	}
	c.updatePrefix()
}

func (c *ChannelStoreUser) removeModes(modes string) {
	for _, mode := range modes {
		c.modes = strings.Replace(c.modes, string(mode), "", 1)
	}
	c.updatePrefix()
}

func (c *ChannelStoreUser) updatePrefix() {
	for i, mode := range userModeChars {
		if strings.Contains(c.modes, string(mode)) {
			c.prefix = string(userModePrefixes[i])
			return
		}
	}
	c.prefix = ""
}

func NewChannelStore() *ChannelStore {
	return &ChannelStore{
		users: make(map[string]map[string][]*ChannelStoreUser),
		topic: make(map[string]map[string]string),
	}
}

func (c *ChannelStore) GetUsers(server, channel string) []string {
	c.userLock.Lock()

	users := make([]string, len(c.users[server][channel]))
	for i, user := range c.users[server][channel] {
		users[i] = user.String()
	}

	c.userLock.Unlock()

	return users
}

func (c *ChannelStore) SetUsers(users []string, server, channel string) {
	c.userLock.Lock()

	if _, ok := c.users[server]; !ok {
		c.users[server] = make(map[string][]*ChannelStoreUser)
	}

	c.users[server][channel] = make([]*ChannelStoreUser, len(users))
	for i, nick := range users {
		c.users[server][channel][i] = NewChannelStoreUser(nick)
	}

	c.userLock.Unlock()
}

func (c *ChannelStore) AddUser(user, server, channel string) {
	c.userLock.Lock()

	if _, ok := c.users[server]; !ok {
		c.users[server] = make(map[string][]*ChannelStoreUser)
	}

	if users, ok := c.users[server][channel]; ok {
		for _, u := range users {
			if u.nick == user {
				c.userLock.Unlock()
				return
			}
		}

		c.users[server][channel] = append(users, NewChannelStoreUser(user))
	} else {
		c.users[server][channel] = []*ChannelStoreUser{NewChannelStoreUser(user)}
	}

	c.userLock.Unlock()
}

func (c *ChannelStore) RemoveUser(user, server, channel string) {
	c.userLock.Lock()
	c.removeUser(user, server, channel)
	c.userLock.Unlock()
}

func (c *ChannelStore) RemoveUserAll(user, server string) {
	c.userLock.Lock()

	for channel := range c.users[server] {
		c.removeUser(user, server, channel)
	}

	c.userLock.Unlock()
}

func (c *ChannelStore) RenameUser(oldNick, newNick, server string) {
	c.userLock.Lock()
	c.renameAll(server, oldNick, newNick)
	c.userLock.Unlock()
}

func (c *ChannelStore) SetMode(server, channel, user, add, remove string) {
	c.userLock.Lock()

	for _, u := range c.users[server][channel] {
		if u.nick == user {
			u.addModes(add)
			u.removeModes(remove)

			c.userLock.Unlock()
			return
		}
	}

	c.userLock.Unlock()
}

func (c *ChannelStore) GetTopic(server, channel string) string {
	c.topicLock.Lock()
	topic := c.topic[server][channel]
	c.topicLock.Unlock()
	return topic
}

func (c *ChannelStore) SetTopic(topic, server, channel string) {
	c.topicLock.Lock()

	if _, ok := c.topic[server]; !ok {
		c.topic[server] = make(map[string]string)
	}

	c.topic[server][channel] = topic
	c.topicLock.Unlock()
}

func (c *ChannelStore) rename(server, channel, oldNick, newNick string) {
	for _, user := range c.users[server][channel] {
		if user.nick == oldNick {
			user.nick = newNick
			return
		}
	}
}

func (c *ChannelStore) renameAll(server, oldNick, newNick string) {
	for channel := range c.users[server] {
		c.rename(server, channel, oldNick, newNick)
	}
}

func (c *ChannelStore) removeUser(user, server, channel string) {
	for i, u := range c.users[server][channel] {
		if u.nick == user {
			users := c.users[server][channel]
			c.users[server][channel] = append(users[:i], users[i+1:]...)
			return
		}
	}
}
