package storage

import (
	"sync"
)

type ChannelStore struct {
	users    map[string]map[string][]string
	userLock sync.Mutex

	topic     map[string]map[string]string
	topicLock sync.Mutex
}

func NewChannelStore() *ChannelStore {
	return &ChannelStore{
		users: make(map[string]map[string][]string),
		topic: make(map[string]map[string]string),
	}
}

func (c *ChannelStore) GetUsers(server, channel string) []string {
	c.userLock.Lock()

	users := make([]string, len(c.users[server][channel]))
	copy(users, c.users[server][channel])

	c.userLock.Unlock()

	return users
}

func (c *ChannelStore) SetUsers(users []string, server, channel string) {
	c.userLock.Lock()

	if _, ok := c.users[server]; !ok {
		c.users[server] = make(map[string][]string)
	}

	c.users[server][channel] = users

	c.userLock.Unlock()
}

func (c *ChannelStore) AddUser(user, server, channel string) {
	c.userLock.Lock()

	if _, ok := c.users[server]; !ok {
		c.users[server] = make(map[string][]string)
	}

	if users, ok := c.users[server][channel]; ok {
		c.users[server][channel] = append(users, user)
	} else {
		c.users[server][channel] = []string{user}
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

	for channel, _ := range c.users[server] {
		c.removeUser(user, server, channel)
	}

	c.userLock.Unlock()
}

func (c *ChannelStore) GetTopic(server, channel string) string {
	c.topicLock.Lock()
	defer c.topicLock.Unlock()

	return c.topic[server][channel]
}

func (c *ChannelStore) SetTopic(topic, server, channel string) {
	c.topicLock.Lock()

	if _, ok := c.topic[server]; !ok {
		c.topic[server] = make(map[string]string)
	}

	c.topic[server][channel] = topic
	c.topicLock.Unlock()
}

func (c *ChannelStore) removeUser(user, server, channel string) {
	for i, u := range c.users[server][channel] {
		if u == user {
			users := c.users[server][channel]
			c.users[server][channel] = append(users[:i], users[i+1:]...)
			return
		}
	}
}
