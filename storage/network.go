package storage

import (
	"fmt"
	"sync"

	"github.com/khlieng/dispatch/pkg/irc"
	"github.com/khlieng/dispatch/version"
)

type Network struct {
	Name           string
	Host           string
	Port           string
	TLS            bool
	ServerPassword string
	Nick           string
	Username       string
	Realname       string
	Account        string
	Password       string

	Features  map[string]interface{}
	Connected bool
	Error     string

	user     *User
	client   *irc.Client
	channels map[string]*Channel
	lock     *sync.Mutex
}

func (n *Network) Save() error {
	return n.user.SaveNetwork(n.Copy())
}

func (n *Network) Copy() *Network {
	n.lock.Lock()
	network := Network{
		Name:           n.Name,
		Host:           n.Host,
		Port:           n.Port,
		TLS:            n.TLS,
		ServerPassword: n.ServerPassword,
		Nick:           n.Nick,
		Username:       n.Username,
		Realname:       n.Realname,
		Account:        n.Account,
		Password:       n.Password,
		Features:       n.Features,
		Connected:      n.Connected,
		Error:          n.Error,
		user:           n.user,
		client:         n.client,
		channels:       n.channels,
		lock:           &sync.Mutex{},
	}
	n.lock.Unlock()

	return &network
}

func (n *Network) Client() *irc.Client {
	return n.client
}

func (n *Network) IRCConfig() *irc.Config {
	return &irc.Config{
		Host:     n.Host,
		Port:     n.Port,
		TLS:      n.TLS,
		Nick:     n.Nick,
		Username: n.Username,
		Realname: n.Realname,
		Account:  n.Account,
		Password: n.Password,
		Version:  fmt.Sprintf("Dispatch %s (git: %s)", version.Tag, version.Commit),
		Source:   "https://github.com/khlieng/dispatch",
	}
}

func (n *Network) SetName(name string) {
	n.lock.Lock()
	n.Name = name
	n.lock.Unlock()
}

func (n *Network) SetNick(nick string) {
	n.lock.Lock()
	n.Nick = nick
	n.lock.Unlock()
}

func (n *Network) SetFeatures(features map[string]interface{}) {
	n.lock.Lock()
	n.Features = features
	n.lock.Unlock()
}

func (n *Network) SetStatus(connected bool, err string) {
	n.lock.Lock()
	n.Connected = connected
	n.Error = err
	n.lock.Unlock()
}

func (n *Network) Channel(name string) *Channel {
	n.lock.Lock()
	ch := n.channels[name]
	n.lock.Unlock()
	return ch
}

func (n *Network) Channels() []*Channel {
	n.lock.Lock()
	channels := make([]*Channel, 0, len(n.channels))
	for _, ch := range n.channels {
		channels = append(channels, ch.Copy())
	}
	n.lock.Unlock()

	return channels
}

func (n *Network) ChannelNames() []string {
	n.lock.Lock()
	names := make([]string, 0, len(n.channels))
	for _, ch := range n.channels {
		names = append(names, ch.Name)
	}
	n.lock.Unlock()

	return names
}

func (n *Network) NewChannel(name string) *Channel {
	return &Channel{
		Network: n.Host,
		Name:    name,
		user:    n.user,
		lock:    &sync.Mutex{},
	}
}

func (n *Network) AddChannel(channel *Channel) {
	n.lock.Lock()
	n.channels[channel.Name] = channel
	n.lock.Unlock()
}

func (n *Network) RemoveChannels(channels ...string) {
	n.lock.Lock()
	for _, name := range channels {
		delete(n.channels, name)
	}
	n.lock.Unlock()
}

type Channel struct {
	Network string
	Name    string

	Topic  string
	Joined bool

	user *User
	lock *sync.Mutex
}

func (c *Channel) Save() error {
	return c.user.SaveChannel(c.Copy())
}

func (c *Channel) Copy() *Channel {
	c.lock.Lock()
	ch := Channel{
		Network: c.Network,
		Name:    c.Name,
		Topic:   c.Topic,
		Joined:  c.Joined,
		user:    c.user,
		lock:    &sync.Mutex{},
	}
	c.lock.Unlock()

	return &ch
}

func (c *Channel) SetTopic(topic string) {
	if c == nil {
		return
	}

	c.lock.Lock()
	c.Topic = topic
	c.lock.Unlock()
}

func (c *Channel) IsJoined() bool {
	if c == nil {
		return false
	}

	c.lock.Lock()
	joined := c.Joined
	c.lock.Unlock()

	return joined
}

func (c *Channel) SetJoined(joined bool) {
	if c == nil {
		return
	}

	c.lock.Lock()
	c.Joined = joined
	c.lock.Unlock()
}
