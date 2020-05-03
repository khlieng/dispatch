package irc

import (
	"bufio"
	"crypto/tls"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"net"
	"os"
	"path/filepath"
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

	DownloadFolder string
	Autoget        bool

	Messages          chan *Message
	ConnectionChanged chan ConnectionState
	Progress          chan DownloadProgress
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
		Progress:          make(chan DownloadProgress, 16),
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

func byteRead(totalBytes uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, totalBytes)
	return b
}

func round2(source float64) float64 {
	return math.Round(100*source) / 100
}

const (
	_ = 1.0 << (10 * iota)
	kibibyte
	mebibyte
	gibibyte
)

func humanReadableByteCount(b float64, speed bool) string {
	unit := ""
	value := b

	switch {
	case b >= gibibyte:
		unit = "GiB"
		value = value / gibibyte
	case b >= mebibyte:
		unit = "MiB"
		value = value / mebibyte
	case b >= kibibyte:
		unit = "KiB"
		value = value / kibibyte
	case b > 1 || b == 0:
		unit = "bytes"
	case b == 1:
		unit = "byte"
	}

	if speed {
		unit = unit + "/s"
	}

	stringValue := strings.TrimSuffix(
		fmt.Sprintf("%.2f", value), ".00",
	)

	return fmt.Sprintf("%s %s", stringValue, unit)
}

func (c *Client) Download(pack *CTCP) {
	if !c.Autoget {
		// TODO: ask user if he/she wants to download the file
		return
	}
	c.Progress <- DownloadProgress{
		PercCompletion: 0,
		File:           pack.File,
	}
	file, err := os.OpenFile(filepath.Join(c.DownloadFolder, pack.File), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		c.Progress <- DownloadProgress{
			PercCompletion: -1,
			File:           pack.File,
			Error:          err,
		}
		return
	}
	defer file.Close()

	con, err := net.Dial("tcp", fmt.Sprintf("%s:%d", pack.IP, pack.Port))

	if err != nil {
		c.Progress <- DownloadProgress{
			PercCompletion: -1,
			File:           pack.File,
			Error:          err,
		}
		return
	}

	defer con.Close()

	var avgSpeed float64
	var prevTime int64 = -1
	secondsElapsed := int64(0)
	totalBytes := uint64(0)
	buf := make([]byte, 0, 4*1024)
	start := time.Now().UnixNano()
	for {
		n, err := con.Read(buf[:cap(buf)])
		buf = buf[:n]
		if n == 0 {
			if err == nil {
				continue
			}
			if err == io.EOF {
				break
			}
		}

		if _, err := file.Write(buf); err != nil {
			c.Progress <- DownloadProgress{
				PercCompletion: -1,
				File:           pack.File,
				Error:          err,
			}
			return
		}

		cycleBytes := uint64(len(buf))
		totalBytes += cycleBytes
		percentage := round2(100 * float64(totalBytes) / float64(pack.Length))

		now := time.Now().UnixNano()
		secondsElapsed = (now - start) / 1e9
		avgSpeed = round2(float64(totalBytes) / (float64(secondsElapsed)))
		speed := 0.0

		if prevTime < 0 {
			speed = avgSpeed
		} else {
			speed = round2(1e9 * float64(cycleBytes) / (float64(now - prevTime)))
		}
		secondsToGo := (float64(pack.Length) - float64(totalBytes)) / speed
		prevTime = now
		con.Write(byteRead(totalBytes))
		c.Progress <- DownloadProgress{
			InstSpeed:      humanReadableByteCount(speed, true),
			AvgSpeed:       humanReadableByteCount(avgSpeed, true),
			PercCompletion: percentage,
			BytesRemaining: humanReadableByteCount(float64(pack.Length-totalBytes), false),
			BytesCompleted: humanReadableByteCount(float64(totalBytes), false),
			SecondsElapsed: secondsElapsed,
			SecondsToGo:    secondsToGo,
			File:           pack.File,
		}
	}
	con.Write(byteRead(totalBytes))
	c.Progress <- DownloadProgress{
		AvgSpeed:       humanReadableByteCount(avgSpeed, true),
		PercCompletion: 100,
		BytesCompleted: humanReadableByteCount(float64(totalBytes), false),
		SecondsElapsed: secondsElapsed,
		File:           pack.File,
	}
}
