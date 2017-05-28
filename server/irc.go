package server

import (
	"crypto/tls"
	"net"

	"github.com/spf13/viper"

	"github.com/khlieng/dispatch/irc"
	"github.com/khlieng/dispatch/storage"
)

func createNickInUseHandler(i *irc.Client, session *Session) func(string) string {
	return func(nick string) string {
		newNick := nick + "_"
		session.printError("Nickname", nick, "is already in use, using", newNick, "instead")

		return newNick
	}
}

func reconnectIRC() {
	for _, user := range storage.LoadUsers() {
		session := NewSession(user)
		sessions.set(user.ID, session)
		go session.run()

		channels := user.GetChannels()

		for _, server := range user.GetServers() {
			i := irc.NewClient(server.Nick, server.Username)
			i.TLS = server.TLS
			i.Password = server.Password
			i.Realname = server.Realname
			i.HandleNickInUse = createNickInUseHandler(i, session)

			if cert := user.GetCertificate(); cert != nil {
				i.TLSConfig = &tls.Config{
					Certificates:       []tls.Certificate{*cert},
					InsecureSkipVerify: !viper.GetBool("verify_client_certificates"),
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
