package server

import (
	"log"
	"strings"

	"github.com/khlieng/name_pending/irc"
	"github.com/khlieng/name_pending/storage"
)

func handleIRC(client *irc.Client, session *Session) {
	var whois WhoisReply
	userBuffers := make(map[string][]string)
	var motd MOTD

	for {
		msg, ok := <-client.Messages
		if !ok {
			session.deleteIRC(client.Host)
			return
		}

		switch msg.Command {
		case irc.Nick:
			session.sendJSON("nick", Nick{
				Server: client.Host,
				Old:    msg.Nick,
				New:    msg.Trailing,
			})

			channelStore.RenameUser(msg.Nick, msg.Trailing, client.Host)

		case irc.Join:
			session.sendJSON("join", Join{
				Server:   client.Host,
				User:     msg.Nick,
				Channels: msg.Params,
			})

			channelStore.AddUser(msg.Nick, client.Host, msg.Params[0])

			if msg.Nick == client.GetNick() {
				session.user.AddChannel(storage.Channel{
					Server: client.Host,
					Name:   msg.Params[0],
				})
			}

		case irc.Part:
			session.sendJSON("part", Part{
				Join: Join{
					Server:   client.Host,
					User:     msg.Nick,
					Channels: msg.Params,
				},
				Reason: msg.Trailing,
			})

			channelStore.RemoveUser(msg.Nick, client.Host, msg.Params[0])

			if msg.Nick == client.GetNick() {
				session.user.RemoveChannel(client.Host, msg.Params[0])
			}

		case irc.Mode:
			target := msg.Params[0]
			if len(msg.Params) > 2 && isChannel(target) {
				mode := parseMode(msg.Params[1])
				mode.Server = client.Host
				mode.Channel = target
				mode.User = msg.Params[2]

				session.sendJSON("mode", mode)

				channelStore.SetMode(client.Host, target, msg.Params[2], mode.Add, mode.Remove)
			}

		case irc.Privmsg, irc.Notice:
			if msg.Params[0] == client.GetNick() {
				session.sendJSON("pm", Chat{
					Server:  client.Host,
					From:    msg.Nick,
					Message: msg.Trailing,
				})
			} else {
				session.sendJSON("message", Chat{
					Server:  client.Host,
					From:    msg.Nick,
					To:      msg.Params[0],
					Message: msg.Trailing,
				})
			}

			if msg.Params[0] != "*" {
				go session.user.LogMessage(client.Host, msg.Nick, msg.Params[0], msg.Trailing)
			}

		case irc.Quit:
			session.sendJSON("quit", Quit{
				Server: client.Host,
				User:   msg.Nick,
				Reason: msg.Trailing,
			})

			channelStore.RemoveUserAll(msg.Nick, client.Host)

		case irc.ReplyWelcome,
			irc.ReplyYourHost,
			irc.ReplyCreated,
			irc.ReplyLUserClient,
			irc.ReplyLUserOp,
			irc.ReplyLUserUnknown,
			irc.ReplyLUserChannels,
			irc.ReplyLUserMe:
			session.sendJSON("pm", Chat{
				Server:  client.Host,
				From:    msg.Nick,
				Message: strings.Join(msg.Params[1:], " "),
			})

		case irc.ReplyWhoisUser:
			whois.Nick = msg.Params[1]
			whois.Username = msg.Params[2]
			whois.Host = msg.Params[3]
			whois.Realname = msg.Params[5]

		case irc.ReplyWhoisServer:
			whois.Server = msg.Params[2]

		case irc.ReplyWhoisChannels:
			whois.Channels = append(whois.Channels, strings.Split(strings.TrimRight(msg.Trailing, " "), " ")...)

		case irc.ReplyEndOfWhois:
			session.sendJSON("whois", whois)

			whois = WhoisReply{}

		case irc.ReplyTopic:
			session.sendJSON("topic", Topic{
				Server:  client.Host,
				Channel: msg.Params[1],
				Topic:   msg.Trailing,
			})

			channelStore.SetTopic(msg.Trailing, client.Host, msg.Params[1])

		case irc.ReplyNamReply:
			users := strings.Split(msg.Trailing, " ")
			userBuffer := userBuffers[msg.Params[2]]
			userBuffers[msg.Params[2]] = append(userBuffer, users...)

		case irc.ReplyEndOfNames:
			channel := msg.Params[1]
			users := userBuffers[channel]

			session.sendJSON("users", Userlist{
				Server:  client.Host,
				Channel: channel,
				Users:   users,
			})

			channelStore.SetUsers(users, client.Host, channel)
			delete(userBuffers, channel)

		case irc.ReplyMotdStart:
			motd.Server = client.Host
			motd.Title = msg.Trailing

		case irc.ReplyMotd:
			motd.Content = append(motd.Content, msg.Trailing)

		case irc.ReplyEndOfMotd:
			session.sendJSON("motd", motd)

			motd = MOTD{}

		default:
			printMessage(msg, client)
		}
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
