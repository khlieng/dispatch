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
	motdBuffer  MOTD
	listBuffer  storage.ChannelListIndex
	dccProgress chan irc.DownloadProgress

	handlers map[string]func(*irc.Message)
}

func newIRCHandler(client *irc.Client, state *State) *ircHandler {
	i := &ircHandler{
		client:      client,
		state:       state,
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
				i.state.deleteNetwork(i.client.Host())
				return
			}

			i.dispatchMessage(msg)

		case state := <-i.client.ConnectionChanged:
			i.state.sendJSON("connection_update", newConnectionUpdate(i.client.Host(), state))

			if network, ok := i.state.network(i.client.Host()); ok {
				var err string
				if state.Error != nil {
					err = state.Error.Error()
				}
				network.SetStatus(state.Connected, err)
			}

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
			Network: i.client.Host(),
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
	nick := Nick{
		Network: i.client.Host(),
		Old:     msg.Sender,
		New:     msg.LastParam(),
	}

	i.state.sendJSON("nick", nick)

	if i.client.Is(nick.New) {
		if network, ok := i.state.network(nick.Network); ok {
			network.SetNick(nick.New)
			go network.Save()
		}
	}

	channels := irc.GetNickChannels(msg)
	go i.state.user.LogEvent(nick.Network, "nick", []string{nick.Old, nick.New}, channels...)
}

func (i *ircHandler) join(msg *irc.Message) {
	host := i.client.Host()

	i.state.sendJSON("join", Join{
		Network:  host,
		User:     msg.Sender,
		Channels: msg.Params,
	})

	channel := msg.Params[0]

	if i.client.Is(msg.Sender) {
		// In case no topic is set and there's a cached one that needs to be cleared
		i.client.Topic(channel)

		if network, ok := i.state.network(host); ok {
			if ch := network.Channel(channel); ch != nil {
				ch.SetJoined(true)
			} else {
				i.state.sendLastMessages(host, channel, 50)

				ch = network.NewChannel(channel)
				ch.SetJoined(true)
				network.AddChannel(ch)
				go ch.Save()
			}
		}
	}

	go i.state.user.LogEvent(host, "join", []string{msg.Sender}, channel)
}

func (i *ircHandler) part(msg *irc.Message) {
	part := Part{
		Network: i.client.Host(),
		User:    msg.Sender,
		Channel: msg.Params[0],
	}

	params := []string{part.User}

	if len(msg.Params) == 2 {
		part.Reason = msg.Params[1]
		params = append(params, part.Reason)
	}

	i.state.sendJSON("part", part)

	if i.client.Is(msg.Sender) {
		go i.state.user.RemoveChannel(part.Network, part.Channel)
	}

	go i.state.user.LogEvent(part.Network, "part", params, part.Channel)
}

func (i *ircHandler) kick(msg *irc.Message) {
	if len(msg.Params) < 2 {
		return
	}

	kick := Kick{
		Network: i.client.Host(),
		Channel: msg.Params[0],
		Sender:  msg.Sender,
		User:    msg.Params[1],
	}

	params := []string{kick.User, kick.Sender}

	if len(msg.Params) > 2 {
		kick.Reason = msg.Params[2]
		params = append(params, kick.Reason)
	}

	i.state.sendJSON("kick", kick)

	go i.state.user.LogEvent(kick.Network, "kick", params, kick.Channel)

	if i.client.Is(kick.User) {
		if network, ok := i.state.network(kick.Network); ok {
			network.Channel(kick.Channel).SetJoined(false)
		}
	}
}

func (i *ircHandler) mode(msg *irc.Message) {
	if mode := irc.GetMode(msg); mode != nil {
		i.state.sendJSON("mode", Mode{
			Mode: mode,
		})
	}
}

func (i *ircHandler) message(msg *irc.Message) {
	if ctcp := msg.ToCTCP(); ctcp != nil {
		if ctcp.Command == "DCC" && strings.HasPrefix(ctcp.Params, "SEND") {
			if pack := i.client.ParseDCCSend(ctcp); pack != nil {
				go i.receiveDCCSend(pack, msg)
				return
			}
		} else if ctcp.Command != "ACTION" {
			return
		}
	}

	message := Message{
		ID:      betterguid.New(),
		Network: i.client.Host(),
		From:    msg.Sender,
		Content: msg.LastParam(),
	}
	target := msg.Params[0]

	if i.client.Is(target) {
		i.state.sendJSON("pm", message)

		if !msg.IsFromServer() {
			i.state.user.AddOpenDM(i.client.Host(), message.From)
		}

		target = message.From
	} else {
		message.To = target
		i.state.sendJSON("message", message)
	}

	if target != "*" && !msg.IsFromServer() {
		go i.state.user.LogMessage(&storage.Message{
			ID:      message.ID,
			Network: message.Network,
			From:    message.From,
			To:      target,
			Content: message.Content,
		})
	}
}

