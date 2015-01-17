package main

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/khlieng/irc/storage"
)

func handleMessages(irc *IRC, session *Session) {
	userBuffers := make(map[string][]string)
	var motd MOTD
	var motdContent bytes.Buffer

	for msg := range irc.Messages {
		switch msg.Command {
		case JOIN:
			user := parseUser(msg.Prefix)

			session.sendJSON("join", Join{
				Server:   irc.Host,
				User:     user,
				Channels: msg.Params,
			})

			channelStore.AddUser(user, irc.Host, msg.Params[0])

			if user == irc.nick {
				session.user.AddChannel(storage.Channel{
					Server: irc.Host,
					Name:   msg.Params[0],
				})
			}

		case PART:
			user := parseUser(msg.Prefix)

			session.sendJSON("part", Join{
				Server:   irc.Host,
				User:     user,
				Channels: msg.Params,
			})

			channelStore.RemoveUser(user, irc.Host, msg.Params[0])

			if user == irc.nick {
				session.user.RemoveChannel(irc.Host, msg.Params[0])
			}

		case PRIVMSG, NOTICE:
			if msg.Params[0] == irc.nick {
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
			/*
				session.sendJSON("quit", Quit{
					Server: irc.Host,
					User: user,
				})
			*/

		case RPL_WELCOME,
			RPL_YOURHOST,
			RPL_LUSERCLIENT,
			RPL_LUSEROP,
			RPL_LUSERUNKNOWN,
			RPL_LUSERCHANNELS,
			RPL_LUSERME:
			session.sendJSON("pm", Chat{
				Server:  irc.Host,
				From:    msg.Prefix,
				Message: strings.Join(append(msg.Params[1:], msg.Trailing), " "),
			})

		case RPL_TOPIC:
			session.sendJSON("topic", Topic{
				Server:  irc.Host,
				Channel: msg.Params[1],
				Topic:   msg.Trailing,
			})

		case RPL_NAMREPLY:
			users := strings.Split(msg.Trailing, " ")

			for i, user := range users {
				users[i] = strings.TrimLeft(user, "@+")
			}

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
			motdContent.WriteString(msg.Trailing)
			motdContent.WriteRune('\n')

		case RPL_ENDOFMOTD:
			motd.Content = motdContent.String()

			session.sendJSON("motd", motd)

			motdContent.Reset()
			motd = MOTD{}

		default:
			printMessage(msg, irc)
		}
	}
}

func printMessage(msg *Message, irc *IRC) {
	fmt.Printf("%s: %s %s %s\n", irc.nick, msg.Prefix, msg.Command, msg.Params)

	if msg.Trailing != "" {
		fmt.Println(msg.Trailing)
	}
	fmt.Println()
}
