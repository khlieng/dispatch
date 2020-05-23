package server

import (
	"crypto/tls"
	"encoding/hex"
	"fmt"

	"github.com/khlieng/dispatch/pkg/irc"
	"github.com/khlieng/dispatch/storage"
	"github.com/khlieng/dispatch/version"
)

func createNickInUseHandler(i *irc.Client, state *State) func(string) string {
	return func(nick string) string {
		newNick := nick + "_"

		if newNick == i.GetNick() {
			state.sendJSON("nick_fail", NickFail{
				Server: i.Host(),
			})
		}

		state.sendJSON("error", IRCError{
			Server:  i.Host(),
			Message: fmt.Sprintf("Nickname %s is unavailable, trying %s instead", nick, newNick),
		})

		return newNick
	}
}

func connectIRC(server *storage.Server, state *State, srcIP []byte) *irc.Client {
	cfg := state.srv.Config()

	ircCfg := irc.Config{
		Host:     server.Host,
		Port:     server.Port,
		TLS:      server.TLS,
		Nick:     server.Nick,
		Username: server.Username,
		Realname: server.Realname,
		Version:  fmt.Sprintf("Dispatch %s (git: %s)", version.Tag, version.Commit),
		Source:   "https://github.com/khlieng/dispatch",
	}

	if server.TLS {
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

	if server.Account != "" && server.Password != "" {
		ircCfg.SASL = &irc.SASLPlain{
			Username: server.Account,
			Password: server.Password,
		}
	}

	if server.ServerPassword == "" &&
		cfg.Defaults.ServerPassword != "" &&
		server.Host == cfg.Defaults.Host {
		ircCfg.Password = cfg.Defaults.ServerPassword
	} else {
		ircCfg.Password = server.ServerPassword
	}

	i := irc.NewClient(ircCfg)
	i.Config.HandleNickInUse = createNickInUseHandler(i, state)

	state.setIRC(server.Host, i)
	i.Connect()
	go newIRCHandler(i, state).run()

	return i
}
