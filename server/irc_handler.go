package server

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/kjk/betterguid"

	"github.com/khlieng/dispatch/pkg/irc"
	"github.com/khlieng/dispatch/storage"
)

var excludedErrors = []string{
	irc.ERR_NICKNAMEINUSE,
	irc.ERR_NICKCOLLISION,
	irc.ERR_UNAVAILRESOURCE,
	irc.ERR_FORWARD,
}

type ircHandler struct {
	client *irc.Client
	state  *State

	whois       WhoisReply
	userBuffers map[string][]string
	motdBuffer  MOTD
	listBuffer  storage.ChannelListIndex
	dccProgress chan irc.DownloadProgress

	handlers map[string]func(*irc.Message)
}

func newIRCHandler(client *irc.Client, state *State) *ircHandler {
	i := &ircHandler{
		client:      client,
		state:       state,
		userBuffers: make(map[string][]string),
		dccProgress: make(chan irc.DownloadProgress, 4),
	}
	i.initHandlers()
	return i
}

func (i *ircHandler) run() {
	var lastConnErr error
	for {
		select {
		case msg, ok := <-i.client.Messages:
			if !ok {
				i.state.deleteIRC(i.client.Host)
				return
			}

			i.dispatchMessage(msg)

		case state := <-i.client.ConnectionChanged:
			i.state.sendJSON("connection_update", newConnectionUpdate(i.client.Host, state))
			i.state.setConnectionState(i.client.Host, state)

			if state.Error != nil && (lastConnErr == nil ||
				state.Error.Error() != lastConnErr.Error()) {
				lastConnErr = state.Error
				i.log("Connection error:", state.Error)
			} else if state.Connected {
				i.log("Connected")
			}

		case progress := <-i.dccProgress:
			if progress.Error != nil {
				i.sendDCCInfo("%s: Download failed (%s)", true, progress.File, progress.Error)
			} else if progress.PercCompletion == 100 {
				i.sendDCCInfo("Download finished, get it here: %s://%s/downloads/%s/%s", true,
					i.state.String("scheme"), i.state.String("host"), i.state.user.Username, progress.File)
			} else if progress.PercCompletion == 0 {
				i.sendDCCInfo("%s: Starting download", true, progress.File)
			} else {
				i.sendDCCInfo("%s: %.1f%%, %s, %s remaining, %.1fs left", false, progress.File,
					progress.PercCompletion, progress.Speed, progress.BytesRemaining, progress.SecondsToGo)
			}
		}
	}
}

func (i *ircHandler) dispatchMessage(msg *irc.Message) {
	if msg.Command[0] == '4' && !isExcludedError(msg.Command) {
		err := IRCError{
			Server:  i.client.Host,
			Message: msg.LastParam(),
		}

		if len(msg.Params) > 2 {
			for i := 1; i < len(msg.Params); i++ {
				if isChannel(msg.Params[i]) {
					err.Target = msg.Params[i]
					break
				}
			}
		}

		i.state.sendJSON("error", err)
	}

	if handler, ok := i.handlers[msg.Command]; ok {
		handler(msg)
	}
}

func (i *ircHandler) nick(msg *irc.Message) {
	i.state.sendJSON("nick", Nick{
		Server: i.client.Host,
		Old:    msg.Nick,
		New:    msg.LastParam(),
	})

	channelStore.RenameUser(msg.Nick, msg.LastParam(), i.client.Host)

	if i.client.Is(msg.LastParam()) {
		go i.state.user.SetNick(msg.LastParam(), i.client.Host)
	}
}

func (i *ircHandler) join(msg *irc.Message) {
	i.state.sendJSON("join", Join{
		Server:   i.client.Host,
		User:     msg.Nick,
		Channels: msg.Params,
	})

	channel := msg.Params[0]
	channelStore.AddUser(msg.Nick, i.client.Host, channel)

	if i.client.Is(msg.Nick) {
		// In case no topic is set and there's a cached one that needs to be cleared
		i.client.Topic(channel)

		i.state.sendLastMessages(i.client.Host, channel, 50)

		go i.state.user.AddChannel(&storage.Channel{
			Server: i.client.Host,
			Name:   channel,
		})
	}
}

