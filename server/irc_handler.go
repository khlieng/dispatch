package server

import (
	"log"
	"strings"

	"github.com/khlieng/dispatch/irc"
	"github.com/khlieng/dispatch/storage"
)

type ircHandler struct {
	client  *irc.Client
	session *Session

	whois       WhoisReply
	userBuffers map[string][]string
	motdBuffer  MOTD

	handlers map[string]func(*irc.Message)
}

func newIRCHandler(client *irc.Client, session *Session) *ircHandler {
	i := &ircHandler{
		client:      client,
		session:     session,
		userBuffers: make(map[string][]string),
	}
	i.initHandlers()
	return i
}

func (i *ircHandler) run() {
	for {
		select {
		case msg, ok := <-i.client.Messages:
			if !ok {
				i.session.deleteIRC(i.client.Host)
				return
			}

			i.dispatchMessage(msg)

		case connected := <-i.client.ConnectionChanged:
			i.session.sendJSON("connection_update", map[string]bool{
				i.client.Host: connected,
			})
			i.session.setConnectionState(i.client.Host, connected)
		}
	}
}

func (i *ircHandler) dispatchMessage(msg *irc.Message) {
	if handler, ok := i.handlers[msg.Command]; ok {
		handler(msg)
	}
}

func (i *ircHandler) nick(msg *irc.Message) {
	i.session.sendJSON("nick", Nick{
		Server:   i.client.Host,
		Old:      msg.Nick,
		New:      msg.Trailing,
		Channels: channelStore.FindUserChannels(msg.Nick, i.client.Host),
	})

	channelStore.RenameUser(msg.Nick, msg.Trailing, i.client.Host)
}

func (i *ircHandler) join(msg *irc.Message) {
	i.session.sendJSON("join", Join{
		Server:   i.client.Host,
		User:     msg.Nick,
		Channels: msg.Params,
	})

	channelStore.AddUser(msg.Nick, i.client.Host, msg.Params[0])

	if msg.Nick == i.client.GetNick() {
		i.session.user.AddChannel(storage.Channel{
			Server: i.client.Host,
			Name:   msg.Params[0],
		})
	}
}

func (i *ircHandler) part(msg *irc.Message) {
	i.session.sendJSON("part", Part{
		Join: Join{
			Server:   i.client.Host,
			User:     msg.Nick,
			Channels: msg.Params,
		},
		Reason: msg.Trailing,
	})

	channelStore.RemoveUser(msg.Nick, i.client.Host, msg.Params[0])

	if msg.Nick == i.client.GetNick() {
		i.session.user.RemoveChannel(i.client.Host, msg.Params[0])
	}
}

func (i *ircHandler) mode(msg *irc.Message) {
	target := msg.Params[0]
	if len(msg.Params) > 2 && isChannel(target) {
		mode := parseMode(msg.Params[1])
		mode.Server = i.client.Host
		mode.Channel = target
		mode.User = msg.Params[2]

		i.session.sendJSON("mode", mode)

		channelStore.SetMode(i.client.Host, target, msg.Params[2], mode.Add, mode.Remove)
	}
}

func (i *ircHandler) message(msg *irc.Message) {
	message := Chat{
		Server:  i.client.Host,
		From:    msg.Nick,
		Message: msg.Trailing,
	}

	if msg.Params[0] == i.client.GetNick() {
		i.session.sendJSON("pm", message)
	} else {
		message.To = msg.Params[0]
		i.session.sendJSON("message", message)
	}

	if msg.Params[0] != "*" {
		go i.session.user.LogMessage(i.client.Host, msg.Nick, msg.Params[0], msg.Trailing)
	}
}

func (i *ircHandler) quit(msg *irc.Message) {
	i.session.sendJSON("quit", Quit{
		Server:   i.client.Host,
		User:     msg.Nick,
		Reason:   msg.Trailing,
		Channels: channelStore.FindUserChannels(msg.Nick, i.client.Host),
	})

	channelStore.RemoveUserAll(msg.Nick, i.client.Host)
}

