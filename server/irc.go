package server

import (
	"crypto/tls"

	"github.com/khlieng/dispatch/irc"
	"github.com/khlieng/dispatch/storage"
)

func reconnectIRC() {
	for _, user := range storage.LoadUsers() {
		session := NewSession()
		session.user = user
		sessions[user.UUID] = session
		go session.write()

		channels := user.GetChannels()

		for _, server := range user.GetServers() {
			i := irc.NewClient(server.Nick, server.Username)
			i.TLS = server.TLS
			i.Password = server.Password
			i.Realname = server.Realname

			if cert := user.GetCertificate(); cert != nil {
				i.TLSConfig = &tls.Config{
					Certificates: []tls.Certificate{*cert},
				}
			}

			i.Connect(server.Address)
			session.setIRC(i.Host, i)
			go newIRCHandler(i, session).run()

			var joining []string
			for _, channel := range channels {
				if channel.Server == server.Address {
					joining = append(joining, channel.Name)
				}
			}
			i.Join(joining...)
		}
	}
}
