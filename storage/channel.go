package storage

import (
	"sync"
)

type ChannelStore struct {
	data map[string]map[string][]string
	lock sync.Mutex
}

func NewChannelStore() *ChannelStore {
	return &ChannelStore{
		data: make(map[string]map[string][]string),
	}
}

func (c *ChannelStore) GetUsers(server, channel string) []string {
	c.lock.Lock()

	users := make([]string, len(c.data[server][channel]))
	copy(users, c.data[server][channel])

	c.lock.Unlock()

	return users
}

func (c *ChannelStore) SetUsers(users []string, server, channel string) {
	c.lock.Lock()

	if _, ok := c.data[server]; !ok {
		c.data[server] = make(map[string][]string)
	}

	c.data[server][channel] = users

	c.lock.Unlock()
}

func (c *ChannelStore) AddUser(user, server, channel string) {
	c.lock.Lock()

	if _, ok := c.data[server]; !ok {
		c.data[server] = make(map[string][]string)
	}

	if users, ok := c.data[server][channel]; ok {
		c.data[server][channel] = append(users, user)
	} else {
		c.data[server][channel] = []string{user}
	}

	c.lock.Unlock()
}

func (c *ChannelStore) RemoveUser(user, server, channel string) {
	c.lock.Lock()
	defer c.lock.Unlock()

	for i, u := range c.data[server][channel] {
		if u == user {
			users := c.data[server][channel]
			c.data[server][channel] = append(users[:i], users[i+1:]...)
			return
		}
	}
}