func (i *ircHandler) info(msg *irc.Message) {
	i.session.sendJSON("pm", Chat{
		Server:  i.client.Host,
		From:    msg.Nick,
		Message: strings.Join(msg.Params[1:], " "),
	})
}

func (i *ircHandler) whoisUser(msg *irc.Message) {
	i.whois.Nick = msg.Params[1]
	i.whois.Username = msg.Params[2]
	i.whois.Host = msg.Params[3]
	i.whois.Realname = msg.Params[5]
}

func (i *ircHandler) whoisServer(msg *irc.Message) {
	i.whois.Server = msg.Params[2]
}

func (i *ircHandler) whoisChannels(msg *irc.Message) {
	i.whois.Channels = append(i.whois.Channels, strings.Split(strings.TrimRight(msg.Trailing, " "), " ")...)
}

func (i *ircHandler) whoisEnd(msg *irc.Message) {
	i.session.sendJSON("whois", i.whois)
	i.whois = WhoisReply{}
}

func (i *ircHandler) topic(msg *irc.Message) {
	i.session.sendJSON("topic", Topic{
		Server:  i.client.Host,
		Channel: msg.Params[1],
		Topic:   msg.Trailing,
	})

	channelStore.SetTopic(msg.Trailing, i.client.Host, msg.Params[1])
}

func (i *ircHandler) names(msg *irc.Message) {
	users := strings.Split(msg.Trailing, " ")
	userBuffer := i.userBuffers[msg.Params[2]]
	i.userBuffers[msg.Params[2]] = append(userBuffer, users...)
}

func (i *ircHandler) namesEnd(msg *irc.Message) {
	channel := msg.Params[1]
	users := i.userBuffers[channel]

	i.session.sendJSON("users", Userlist{
		Server:  i.client.Host,
		Channel: channel,
		Users:   users,
	})

	channelStore.SetUsers(users, i.client.Host, channel)
	delete(i.userBuffers, channel)
}

func (i *ircHandler) motdStart(msg *irc.Message) {
	i.motdBuffer.Server = i.client.Host
	i.motdBuffer.Title = msg.Trailing
}

func (i *ircHandler) motd(msg *irc.Message) {
	i.motdBuffer.Content = append(i.motdBuffer.Content, msg.Trailing)
}

func (i *ircHandler) motdEnd(msg *irc.Message) {
	i.session.sendJSON("motd", i.motdBuffer)
	i.motdBuffer = MOTD{}
}

func (i *ircHandler) initHandlers() {
	i.handlers = map[string]func(*irc.Message){
		irc.Nick:               i.nick,
		irc.Join:               i.join,
		irc.Part:               i.part,
		irc.Mode:               i.mode,
		irc.Privmsg:            i.message,
		irc.Notice:             i.message,
		irc.Quit:               i.quit,
		irc.ReplyWelcome:       i.info,
		irc.ReplyYourHost:      i.info,
		irc.ReplyCreated:       i.info,
		irc.ReplyLUserClient:   i.info,
		irc.ReplyLUserOp:       i.info,
		irc.ReplyLUserUnknown:  i.info,
		irc.ReplyLUserChannels: i.info,
		irc.ReplyLUserMe:       i.info,
		irc.ReplyWhoisUser:     i.whoisUser,
		irc.ReplyWhoisServer:   i.whoisServer,
		irc.ReplyWhoisChannels: i.whoisChannels,
		irc.ReplyEndOfWhois:    i.whoisEnd,
		irc.ReplyTopic:         i.topic,
		irc.ReplyNamReply:      i.names,
		irc.ReplyEndOfNames:    i.namesEnd,
		irc.ReplyMotdStart:     i.motdStart,
		irc.ReplyMotd:          i.motd,
		irc.ReplyEndOfMotd:     i.motdEnd,
	}
}

func parseMode(mode string) *Mode {
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

func printMessage(msg *irc.Message, i *irc.Client) {
	log.Println(i.GetNick()+":", msg.Prefix, msg.Command, msg.Params, msg.Trailing)
}
