package irc

import (
	"strings"
	"time"
)

func (c *Client) handleMessage(msg *Message) {
	switch msg.Command {
	case CAP:
		c.handleCAP(msg)

	case PING:
		go c.write("PONG :" + msg.LastParam())

	case JOIN:
		if len(msg.Params) > 0 {
			channel := msg.Params[0]

			if c.Is(msg.Sender) {
				c.addChannel(channel)
			}

			c.state.addUser(msg.Sender, channel)
		}

	case PART:
		if len(msg.Params) > 0 {
			channel := msg.Params[0]

			if c.Is(msg.Sender) {
				c.state.removeChannel(channel)
			} else {
				c.state.removeUser(msg.Sender, channel)
			}
		}

	case QUIT:
		msg.meta = c.state.removeUserAll(msg.Sender)

	case KICK:
		if len(msg.Params) > 1 {
			channel, nick := msg.Params[0], msg.Params[1]

			if c.Is(nick) {
				c.removeChannels(channel)
				c.state.removeChannel(channel)
			} else {
				c.state.removeUser(nick, channel)
			}
		}

	case NICK:
		if c.Is(msg.Sender) {
			c.setNick(msg.LastParam())
		}

		msg.meta = c.state.renameUser(msg.Sender, msg.LastParam())

	case PRIVMSG:
		if ctcp := msg.ToCTCP(); ctcp != nil {
			c.handleCTCP(ctcp, msg)
		}

	case MODE:
		if len(msg.Params) > 1 {
			target := msg.Params[0]
			if len(msg.Params) > 2 && isChannel(target) {
				mode := ParseMode(msg.Params[1])
				mode.Network = c.Host()
				mode.Channel = target
				mode.User = msg.Params[2]

				c.state.setMode(target, msg.Params[2], mode.Add, mode.Remove)

				msg.meta = mode
			}
		}

	case TOPIC, RPL_TOPIC:
		chIndex := 0
		if msg.Command == RPL_TOPIC {
			chIndex = 1
		}

		if len(msg.Params) > chIndex {
			c.state.setTopic(msg.LastParam(), msg.Params[chIndex])
		}

	case RPL_NOTOPIC:
		if len(msg.Params) > 1 {
			channel := msg.Params[1]
			c.state.setTopic("", channel)
		}

	case RPL_WELCOME:
		if len(msg.Params) > 0 {
			c.setNick(msg.Params[0])
		}
		c.negotiating = false
		c.setRegistered(true)
		c.flushChannels()

		c.backoff.Reset()
		c.sendRecv.Add(1)
		go c.send()

	case RPL_ISUPPORT:
		c.Features.Parse(msg.Params)

	case ERR_NICKNAMEINUSE, ERR_NICKCOLLISION, ERR_UNAVAILRESOURCE:
		if c.Config.HandleNickInUse != nil && len(msg.Params) > 1 {
			go c.writeNick(c.Config.HandleNickInUse(msg.Params[1]))
		}

	case RPL_NAMREPLY:
		if len(msg.Params) > 2 {
			channel := msg.Params[2]
			users := strings.Split(strings.TrimSuffix(msg.LastParam(), " "), " ")

			userBuffer := c.state.userBuffers[channel]
			c.state.userBuffers[channel] = append(userBuffer, users...)
		}

	case RPL_ENDOFNAMES:
		if len(msg.Params) > 1 {
			channel := msg.Params[1]
			users := c.state.userBuffers[channel]

			c.state.setUsers(users, channel)
			delete(c.state.userBuffers, channel)
			msg.meta = users
		}

	case ERROR:
		c.Messages <- msg
		c.connChange(false, nil)
		time.Sleep(5 * time.Second)
		close(c.quit)
		return
	}

	c.handleSASL(msg)
}

type Mode struct {
	Network string
	Channel string
	User    string
	Add     string
	Remove  string
}

func ParseMode(mode string) *Mode {
	m := Mode{}
	add := false

	for _, c := range mode {
		if c == '+' {
			add = true
		} else if c == '-' {
			add = false
		} else if add {
			m.Add += string(c)
		} else {
			m.Remove += string(c)
		}
	}

	return &m
}

func isChannel(s string) bool {
	return strings.IndexAny(s, "&#+!") == 0
}
