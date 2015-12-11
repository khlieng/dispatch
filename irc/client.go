package irc

import (
	"bufio"
	"crypto/tls"
	"net"
	"strings"
	"sync"

	"github.com/khlieng/dispatch/Godeps/_workspace/src/github.com/matryer/resync"
)

type Client struct {
	Server    string
	Host      string
	TLS       bool
	TLSConfig *tls.Config
	Password  string
	Username  string
	Realname  string
	Messages  chan *Message

	nick string

	conn      net.Conn
	connected bool
	dialer    *net.Dialer
	reader    *bufio.Reader
	out       chan string

	quit      chan struct{}
	reconnect chan struct{}
	ready     sync.WaitGroup
	once      resync.Once
	lock      sync.Mutex
}

func NewClient(nick, username string) *Client {
	return &Client{
		nick:      nick,
		Username:  username,
		Realname:  nick,
		Messages:  make(chan *Message, 32),
		out:       make(chan string, 32),
		quit:      make(chan struct{}),
		reconnect: make(chan struct{}),
	}
}

func (c *Client) GetNick() string {
	c.lock.Lock()
	defer c.lock.Unlock()

	return c.nick
}

func (c *Client) Connected() bool {
	c.lock.Lock()
	defer c.lock.Unlock()

	return c.connected
}

func (c *Client) Nick(nick string) {
	c.Write("NICK " + nick)

	c.lock.Lock()
	c.nick = nick
	c.lock.Unlock()
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
		c.lock.Lock()
		c.conn.Close()
		c.lock.Unlock()
	}()
}

func (c *Client) Join(channels ...string) {
	c.Write("JOIN " + strings.Join(channels, ","))
}

func (c *Client) Part(channels ...string) {
	c.Write("PART " + strings.Join(channels, ","))
}

func (c *Client) Topic(channel string) {
	c.Write("TOPIC " + channel)
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
