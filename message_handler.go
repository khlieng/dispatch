package main

import (
	"log"
	"strings"

	"github.com/khlieng/name_pending/storage"
)

func handleMessages(irc *IRC, session *Session) {
	var whois WhoisReply
	userBuffers := make(map[string][]string)
	var motd MOTD

	for msg := range irc.Messages {
		switch msg.Command {
		case NICK:
			session.sendJSON("nick", Nick{
				Server: irc.Host,
				Old:    msg.Prefix,
				New:    msg.Trailing,
			})

			channelStore.RenameUser(msg.Prefix, msg.Trailing, irc.Host)

		case JOIN:
			user := msg.Prefix

			session.sendJSON("join", Join{
				Server:   irc.Host,
				User:     user,
				Channels: msg.Params,
			})

			channelStore.AddUser(user, irc.Host, msg.Params[0])

			if user == irc.GetNick() {
				session.user.AddChannel(storage.Channel{
					Server: irc.Host,
					Name:   msg.Params[0],
				})
			}

		case PART:
			user := msg.Prefix

			session.sendJSON("part", Part{
				Join: Join{
					Server:   irc.Host,
					User:     user,
					Channels: msg.Params,
				},
				Reason: msg.Trailing,
			})

			channelStore.RemoveUser(user, irc.Host, msg.Params[0])

			if user == irc.GetNick() {
				session.user.RemoveChannel(irc.Host, msg.Params[0])
			}

		case MODE:
			target := msg.Params[0]
			if len(msg.Params) > 2 && isChannel(target) {
				mode := parseMode(msg.Params[1])
				mode.Server = irc.Host
				mode.Channel = target
				mode.User = msg.Params[2]

				session.sendJSON("mode", mode)

				channelStore.SetMode(irc.Host, target, msg.Params[2], mode.Add, mode.Remove)
			}

		case PRIVMSG, NOTICE:
			if msg.Params[0] == irc.GetNick() {
				session.sendJSON("pm", Chat{
					Server:  irc.Host,
					From:    msg.Prefix,
					Message: msg.Trailing,
				})
			} else {
				session.sendJSON("message", Chat{
					Server:  irc.Host,
					From:    msg.Prefix,
					To:      msg.Params[0],
					Message: msg.Trailing,
				})
			}

		case QUIT:
			user := msg.Prefix

			session.sendJSON("quit", Quit{
				Server: irc.Host,
				User:   user,
				Reason: msg.Trailing,
			})

			channelStore.RemoveUserAll(user, irc.Host)

		case RPL_WELCOME,
			RPL_YOURHOST,
			RPL_CREATED,
			RPL_LUSERCLIENT,
			RPL_LUSEROP,
			RPL_LUSERUNKNOWN,
			RPL_LUSERCHANNELS,
			RPL_LUSERME:
			session.sendJSON("pm", Chat{
				Server:  irc.Host,
				From:    msg.Prefix,
				Message: strings.Join(msg.Params[1:], " "),
			})

		case RPL_WHOISUSER:
			whois.Nick = msg.Params[1]
			whois.Username = msg.Params[2]
			whois.Host = msg.Params[3]
			whois.Realname = msg.Params[5]

		case RPL_WHOISSERVER:
			whois.Server = msg.Params[2]

		case RPL_WHOISCHANNELS:
			whois.Channels = append(whois.Channels, strings.Split(strings.TrimRight(msg.Trailing, " "), " ")...)

		case RPL_ENDOFWHOIS:
			session.sendJSON("whois", whois)

			whois = WhoisReply{}

		case RPL_TOPIC:
			session.sendJSON("topic", Topic{
				Server:  irc.Host,
				Channel: msg.Params[1],
				Topic:   msg.Trailing,
			})

			channelStore.SetTopic(msg.Trailing, irc.Host, msg.Params[1])

		case RPL_NAMREPLY:
			users := strings.Split(msg.Trailing, " ")

			/*for i, user := range users {
				users[i] = strings.TrimLeft(user, "@+")
			}*/

			userBuffer := userBuffers[msg.Params[2]]
			userBuffers[msg.Params[2]] = append(userBuffer, users...)

		case RPL_ENDOFNAMES:
			channel := msg.Params[1]
			users := userBuffers[channel]

			session.sendJSON("users", Userlist{
				Server:  irc.Host,
				Channel: channel,
				Users:   users,
			})

			channelStore.SetUsers(users, irc.Host, channel)
			delete(userBuffers, channel)

		case RPL_MOTDSTART:
			motd.Server = irc.Host
			motd.Title = msg.Trailing

		case RPL_MOTD:
			motd.Content = append(motd.Content, msg.Trailing)

		case RPL_ENDOFMOTD:
			session.sendJSON("motd", motd)

			motd = MOTD{}

		default:
			printMessage(msg, irc)
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

func printMessage(msg *Message, irc *IRC) {
	log.Println(irc.GetNick()+":", msg.Prefix, msg.Command, msg.Params, msg.Trailing)
}
