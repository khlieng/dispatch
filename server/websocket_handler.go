package server

import (
	"crypto/tls"
	"encoding/json"
	"log"
	"strings"

	"github.com/khlieng/dispatch/Godeps/_workspace/src/github.com/gorilla/websocket"

	"github.com/khlieng/dispatch/irc"
	"github.com/khlieng/dispatch/storage"
)

type wsHandler struct {
	ws       *wsConn
	session  *Session
	addr     string
	handlers map[string]func([]byte)
}

func newWSHandler(conn *websocket.Conn, uuid string) *wsHandler {
	h := &wsHandler{
		ws:   newWSConn(conn),
		addr: conn.RemoteAddr().String(),
	}
	h.init(uuid)
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

func (h *wsHandler) init(uuid string) {
	log.Println(h.addr, "set UUID", uuid)

	sessionLock.Lock()
	if storedSession, exists := sessions[uuid]; exists {
		sessionLock.Unlock()
		h.session = storedSession
		h.session.setWS(h.addr, h.ws)

		log.Println(h.addr, "attached to", h.session.numIRC(), "existing IRC connections")

		channels := h.session.user.GetChannels()
		for i, channel := range channels {
			channels[i].Topic = channelStore.GetTopic(channel.Server, channel.Name)
		}

		h.session.sendJSON("channels", channels)
		h.session.sendJSON("servers", h.session.user.GetServers())

		for _, channel := range channels {
			h.session.sendJSON("users", Userlist{
				Server:  channel.Server,
				Channel: channel.Name,
				Users:   channelStore.GetUsers(channel.Server, channel.Name),
			})
		}
	} else {
		h.session = NewSession()
		h.session.user = storage.NewUser(uuid)

		sessions[uuid] = h.session
		sessionLock.Unlock()

		h.session.setWS(h.addr, h.ws)
		h.session.sendJSON("servers", nil)

		go h.session.write()
	}
}

func (h *wsHandler) connect(b []byte) {
	var data Connect
	json.Unmarshal(b, &data)

	if _, ok := h.session.getIRC(data.Server); !ok {
		log.Println(h.addr, "connecting to", data.Server)

		i := irc.NewClient(data.Nick, data.Username)
		i.TLS = data.TLS
		i.Password = data.Password
		i.Realname = data.Realname

		if cert := h.session.user.GetCertificate(); cert != nil {
			i.TLSConfig = &tls.Config{
				Certificates: []tls.Certificate{*cert},
			}
		}

		if idx := strings.Index(data.Server, ":"); idx < 0 {
			h.session.setIRC(data.Server, i)
		} else {
			h.session.setIRC(data.Server[:idx], i)
		}

		i.Connect(data.Server)
		go newIRCHandler(i, h.session).run()

		h.session.user.AddServer(storage.Server{
			Name:     data.Name,
			Address:  i.Host,
			TLS:      data.TLS,
			Password: data.Password,
			Nick:     data.Nick,
			Username: data.Username,
			Realname: data.Realname,
		})
	} else {
		log.Println(h.addr, "already connected to", data.Server)
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
		i.Quit()
		h.session.deleteIRC(data.Server)
		channelStore.RemoveUserAll(i.GetNick(), data.Server)
		h.session.user.RemoveServer(data.Server)
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
		h.session.user.SetNick(data.New, data.Server)
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
		"search":  h.search,
		"cert":    h.cert,
	}
}