func (i *ircHandler) part(msg *irc.Message) {
	part := Part{
		Server:  i.client.Host,
		User:    msg.Nick,
		Channel: msg.Params[0],
	}

	if len(msg.Params) == 2 {
		part.Reason = msg.Params[1]
	}

	i.state.sendJSON("part", part)

	channelStore.RemoveUser(msg.Nick, i.client.Host, part.Channel)

	if i.client.Is(msg.Nick) {
		go i.state.user.RemoveChannel(i.client.Host, part.Channel)
	}
}

func (i *ircHandler) mode(msg *irc.Message) {
	target := msg.Params[0]
	if len(msg.Params) > 2 && isChannel(target) {
		mode := parseMode(msg.Params[1])
		mode.Server = i.client.Host
		mode.Channel = target
		mode.User = msg.Params[2]

		i.state.sendJSON("mode", mode)

		channelStore.SetMode(i.client.Host, target, msg.Params[2], mode.Add, mode.Remove)
	}
}

func (i *ircHandler) message(msg *irc.Message) {
	if ctcp := msg.ToCTCP(); ctcp != nil {
		if ctcp.Command == "DCC" && strings.HasPrefix(ctcp.Params, "SEND") {
			if pack := irc.ParseDCCSend(ctcp); pack != nil {
				go i.receiveDCCSend(pack, msg)
				return
			}
		} else if ctcp.Command != "ACTION" {
			return
		}
	}

	message := Message{
		ID:      betterguid.New(),
		Server:  i.client.Host,
		From:    msg.Nick,
		Content: msg.LastParam(),
	}
	target := msg.Params[0]

	if i.client.Is(target) {
		i.state.sendJSON("pm", message)

		if !msg.IsFromServer() {
			i.state.user.AddOpenDM(i.client.Host, message.From)
		}

		target = message.From
	} else {
		message.To = target
		i.state.sendJSON("message", message)
	}

	if target != "*" && !msg.IsFromServer() {
		go i.state.user.LogMessage(message.ID,
			i.client.Host, msg.Nick, target, msg.LastParam())
	}
}

func (i *ircHandler) quit(msg *irc.Message) {
	i.state.sendJSON("quit", Quit{
		Server: i.client.Host,
		User:   msg.Nick,
		Reason: msg.LastParam(),
	})

	channelStore.RemoveUserAll(msg.Nick, i.client.Host)
}

func (i *ircHandler) info(msg *irc.Message) {
	if msg.Command == irc.RPL_WELCOME {
		i.state.sendJSON("nick", Nick{
			Server: i.client.Host,
			New:    msg.Params[0],
		})

		_, needsUpdate := channelIndexes.Get(i.client.Host)
		if needsUpdate {
			i.listBuffer = storage.NewMapChannelListIndex()
			i.client.List()
		}

		go i.state.user.SetNick(msg.Params[0], i.client.Host)
	}

	i.state.sendJSON("pm", Message{
		Server:  i.client.Host,
		From:    msg.Nick,
		Content: strings.Join(msg.Params[1:], " "),
	})
}

func (i *ircHandler) features(msg *irc.Message) {
	i.state.sendJSON("features", Features{
		Server:   i.client.Host,
		Features: i.client.Features.Map(),
	})

	if name := i.client.Features.String("NETWORK"); name != "" {
		go func() {
			server, err := i.state.user.GetServer(i.client.Host)
			if err == nil && server.Name == "" {
				i.state.user.SetServerName(name, server.Host)
			}
		}()
	}
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
	i.whois.Channels = append(i.whois.Channels, strings.Split(strings.TrimRight(msg.LastParam(), " "), " ")...)
}

