package server

import (
	"crypto/tls"
	"net"

	"github.com/khlieng/dispatch/irc"
	"github.com/khlieng/dispatch/storage"
	"github.com/spf13/viper"
)

func createNickInUseHandler(i *irc.Client, session *Session) func(string) string {
	return func(nick string) string {
		newNick := nick + "_"

		if newNick == i.GetNick() {
			session.sendJSON("nick_fail", NickFail{
				Server: i.Host,
			})
		}

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
			i := connectIRC(server, session)

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

func connectIRC(server storage.Server, session *Session) *irc.Client {
	i := irc.NewClient(server.Nick, server.Username)
	i.TLS = server.TLS
	i.Realname = server.Realname
	i.HandleNickInUse = createNickInUseHandler(i, session)

	address := server.Host
	if server.Port != "" {
		address = net.JoinHostPort(server.Host, server.Port)
	}

	if server.Password == "" &&
		viper.GetString("defaults.password") != "" &&
		address == viper.GetString("defaults.address") {
		i.Password = viper.GetString("defaults.password")
	} else {
		i.Password = server.Password
	}

	if i.TLS {
		i.TLSConfig = &tls.Config{
			InsecureSkipVerify: !viper.GetBool("verify_certificates"),
		}

		if cert := session.user.GetCertificate(); cert != nil {
			i.TLSConfig.Certificates = []tls.Certificate{*cert}
		}
	}

	session.setIRC(server.Host, i)
	i.Connect(address)
	go newIRCHandler(i, session).run()

	return i
}
