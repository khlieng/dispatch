package server

import (
	"crypto/tls"
	"net"

	"github.com/khlieng/dispatch/irc"
	"github.com/khlieng/dispatch/storage"
)

func reconnectIRC() {
	for _, user := range storage.LoadUsers() {
		session := NewSession(user)
		sessions[user.ID] = session
		go session.run()

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

			session.setIRC(server.Host, i)

			if server.Port != "" {
				i.Connect(net.JoinHostPort(server.Host, server.Port))
			} else {
				i.Connect(server.Host)
			}

			go newIRCHandler(i, session).run()

			var joining []string
			for _, channel := range channels {
				if channel.Server == server.Host {
					joining = append(joining, channel.Name)
				}
			}
			i.Join(joining...)
		}
	}
}