func (i *ircHandler) whoisEnd(msg *irc.Message) {
	if i.whois.Nick != "" {
		i.state.sendJSON("whois", i.whois)
	}
	i.whois = WhoisReply{}
}

func (i *ircHandler) topic(msg *irc.Message) {
	var channel string
	var nick string

	if msg.Command == irc.TOPIC {
		channel = msg.Params[0]
		nick = msg.Nick
	} else {
		channel = msg.Params[1]
	}

	i.state.sendJSON("topic", Topic{
		Server:  i.client.Host,
		Channel: channel,
		Topic:   msg.LastParam(),
		Nick:    nick,
	})

	channelStore.SetTopic(msg.LastParam(), i.client.Host, channel)
}

func (i *ircHandler) noTopic(msg *irc.Message) {
	channel := msg.Params[1]

	i.state.sendJSON("topic", Topic{
		Server:  i.client.Host,
		Channel: channel,
	})

	channelStore.SetTopic("", i.client.Host, channel)
}

func (i *ircHandler) names(msg *irc.Message) {
	users := strings.Split(strings.TrimSuffix(msg.LastParam(), " "), " ")
	userBuffer := i.userBuffers[msg.Params[2]]
	i.userBuffers[msg.Params[2]] = append(userBuffer, users...)
}

func (i *ircHandler) namesEnd(msg *irc.Message) {
	channel := msg.Params[1]
	users := i.userBuffers[channel]

	i.state.sendJSON("users", Userlist{
		Server:  i.client.Host,
		Channel: channel,
		Users:   users,
	})

	channelStore.SetUsers(users, i.client.Host, channel)
	delete(i.userBuffers, channel)
}

func (i *ircHandler) motdStart(msg *irc.Message) {
	i.motdBuffer.Server = i.client.Host
	i.motdBuffer.Title = msg.LastParam()
}

func (i *ircHandler) motd(msg *irc.Message) {
	i.motdBuffer.Content = append(i.motdBuffer.Content, msg.LastParam())
}

func (i *ircHandler) motdEnd(msg *irc.Message) {
	i.state.sendJSON("motd", i.motdBuffer)
	i.motdBuffer = MOTD{}
}

func (i *ircHandler) list(msg *irc.Message) {
	if i.listBuffer == nil && i.state.Bool("update_chanlist_"+i.client.Host) {
		i.listBuffer = storage.NewMapChannelListIndex()
	}

	if i.listBuffer != nil {
		c, _ := strconv.Atoi(msg.Params[2])
		i.listBuffer.Add(&storage.ChannelListItem{
			Name:      msg.Params[1],
			UserCount: c,
			Topic:     msg.LastParam(),
		})
	}
}

func (i *ircHandler) listEnd(msg *irc.Message) {
	if i.listBuffer != nil {
		i.state.Set("update_chanlist_"+i.client.Host, false)

		go func(idx storage.ChannelListIndex) {
			idx.Finish()
			channelIndexes.Set(i.client.Host, idx)
		}(i.listBuffer)

		i.listBuffer = nil
	}
}

func (i *ircHandler) badNick(msg *irc.Message) {
	i.state.sendJSON("nick_fail", NickFail{
		Server: i.client.Host,
	})
}

func (i *ircHandler) forward(msg *irc.Message) {
	if len(msg.Params) > 2 {
		i.state.sendJSON("channel_forward", ChannelForward{
			Server: i.client.Host,
			Old:    msg.Params[1],
			New:    msg.Params[2],
		})
	}
}

func (i *ircHandler) error(msg *irc.Message) {
	i.state.sendJSON("error", IRCError{
		Server:  i.client.Host,
		Message: msg.LastParam(),
	})
}

