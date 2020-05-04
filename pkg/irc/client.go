package irc

import (
	"bufio"
	"crypto/tls"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/jpillora/backoff"
)

type Client struct {
	Server          string
	Host            string
	TLS             bool
	TLSConfig       *tls.Config
	Password        string
	Username        string
	Realname        string
	HandleNickInUse func(string) string

	Messages          chan *Message
	ConnectionChanged chan ConnectionState
	Features          *Features
	nick              string
	channels          []string

	conn       net.Conn
	connected  bool
	registered bool
	dialer     *net.Dialer
	recvBuf    []byte
	scan       *bufio.Scanner
	backoff    *backoff.Backoff
	out        chan string

	quit      chan struct{}
	reconnect chan struct{}
	sendRecv  sync.WaitGroup
	lock      sync.Mutex
}

func NewClient(nick, username string) *Client {
	return &Client{
		nick:              nick,
		Features:          NewFeatures(),
		Username:          username,
		Realname:          nick,
		Messages:          make(chan *Message, 32),
		ConnectionChanged: make(chan ConnectionState, 16),
		out:               make(chan string, 32),
		quit:              make(chan struct{}),
		reconnect:         make(chan struct{}),
		recvBuf:           make([]byte, 0, 4096),
		backoff: &backoff.Backoff{
			Min:    500 * time.Millisecond,
			Max:    30 * time.Second,
			Jitter: true,
		},
	}
}

func (c *Client) GetNick() string {
	c.lock.Lock()
	nick := c.nick
	c.lock.Unlock()
	return nick
}

func (c *Client) setNick(nick string) {
	c.lock.Lock()
	c.nick = nick
	c.lock.Unlock()
}

func (c *Client) Connected() bool {
	c.lock.Lock()
	connected := c.connected
	c.lock.Unlock()
	return connected
}

func (c *Client) Registered() bool {
	c.lock.Lock()
	reg := c.registered
	c.lock.Unlock()
	return reg
}

func (c *Client) setRegistered(reg bool) {
	c.lock.Lock()
	c.registered = reg
	c.lock.Unlock()
}

func (c *Client) Nick(nick string) {
	c.Write("NICK " + nick)
}

func (c *Client) Oper(name, password string) {
	c.Write("OPER " + name + " " + password)
}

func (c *Client) Mode(target, modes, params string) {
	c.Write(strings.TrimRight("MODE "+target+" "+modes+" "+params, " "))
}

func (c *Client) Quit() {
	go func() {
		if c.Connected() {
			c.write("QUIT")
		}
		close(c.quit)
	}()
}

func (c *Client) Join(channels ...string) {
	c.Write("JOIN " + strings.Join(channels, ","))
}

func (c *Client) Part(channels ...string) {
	c.Write("PART " + strings.Join(channels, ","))
	c.removeChannels(channels...)
}

func (c *Client) Topic(channel string, topic ...string) {
	msg := "TOPIC " + channel
	if len(topic) > 0 {
		msg += " :" + topic[0]
	}
	c.Write(msg)
}

func (c *Client) Invite(nick, channel string) {
	c.Write("INVITE " + nick + " " + channel)
}

func (c *Client) Kick(channel string, users ...string) {
	c.Write("KICK " + channel + " " + strings.Join(users, ","))
}

func (c *Client) Privmsg(target, msg string) {
	c.Writef("PRIVMSG %s :%s", target, msg)
}

func (c *Client) Notice(target, msg string) {
	c.Writef("NOTICE %s :%s", target, msg)
}

func (c *Client) Whois(nick string) {
	c.Write("WHOIS " + nick)
}

func (c *Client) Away(message string) {
	c.Write("AWAY :" + message)
}

func (c *Client) List() {
	c.Write("LIST")
}

func (c *Client) writePass(password string) {
	c.write("PASS " + password)
}

func (c *Client) writeNick(nick string) {
	c.write("NICK " + nick)
}

func (c *Client) writeUser(username, realname string) {
	c.writef("USER %s 0 * :%s", username, realname)
}

func (c *Client) register() {
	if c.Password != "" {
		c.writePass(c.Password)
	}
	c.writeNick(c.nick)
	c.writeUser(c.Username, c.Realname)
}

func (c *Client) addChannel(channel string) {
	c.lock.Lock()
	c.channels = append(c.channels, channel)
	c.lock.Unlock()
}

func (c *Client) removeChannels(channels ...string) {
	c.lock.Lock()
	for _, removeCh := range channels {
		for i, ch := range c.channels {
			if c.EqualFold(removeCh, ch) {
				c.channels = append(c.channels[:i], c.channels[i+1:]...)
			}
		}
	}
	c.lock.Unlock()
}

func (c *Client) flushChannels() {
	c.lock.Lock()
	if len(c.channels) > 0 {
		c.Join(c.channels...)
		c.channels = []string{}
	}
	c.lock.Unlock()
}
