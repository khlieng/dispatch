package main

import (
	"encoding/json"
	"log"

	"golang.org/x/net/websocket"

	"github.com/khlieng/irc/storage"
)

func handleWS(ws *websocket.Conn) {
	defer ws.Close()

	var session *Session
	var UUID string
	var req WSRequest

	addr := ws.Request().RemoteAddr

	log.Println(addr, "connected")

	for {
		err := websocket.JSON.Receive(ws, &req)
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

				log.Println(addr, "attached to existing IRC connections")

				channels := session.user.GetChannels()
				for i, channel := range channels {
					channels[i].Users = channelStore.GetUsers(channel.Server, channel.Name)
				}

				session.sendJSON("channels", channels)
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
				irc.TLS = true
				irc.Connect(data.Server)

				session.setIRC(irc.Host, irc)

				go handleMessages(irc, session)

				session.user.AddServer(storage.Server{
					Address:  irc.Host,
					Nick:     data.Nick,
					Username: data.Username,
				})
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
			var data Join

			json.Unmarshal(req.Request, &data)

			if irc, ok := session.getIRC(data.Server); ok {
				irc.Part(data.Channels...)
			}

		case "chat":
			var data Chat

			json.Unmarshal(req.Request, &data)

			if irc, ok := session.getIRC(data.Server); ok {
				irc.Privmsg(data.To, data.Message)
			}
		}
	}
}
