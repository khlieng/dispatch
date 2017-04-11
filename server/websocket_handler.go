package server

import (
	"crypto/tls"
	"encoding/json"
	"log"
	"net"

	"github.com/gorilla/websocket"
	"github.com/spf13/viper"

	"github.com/khlieng/dispatch/irc"
	"github.com/khlieng/dispatch/storage"
)

type wsHandler struct {
	ws       *wsConn
	session  *Session
	addr     string
	handlers map[string]func([]byte)
}

func newWSHandler(conn *websocket.Conn, session *Session) *wsHandler {
	h := &wsHandler{
		ws:      newWSConn(conn),
		session: session,
		addr:    conn.RemoteAddr().String(),
	}
	h.init()
	h.initHandlers()
	return h
}

func (h *wsHandler) run() {
	defer h.ws.close()
	go h.ws.send()
	go h.ws.recv()

	for {
		req, ok := <-h.ws.in
		if !ok {
			if h.session != nil {
				h.session.deleteWS(h.addr)
			}
			return
		}

		h.dispatchRequest(req)
	}
}

func (h *wsHandler) dispatchRequest(req WSRequest) {
	if handler, ok := h.handlers[req.Type]; ok {
		handler(req.Data)
	}
}

func (h *wsHandler) init() {
	h.session.setWS(h.addr, h.ws)

	log.Println(h.addr, "[Session] User ID:", h.session.user.ID, "|",
		h.session.numIRC(), "IRC connections |",
		h.session.numWS(), "WebSocket connections")

	channels := h.session.user.GetChannels()

	for _, channel := range channels {
		h.session.sendJSON("users", Userlist{
			Server:  channel.Server,
			Channel: channel.Name,
			Users:   channelStore.GetUsers(channel.Server, channel.Name),
		})
	}
}

func (h *wsHandler) connect(b []byte) {
	var data Connect
	json.Unmarshal(b, &data)

	host, port, err := net.SplitHostPort(data.Server)
	if err != nil {
		host = data.Server
	}

	if _, ok := h.session.getIRC(host); !ok {
		log.Println(h.addr, "[IRC] Add server", data.Server)

		i := irc.NewClient(data.Nick, data.Username)
		i.TLS = data.TLS
		i.Realname = data.Realname
		i.HandleNickInUse = createNickInUseHandler(i, h.session)

		if data.Password == "" &&
			viper.GetString("defaults.password") != "" &&
			data.Server == viper.GetString("defaults.address") {
			i.Password = viper.GetString("defaults.password")
		} else {
			i.Password = data.Password
		}

		if cert := h.session.user.GetCertificate(); cert != nil {
			i.TLSConfig = &tls.Config{
				Certificates:       []tls.Certificate{*cert},
				InsecureSkipVerify: !viper.GetBool("verify_client_certificates"),
			}
		}

		h.session.setIRC(host, i)
		i.Connect(data.Server)
		go newIRCHandler(i, h.session).run()

		go h.session.user.AddServer(storage.Server{
			Name:     data.Name,
			Host:     host,
			Port:     port,
			TLS:      data.TLS,
			Password: data.Password,
			Nick:     data.Nick,
			Username: data.Username,
			Realname: data.Realname,
		})
	} else {
		log.Println(h.addr, "[IRC]", data.Server, "already added")
	}
}

func (h *wsHandler) join(b []byte) {
	var data Join
	json.Unmarshal(b, &data)

	if i, ok := h.session.getIRC(data.Server); ok {
		i.Join(data.Channels...)
	}
}

func (h *wsHandler) part(b []byte) {
	var data Part
	json.Unmarshal(b, &data)

	if i, ok := h.session.getIRC(data.Server); ok {
		i.Part(data.Channels...)
	}
}

func (h *wsHandler) quit(b []byte) {
	var data Quit
	json.Unmarshal(b, &data)

	if i, ok := h.session.getIRC(data.Server); ok {
		log.Println(h.addr, "[IRC] Remove server", data.Server)

		i.Quit()
		h.session.deleteIRC(data.Server)
		channelStore.RemoveUserAll(i.GetNick(), data.Server)
		go h.session.user.RemoveServer(data.Server)
	}
}

func (h *wsHandler) chat(b []byte) {
	var data Chat
	json.Unmarshal(b, &data)

	if i, ok := h.session.getIRC(data.Server); ok {
		i.Privmsg(data.To, data.Message)
	}
}

func (h *wsHandler) nick(b []byte) {
	var data Nick
	json.Unmarshal(b, &data)

	if i, ok := h.session.getIRC(data.Server); ok {
		i.Nick(data.New)
	}
}

func (h *wsHandler) invite(b []byte) {
	var data Invite
	json.Unmarshal(b, &data)

	if i, ok := h.session.getIRC(data.Server); ok {
		i.Invite(data.User, data.Channel)
	}
}

func (h *wsHandler) kick(b []byte) {
	var data Invite
	json.Unmarshal(b, &data)

	if i, ok := h.session.getIRC(data.Server); ok {
		i.Kick(data.Channel, data.User)
	}
}

func (h *wsHandler) whois(b []byte) {
	var data Whois
	json.Unmarshal(b, &data)

	if i, ok := h.session.getIRC(data.Server); ok {
		i.Whois(data.User)
	}
}

func (h *wsHandler) away(b []byte) {
	var data Away
	json.Unmarshal(b, &data)

	if i, ok := h.session.getIRC(data.Server); ok {
		i.Away(data.Message)
	}
}

func (h *wsHandler) raw(b []byte) {
	var data Raw
	json.Unmarshal(b, &data)

	if i, ok := h.session.getIRC(data.Server); ok {
		i.Write(data.Message)
	}
}

func (h *wsHandler) search(b []byte) {
	go func() {
		var data SearchRequest
		json.Unmarshal(b, &data)

		results, err := h.session.user.SearchMessages(data.Server, data.Channel, data.Phrase)
		if err != nil {
			log.Println(err)
			return
		}

		h.session.sendJSON("search", SearchResult{
			Server:  data.Server,
			Channel: data.Channel,
			Results: results,
		})
	}()
}

func (h *wsHandler) cert(b []byte) {
	var data ClientCert
	json.Unmarshal(b, &data)

	err := h.session.user.SetCertificate(data.Cert, data.Key)
	if err != nil {
		h.session.sendJSON("cert_fail", Error{Message: err.Error()})
		return
	}

	h.session.sendJSON("cert_success", nil)
}

func (h *wsHandler) initHandlers() {
	h.handlers = map[string]func([]byte){
		"connect": h.connect,
		"join":    h.join,
		"part":    h.part,
		"quit":    h.quit,
		"chat":    h.chat,
		"nick":    h.nick,
		"invite":  h.invite,
		"kick":    h.kick,
		"whois":   h.whois,
		"away":    h.away,
		"raw":     h.raw,
		"search":  h.search,
		"cert":    h.cert,
	}
}
