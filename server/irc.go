package server

import (
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"net"

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

		state.sendJSON("error", IRCError{
			Server:  i.Host,
			Message: fmt.Sprintf("Nickname %s is already in use, using %s instead", nick, newNick),
		})

		return newNick
	}
}

func connectIRC(server *storage.Server, state *State, srcIP []byte) *irc.Client {
	i := irc.NewClient(server.Nick, server.Username)
	i.TLS = server.TLS
	i.Realname = server.Realname
	i.HandleNickInUse = createNickInUseHandler(i, state)

	address := server.Host
	if server.Port != "" {
		address = net.JoinHostPort(server.Host, server.Port)
	}

	cfg := state.srv.Config()

	if cfg.HexIP {
		i.Username = hex.EncodeToString(srcIP)
	} else if i.Username == "" {
		i.Username = server.Nick
	}

	if i.Realname == "" {
		i.Realname = server.Nick
	}

	if server.Password == "" &&
		cfg.Defaults.Password != "" &&
		address == cfg.Defaults.Host {
		i.Password = cfg.Defaults.Password
	} else {
		i.Password = server.Password
	}

	if i.TLS {
		i.TLSConfig = &tls.Config{
			InsecureSkipVerify: !cfg.VerifyCertificates,
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
