package server

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/khlieng/name_pending/Godeps/_workspace/src/github.com/gorilla/websocket"

	"github.com/khlieng/name_pending/storage"
)

func handleWS(ws *websocket.Conn) {
	defer ws.Close()

	var session *Session
	var UUID string
	var req WSRequest

	addr := ws.RemoteAddr().String()

	log.Println(addr, "connected")

	for {
		err := ws.ReadJSON(&req)
		if err != nil {
			if session != nil {
				session.deleteWS(addr)
			}

			log.Println(addr, "disconnected")
			return
		}

		switch req.Type {
		case "uuid":
			json.Unmarshal(req.Request, &UUID)

			log.Println(addr, "set UUID", UUID)

			sessionLock.Lock()

			if storedSession, exists := sessions[UUID]; exists {
				sessionLock.Unlock()
				session = storedSession

				log.Println(addr, "attached to", session.numIRC(), "existing IRC connections")

				channels := session.user.GetChannels()
				for i, channel := range channels {
					channels[i].Topic = channelStore.GetTopic(channel.Server, channel.Name)
				}

				session.sendJSON("channels", channels)
				session.sendJSON("servers", session.user.GetServers())

				for _, channel := range channels {
					session.sendJSON("users", Userlist{
						Server:  channel.Server,
						Channel: channel.Name,
						Users:   channelStore.GetUsers(channel.Server, channel.Name),
					})
				}
			} else {
				session = NewSession()
				session.user = storage.NewUser(UUID)

				sessions[UUID] = session
				sessionLock.Unlock()

				go session.write()
			}

			session.setWS(addr, ws)

		case "connect":
			var data Connect

			json.Unmarshal(req.Request, &data)

			if _, ok := session.getIRC(data.Server); !ok {
				log.Println(addr, "connecting to", data.Server)

				irc := NewIRC(data.Nick, data.Username)
				irc.TLS = data.TLS
				irc.Password = data.Password
				irc.Realname = data.Realname

				if idx := strings.Index(data.Server, ":"); idx < 0 {
					session.setIRC(data.Server, irc)
				} else {
					session.setIRC(data.Server[:idx], irc)
				}

				go func() {
					err := irc.Connect(data.Server)
					if err != nil {
						session.deleteIRC(irc.Host)
						session.sendError(err, irc.Host)
						log.Println(err)
					} else {
						go handleMessages(irc, session)

						session.user.AddServer(storage.Server{
							Name:     data.Name,
							Address:  irc.Host,
							TLS:      data.TLS,
							Password: data.Password,
							Nick:     data.Nick,
							Username: data.Username,
							Realname: data.Realname,
						})
					}
				}()
			} else {
				log.Println(addr, "already connected to", data.Server)
			}

		case "join":
			var data Join

			json.Unmarshal(req.Request, &data)

			if irc, ok := session.getIRC(data.Server); ok {
				irc.Join(data.Channels...)
			}

		case "part":
			var data Part

			json.Unmarshal(req.Request, &data)

			if irc, ok := session.getIRC(data.Server); ok {
				irc.Part(data.Channels...)
			}

		case "quit":
			var data Quit

			json.Unmarshal(req.Request, &data)

			if irc, ok := session.getIRC(data.Server); ok {
				irc.Quit()
				session.deleteIRC(data.Server)
				channelStore.RemoveUserAll(irc.GetNick(), data.Server)
				session.user.RemoveServer(data.Server)
			}

		case "chat":
			var data Chat

			json.Unmarshal(req.Request, &data)

			if irc, ok := session.getIRC(data.Server); ok {
				irc.Privmsg(data.To, data.Message)
			}

		case "nick":
			var data Nick

			json.Unmarshal(req.Request, &data)

			if irc, ok := session.getIRC(data.Server); ok {
				irc.Nick(data.New)
				session.user.SetNick(data.New, data.Server)
			}

		case "invite":
			var data Invite

			json.Unmarshal(req.Request, &data)

			if irc, ok := session.getIRC(data.Server); ok {
				irc.Invite(data.User, data.Channel)
			}

		case "kick":
			var data Invite

			json.Unmarshal(req.Request, &data)

			if irc, ok := session.getIRC(data.Server); ok {
				irc.Kick(data.Channel, data.User)
			}

		case "whois":
			var data Whois

			json.Unmarshal(req.Request, &data)

			if irc, ok := session.getIRC(data.Server); ok {
				irc.Whois(data.User)
			}

		case "away":
			var data Away

			json.Unmarshal(req.Request, &data)

			if irc, ok := session.getIRC(data.Server); ok {
				irc.Away(data.Message)
			}
		}
	}
}