func (i *ircHandler) receiveDCCSend(pack *irc.DCCSend, msg *irc.Message) {
	cfg := i.state.srv.Config()

	if cfg.DCC.Enabled {
		if cfg.DCC.Autoget.Enabled {
			file, err := os.OpenFile(storage.Path.DownloadedFile(i.state.user.Username, pack.File), os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return
			}
			defer file.Close()

			irc.DownloadDCC(file, pack, i.dccProgress)
		} else {
			i.state.setPendingDCC(pack.File, pack)

			i.state.sendJSON("dcc_send", DCCSend{
				Server:   i.client.Host,
				From:     msg.Nick,
				Filename: pack.File,
				URL: fmt.Sprintf("%s://%s/downloads/%s/%s",
					i.state.String("scheme"), i.state.String("host"), i.state.user.Username, pack.File),
			})

			time.Sleep(150 * time.Second)
			i.state.deletePendingDCC(pack.File)
		}
	}
}

func (i *ircHandler) initHandlers() {
	i.handlers = map[string]func(*irc.Message){
		irc.NICK:                 i.nick,
		irc.JOIN:                 i.join,
		irc.PART:                 i.part,
		irc.MODE:                 i.mode,
		irc.PRIVMSG:              i.message,
		irc.NOTICE:               i.message,
		irc.QUIT:                 i.quit,
		irc.TOPIC:                i.topic,
		irc.ERROR:                i.error,
		irc.RPL_WELCOME:          i.info,
		irc.RPL_YOURHOST:         i.info,
		irc.RPL_CREATED:          i.info,
		irc.RPL_ISUPPORT:         i.features,
		irc.RPL_LUSERCLIENT:      i.info,
		irc.RPL_LUSEROP:          i.info,
		irc.RPL_LUSERUNKNOWN:     i.info,
		irc.RPL_LUSERCHANNELS:    i.info,
		irc.RPL_LUSERME:          i.info,
		irc.RPL_WHOISUSER:        i.whoisUser,
		irc.RPL_WHOISSERVER:      i.whoisServer,
		irc.RPL_WHOISCHANNELS:    i.whoisChannels,
		irc.RPL_ENDOFWHOIS:       i.whoisEnd,
		irc.RPL_NOTOPIC:          i.noTopic,
		irc.RPL_TOPIC:            i.topic,
		irc.RPL_NAMREPLY:         i.names,
		irc.RPL_ENDOFNAMES:       i.namesEnd,
		irc.RPL_MOTDSTART:        i.motdStart,
		irc.RPL_MOTD:             i.motd,
		irc.RPL_ENDOFMOTD:        i.motdEnd,
		irc.RPL_LIST:             i.list,
		irc.RPL_LISTEND:          i.listEnd,
		irc.ERR_ERRONEUSNICKNAME: i.badNick,
		irc.ERR_FORWARD:          i.forward,
	}
}

func (i *ircHandler) log(v ...interface{}) {
	log.Println("[IRC]", i.state.user.ID, i.client.Host, fmt.Sprint(v...))
}

func (i *ircHandler) sendDCCInfo(message string, log bool, a ...interface{}) {
	msg := Message{
		Server:  i.client.Host,
		From:    "@dcc",
		Content: fmt.Sprintf(message, a...),
	}
	i.state.sendJSON("pm", msg)

	if log {
		i.state.user.AddOpenDM(msg.Server, msg.From)
		i.state.user.LogMessage(betterguid.New(), msg.Server, msg.From, msg.From, msg.Content)
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

func isExcludedError(cmd string) bool {
	for _, err := range excludedErrors {
		if cmd == err {
			return true
		}
	}
	return false
}

func formatIRCError(msg *irc.Message) string {
	errMsg := strings.TrimSuffix(msg.LastParam(), ".")
	if len(msg.Params) > 2 {
		for _, c := range msg.LastParam() {
			if unicode.IsLower(c) {
				return msg.Params[1] + " " + errMsg
			}
			return msg.Params[1] + ": " + errMsg
		}
	}
	return errMsg
}

func printMessage(msg *irc.Message, i *irc.Client) {
	log.Println(i.GetNick()+":", msg.Prefix, msg.Command, msg.Params)
}