func (i *ircHandler) quit(msg *irc.Message) {
	i.state.sendJSON("quit", Quit{
		Network: i.client.Host(),
		User:    msg.Sender,
		Reason:  msg.LastParam(),
	})

	channels := irc.GetQuitChannels(msg)

	go i.state.user.LogEvent(i.client.Host(), "quit", []string{msg.Sender, msg.LastParam()}, channels...)
}

func (i *ircHandler) info(msg *irc.Message) {
	if msg.Command == irc.RPL_WELCOME {
		i.state.sendJSON("nick", Nick{
			Network: i.client.Host(),
			New:     msg.Params[0],
		})

		_, needsUpdate := channelIndexes.Get(i.client.Host())
		if needsUpdate {
			i.listBuffer = storage.NewMapChannelListIndex()
			i.client.List()
		}

		go i.state.user.SetNick(msg.Params[0], i.client.Host())
	}

	i.state.sendJSON("pm", Message{
		Network: i.client.Host(),
		From:    msg.Sender,
		Content: strings.Join(msg.Params[1:], " "),
	})
}

func (i *ircHandler) features(msg *irc.Message) {
	features := i.client.Features.Map()

	i.state.sendJSON("features", Features{
		Network:  i.client.Host(),
		Features: features,
	})

	if network, ok := i.state.network(i.client.Host()); ok {
		network.SetFeatures(features)

		if name := i.client.Features.String("NETWORK"); name != "" {
			network.SetName(name)
			go network.Save()
		}
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
		nick = msg.Sender

		go i.state.user.LogEvent(i.client.Host(), "topic", []string{nick, msg.LastParam()}, channel)
	} else {
		channel = msg.Params[1]
	}

	i.state.sendJSON("topic", Topic{
		Network: i.client.Host(),
		Channel: channel,
		Topic:   msg.LastParam(),
		Nick:    nick,
	})

	if network, ok := i.state.network(i.client.Host()); ok {
		network.Channel(channel).SetTopic(msg.LastParam())
	}
}

func (i *ircHandler) noTopic(msg *irc.Message) {
	channel := msg.Params[1]

	i.state.sendJSON("topic", Topic{
		Network: i.client.Host(),
		Channel: channel,
	})

	if network, ok := i.state.network(i.client.Host()); ok {
		network.Channel(channel).SetTopic("")
	}

}

func (i *ircHandler) namesEnd(msg *irc.Message) {
	i.state.sendJSON("users", Userlist{
		Network: i.client.Host(),
		Channel: msg.Params[1],
		Users:   irc.GetNamreplyUsers(msg),
	})
}

func (i *ircHandler) motdStart(msg *irc.Message) {
	i.motdBuffer.Network = i.client.Host()
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
	if i.listBuffer == nil && i.state.Bool("update_chanlist_"+i.client.Host()) {
		i.listBuffer = storage.NewMapChannelListIndex()
	}

	if i.listBuffer != nil {
		userCount, _ := strconv.Atoi(msg.Params[2])
		i.listBuffer.Add(&storage.ChannelListItem{
			Name:      msg.Params[1],
			UserCount: userCount,
			Topic:     msg.LastParam(),
		})
	}
}

func (i *ircHandler) listEnd(msg *irc.Message) {
	if i.listBuffer != nil {
		i.state.Set("update_chanlist_"+i.client.Host(), false)

		go func(idx storage.ChannelListIndex) {
			idx.Finish()
			channelIndexes.Set(i.client.Host(), idx)
		}(i.listBuffer)

		i.listBuffer = nil
	}
}

func (i *ircHandler) badNick(msg *irc.Message) {
	i.state.sendJSON("nick_fail", NickFail{
		Network: i.client.Host(),
	})
}

func (i *ircHandler) forward(msg *irc.Message) {
	if len(msg.Params) > 2 {
		i.state.sendJSON("channel_forward", ChannelForward{
			Network: i.client.Host(),
			Old:     msg.Params[1],
			New:     msg.Params[2],
		})
	}
}

func (i *ircHandler) error(msg *irc.Message) {
	i.state.sendJSON("error", IRCError{
		Network: i.client.Host(),
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

			pack.Download(file, i.dccProgress)
		} else {
			i.state.setPendingDCC(pack.File, pack)

			i.state.sendJSON("dcc_send", DCCSend{
				Network:  i.client.Host(),
				From:     msg.Sender,
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
		irc.KICK:                 i.kick,
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
	log.Println("[IRC]", i.state.user.ID, i.client.Host(), fmt.Sprint(v...))
}

func (i *ircHandler) sendDCCInfo(message string, log bool, a ...interface{}) {
	msg := Message{
		Network: i.client.Host(),
		From:    "@dcc",
		Content: fmt.Sprintf(message, a...),
	}
	i.state.sendJSON("pm", msg)

	if log {
		i.state.user.AddOpenDM(msg.Network, msg.From)
		i.state.user.LogMessage(&storage.Message{
			Network: msg.Network,
			From:    msg.From,
			Content: msg.Content,
		})
	}
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
	log.Println(i.GetNick()+":", msg.Sender, msg.Command, msg.Params)
}
