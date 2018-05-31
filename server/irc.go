package server

import (
	"crypto/tls"
	"net"

	"github.com/spf13/viper"

	"github.com/khlieng/dispatch/pkg/irc"
	"github.com/khlieng/dispatch/storage"
)

func createNickInUseHandler(i *irc.Client, state *State) func(string) string {
	return func(nick string) string {
		newNick := nick + "_"

		if newNick == i.GetNick() {
			state.sendJSON("nick_fail", NickFail{
				Server: i.Host,
			})
		}

		state.printError("Nickname", nick, "is already in use, using", newNick, "instead")

		return newNick
	}
}

func connectIRC(server *storage.Server, state *State) *irc.Client {
	i := irc.NewClient(server.Nick, server.Username)
	i.TLS = server.TLS
	i.Realname = server.Realname
	i.HandleNickInUse = createNickInUseHandler(i, state)

	address := server.Host
	if server.Port != "" {
		address = net.JoinHostPort(server.Host, server.Port)
	}

	if i.Username == "" {
		i.Username = server.Nick
	}
	if i.Realname == "" {
		i.Realname = server.Nick
	}

	if server.Password == "" &&
		viper.GetString("defaults.password") != "" &&
		address == viper.GetString("defaults.host") {
		i.Password = viper.GetString("defaults.password")
	} else {
		i.Password = server.Password
	}

	if i.TLS {
		i.TLSConfig = &tls.Config{
			InsecureSkipVerify: !viper.GetBool("verify_certificates"),
		}

		if cert := state.user.GetCertificate(); cert != nil {
			i.TLSConfig.Certificates = []tls.Certificate{*cert}
		}
	}

	state.setIRC(server.Host, i)
	i.Connect(address)
	go newIRCHandler(i, state).run()

	return i
}
