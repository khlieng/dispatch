package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/kjk/betterguid"
)

type wsHandler struct {
	ws       *wsConn
	session  *Session
	addr     string
	handlers map[string]func([]byte)
}

func newWSHandler(conn *websocket.Conn, session *Session, r *http.Request) *wsHandler {
	h := &wsHandler{
		ws:      newWSConn(conn),
		session: session,
		addr:    conn.RemoteAddr().String(),
	}
	h.init(r)
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

func (h *wsHandler) init(r *http.Request) {
	h.session.setWS(h.addr, h.ws)

	log.Println(h.addr, "[Session] User ID:", h.session.user.ID, "|",
		h.session.numIRC(), "IRC connections |",
		h.session.numWS(), "WebSocket connections")

	channels := h.session.user.GetChannels()
	path := r.URL.EscapedPath()[3:]
	pathServer, pathChannel := getTabFromPath(path)
	cookieServer, cookieChannel := parseTabCookie(r, path)

	for _, channel := range channels {
		if (channel.Server == pathServer && channel.Name == pathChannel) ||
			(channel.Server == cookieServer && channel.Name == cookieChannel) {
			continue
		}

		h.session.sendJSON("users", Userlist{
			Server:  channel.Server,
			Channel: channel.Name,
			Users:   channelStore.GetUsers(channel.Server, channel.Name),
		})

		h.session.sendLastMessages(channel.Server, channel.Name, 50)
	}
}

func (h *wsHandler) connect(b []byte) {
	var data Server
	json.Unmarshal(b, &data)

	if _, ok := h.session.getIRC(data.Host); !ok {
		log.Println(h.addr, "[IRC] Add server", data.Server)

		connectIRC(data.Server, h.session)

		go h.session.user.AddServer(data.Server)
	} else {
		log.Println(h.addr, "[IRC]", data.Host, "already added")
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

func (h *wsHandler) message(b []byte) {
	var data Message
	json.Unmarshal(b, &data)

	if i, ok := h.session.getIRC(data.Server); ok {
		i.Privmsg(data.To, data.Content)

		go h.session.user.LogMessage(betterguid.New(),
			data.Server, i.GetNick(), data.To, data.Content)
	}
}

func (h *wsHandler) nick(b []byte) {
	var data Nick
	json.Unmarshal(b, &data)

	if i, ok := h.session.getIRC(data.Server); ok {
		i.Nick(data.New)
	}
}

func (h *wsHandler) topic(b []byte) {
	var data Topic
	json.Unmarshal(b, &data)

	if i, ok := h.session.getIRC(data.Server); ok {
		i.Topic(data.Channel, data.Topic)
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

func (h *wsHandler) fetchMessages(b []byte) {
	var data FetchMessages
	json.Unmarshal(b, &data)

	h.session.sendMessages(data.Server, data.Channel, 200, data.Next)
}

func (h *wsHandler) setServerName(b []byte) {
	var data ServerName
	json.Unmarshal(b, &data)

	if isValidServerName(data.Name) {
		h.session.user.SetServerName(data.Name, data.Server)
	}
}

func (h *wsHandler) initHandlers() {
	h.handlers = map[string]func([]byte){
		"connect":         h.connect,
		"join":            h.join,
		"part":            h.part,
		"quit":            h.quit,
		"message":         h.message,
		"nick":            h.nick,
		"topic":           h.topic,
		"invite":          h.invite,
		"kick":            h.kick,
		"whois":           h.whois,
		"away":            h.away,
		"raw":             h.raw,
		"search":          h.search,
		"cert":            h.cert,
		"fetch_messages":  h.fetchMessages,
		"set_server_name": h.setServerName,
	}
}

func isValidServerName(name string) bool {
	return strings.TrimSpace(name) != ""
}
