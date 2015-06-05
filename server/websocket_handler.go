package server

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/khlieng/name_pending/Godeps/_workspace/src/github.com/gorilla/websocket"

	"github.com/khlieng/name_pending/irc"
	"github.com/khlieng/name_pending/storage"
)

func handleWS(conn *websocket.Conn) {
	defer conn.Close()

	var session *Session
	var UUID string

	addr := conn.RemoteAddr().String()

	ws := newConn(conn)
	defer ws.close()
	go ws.send()
	go ws.recv()

	log.Println(addr, "connected")

	for {
		req, ok := <-ws.in
		if !ok {
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

				session.sendJSON("servers", nil)

				go session.write()
			}

			session.setWS(addr, ws)

		case "connect":
			var data Connect

			json.Unmarshal(req.Request, &data)

			if _, ok := session.getIRC(data.Server); !ok {
				log.Println(addr, "connecting to", data.Server)

				i := irc.NewClient(data.Nick, data.Username)
				i.TLS = data.TLS
				i.Password = data.Password
				i.Realname = data.Realname

				if idx := strings.Index(data.Server, ":"); idx < 0 {
					session.setIRC(data.Server, i)
				} else {
					session.setIRC(data.Server[:idx], i)
				}

				go func() {
					i.Connect(data.Server)
					go handleIRC(i, session)

					session.user.AddServer(storage.Server{
						Name:     data.Name,
						Address:  i.Host,
						TLS:      data.TLS,
						Password: data.Password,
						Nick:     data.Nick,
						Username: data.Username,
						Realname: data.Realname,
					})
				}()
			} else {
				log.Println(addr, "already connected to", data.Server)
			}

		case "join":
			var data Join

			json.Unmarshal(req.Request, &data)

			if i, ok := session.getIRC(data.Server); ok {
				i.Join(data.Channels...)
			}

		case "part":
			var data Part

			json.Unmarshal(req.Request, &data)

			if i, ok := session.getIRC(data.Server); ok {
				i.Part(data.Channels...)
			}

		case "quit":
			var data Quit

			json.Unmarshal(req.Request, &data)

			if i, ok := session.getIRC(data.Server); ok {
				i.Quit()
				session.deleteIRC(data.Server)
				channelStore.RemoveUserAll(i.GetNick(), data.Server)
				session.user.RemoveServer(data.Server)
			}

		case "chat":
			var data Chat

			json.Unmarshal(req.Request, &data)

			if i, ok := session.getIRC(data.Server); ok {
				i.Privmsg(data.To, data.Message)
			}

		case "nick":
			var data Nick

			json.Unmarshal(req.Request, &data)

			if i, ok := session.getIRC(data.Server); ok {
				i.Nick(data.New)
				session.user.SetNick(data.New, data.Server)
			}

		case "invite":
			var data Invite

			json.Unmarshal(req.Request, &data)

			if i, ok := session.getIRC(data.Server); ok {
				i.Invite(data.User, data.Channel)
			}

		case "kick":
			var data Invite

			json.Unmarshal(req.Request, &data)

			if i, ok := session.getIRC(data.Server); ok {
				i.Kick(data.Channel, data.User)
			}

		case "whois":
			var data Whois

			json.Unmarshal(req.Request, &data)

			if i, ok := session.getIRC(data.Server); ok {
				i.Whois(data.User)
			}

		case "away":
			var data Away

			json.Unmarshal(req.Request, &data)

			if i, ok := session.getIRC(data.Server); ok {
				i.Away(data.Message)
			}

		case "search":
			go func() {
				var data SearchRequest

				json.Unmarshal(req.Request, &data)

				results, err := session.user.SearchMessages(data.Server, data.Channel, data.Phrase)
				if err != nil {
					log.Println(err)
					return
				}

				session.sendJSON("search", SearchResult{
					Server:  data.Server,
					Channel: data.Channel,
					Results: results,
				})
			}()
		}
	}
}
