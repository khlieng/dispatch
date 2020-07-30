package server

import (
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/eyedeekay/goSam"
	"github.com/khlieng/dispatch/pkg/irc"
	"github.com/khlieng/dispatch/storage"
	"golang.org/x/net/proxy"
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
	ircCfg.AutoCTCP = cfg.AutoCTCP

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

	if cfg.Proxy.Enabled && strings.ToLower(cfg.Proxy.Protocol) == "socks5" {
		addr := net.JoinHostPort(cfg.Proxy.Host, cfg.Proxy.Port)

		var auth *proxy.Auth
		if cfg.Proxy.Username != "" {
			auth = &proxy.Auth{
				User:     cfg.Proxy.Username,
				Password: cfg.Proxy.Password,
			}
		}

		dialer, err := proxy.SOCKS5("tcp", addr, auth, irc.DefaultDialer)
		if err != nil {
			log.Println(err)
		} else {
			ircCfg.Dialer = dialer
		}
	}

	if cfg.Proxy.Enabled && strings.ToLower(cfg.Proxy.Protocol) == "i2p" {
		addr := net.JoinHostPort(cfg.Proxy.Host, cfg.Proxy.Port)

		client, err := goSam.NewClient(addr)
		if err != nil {
			log.Println(err)
		} else {
			ircCfg.Dialer = client
		}
	}

	i := irc.NewClient(ircCfg)
	i.Config.HandleNickInUse = createNickInUseHandler(i, state)

	state.setNetwork(network.Host, state.user.NewNetwork(network, i))
	i.Connect()
	go newIRCHandler(i, state).run()

	return i
}
