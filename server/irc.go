package server

import (
	"crypto/tls"
	"encoding/hex"
	"fmt"

	"github.com/khlieng/dispatch/pkg/irc"
	"github.com/khlieng/dispatch/storage"
)

func createNickInUseHandler(i *irc.Client, state *State) func(string) string {
	return func(nick string) string {
		newNick := nick + "_"

		if newNick == i.GetNick() {
			state.sendJSON("nick_fail", NickFail{
				Network: i.Host(),
			})
		}

		state.sendJSON("error", IRCError{
			Network: i.Host(),
			Message: fmt.Sprintf("Nickname %s is unavailable, trying %s instead", nick, newNick),
		})

		return newNick
	}
}

func connectIRC(network *storage.Network, state *State, srcIP []byte) *irc.Client {
	cfg := state.srv.Config()
	ircCfg := network.IRCConfig()

	if ircCfg.TLS {
		ircCfg.TLSConfig = &tls.Config{
			InsecureSkipVerify: !cfg.VerifyCertificates,
		}

		if cert := state.user.GetCertificate(); cert != nil {
			ircCfg.TLSConfig.Certificates = []tls.Certificate{*cert}
		}
	}

	if cfg.HexIP {
		ircCfg.Username = hex.EncodeToString(srcIP)
	}

	if ircCfg.ServerPassword == "" &&
		cfg.Defaults.ServerPassword != "" &&
		ircCfg.Host == cfg.Defaults.Host {
		ircCfg.ServerPassword = cfg.Defaults.ServerPassword
	}

	i := irc.NewClient(ircCfg)
	i.Config.HandleNickInUse = createNickInUseHandler(i, state)

	state.setNetwork(network.Host, state.user.NewNetwork(network, i))
	i.Connect()
	go newIRCHandler(i, state).run()

	return i
}
