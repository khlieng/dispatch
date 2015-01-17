package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"net"
	"strings"
)

const (
	PING    = "PING"
	JOIN    = "JOIN"
	PART    = "PART"
	MODE    = "MODE"
	PRIVMSG = "PRIVMSG"
	NOTICE  = "NOTICE"
	TOPIC   = "TOPIC"
	QUIT    = "QUIT"

	RPL_WELCOME       = "001"
	RPL_YOURHOST      = "002"
	RPL_CREATED       = "003"
	RPL_LUSERCLIENT   = "251"
	RPL_LUSEROP       = "252"
	RPL_LUSERUNKNOWN  = "253"
	RPL_LUSERCHANNELS = "254"
	RPL_LUSERME       = "255"

	RPL_WHOISUSER     = "311"
	RPL_WHOISSERVER   = "312"
	RPL_WHOISOPERATOR = "313"
	RPL_WHOISIDLE     = "317"
	RPL_ENDOFWHOIS    = "318"
	RPL_WHOISCHANNELS = "319"

	RPL_TOPIC = "332"

	RPL_NAMREPLY   = "353"
	RPL_ENDOFNAMES = "366"

	RPL_MOTD      = "372"
	RPL_MOTDSTART = "375"
	RPL_ENDOFMOTD = "376"
)

type Message struct {
	Prefix   string
	Command  string
	Params   []string
	Trailing string
}

type IRC struct {
	conn   net.Conn
	reader *bufio.Reader

	Messages  chan *Message
	Server    string
	Host      string
	TLS       bool
	TLSConfig *tls.Config
	nick      string
	Username  string
	Realname  string
}

func NewIRC(nick, username string) *IRC {
	return &IRC{
		nick:     nick,
		Username: username,
		Realname: nick,
		Messages: make(chan *Message, 32),
	}
}

func (i *IRC) Connect(address string) {
	if idx := strings.Index(address, ":"); idx < 0 {
		i.Host = address

		if i.TLS {
			address += ":6697"
		} else {
			address += ":6667"
		}
	} else {
		i.Host = address[:idx]
	}
	i.Server = address

	if i.TLS {
		if i.TLSConfig == nil {
			i.TLSConfig = &tls.Config{InsecureSkipVerify: true}
		}
		i.conn, _ = tls.Dial("tcp", address, i.TLSConfig)
	} else {
		i.conn, _ = net.Dial("tcp", address)
	}

	i.reader = bufio.NewReader(i.conn)

	i.Nick(i.nick)
	i.User(i.Username, i.Realname)

	go i.recv()
}

func (i *IRC) Pass(password string) {
	i.Write("PASS " + password)
}

func (i *IRC) Nick(nick string) {
	i.Write("NICK " + nick)
}

func (i *IRC) User(username, realname string) {
	i.Writef("USER %s 0 * :%s", username, realname)
}

func (i *IRC) Join(channels ...string) {
	i.Write("JOIN " + strings.Join(channels, ","))
}

func (i *IRC) Part(channels ...string) {
	i.Write("PART " + strings.Join(channels, ","))
}

func (i *IRC) Privmsg(target, msg string) {
	i.Writef("PRIVMSG %s :%s", target, msg)
}

func (i *IRC) Notice(target, msg string) {
	i.Writef("NOTICE %s :%s", target, msg)
}

func (i *IRC) Topic(channel string) {
	i.Write("TOPIC " + channel)
}

func (i *IRC) Whois(nick string) {
	i.Write("WHOIS " + nick)
}

func (i *IRC) Quit() {
	i.Write("QUIT")
	i.conn.Close()
}

func (i *IRC) Write(data string) {
	fmt.Fprint(i.conn, data+"\r\n")
}

func (i *IRC) Writef(format string, a ...interface{}) {
	fmt.Fprintf(i.conn, format+"\r\n", a...)
}

func (i *IRC) recv() {
	defer i.conn.Close()
	for {
		line, err := i.reader.ReadString('\n')
		if err != nil {
			return
		}

		msg := parseMessage(line)
		msg.Prefix = parseUser(msg.Prefix)

		switch msg.Command {
		case PING:
			i.Write("PONG :" + msg.Trailing)
		}

		i.Messages <- msg
	}
}

func parseMessage(line string) *Message {
	line = strings.Trim(line, "\r\n")
	msg := Message{}
	cmdStart := 0
	cmdEnd := len(line)

	if strings.HasPrefix(line, ":") {
		cmdStart = strings.Index(line, " ") + 1
		msg.Prefix = line[1 : cmdStart-1]
	}

	if i := strings.LastIndex(line, " :"); i > 0 {
		cmdEnd = i
		msg.Trailing = line[i+2:]
	}

	cmd := strings.Split(line[cmdStart:cmdEnd], " ")
	msg.Command = cmd[0]
	if len(cmd) > 1 {
		msg.Params = cmd[1:]
	}

	return &msg
}

func parseUser(user string) string {
	if i := strings.Index(user, "!"); i > 0 {
		return user[:i]
	}
	return user
}
