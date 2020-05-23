package server

import (
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"net"

	"github.com/khlieng/dispatch/pkg/irc"
	"github.com/khlieng/dispatch/storage"
	"github.com/khlieng/dispatch/version"
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
			Message: fmt.Sprintf("Nickname %s is unavailable, trying %s instead", nick, newNick),
		})

		return newNick
	}
}

func connectIRC(server *storage.Server, state *State, srcIP []byte) *irc.Client {
	i := irc.NewClient(server.Nick, server.Username)
	i.TLS = server.TLS
	i.Realname = server.Realname
	i.Version = fmt.Sprintf("Dispatch %s (git: %s)", version.Tag, version.Commit)
	i.Source = "https://github.com/khlieng/dispatch"
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

	if server.ServerPassword == "" &&
		cfg.Defaults.ServerPassword != "" &&
		address == cfg.Defaults.Host {
		i.Password = cfg.Defaults.ServerPassword
	} else {
		i.Password = server.ServerPassword
	}

	if server.Account != "" && server.Password != "" {
		i.SASL = &irc.SASLPlain{
			Username: server.Account,
			Password: server.Password,
		}
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
