package server

import (
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/kjk/betterguid"
)

type wsHandler struct {
	ws       *wsConn
	state    *State
	addr     net.Addr
	handlers map[string]func([]byte)
}

func newWSHandler(conn *websocket.Conn, state *State, r *http.Request) *wsHandler {
	h := &wsHandler{
		ws:    newWSConn(conn),
		state: state,
		addr:  conn.RemoteAddr(),
	}

	if r.Header.Get("X-Forwarded-For") != "" {
		ip := net.ParseIP(r.Header.Get("X-Forwarded-For"))
		if ip != nil {
			h.addr.(*net.TCPAddr).IP = ip
		}
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
			if h.state != nil {
				h.state.deleteWS(h.addr.String())
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
	h.state.setWS(h.addr.String(), h.ws)
	h.state.user.SetLastIP(addrToIPBytes(h.addr))

	log.Println(h.addr, "[State] User ID:", h.state.user.ID, "|",
		h.state.numIRC(), "IRC connections |",
		h.state.numWS(), "WebSocket connections")

	channels, err := h.state.user.GetChannels()
	if err != nil {
		log.Println(err)
	}
	path := r.URL.EscapedPath()[3:]
	pathServer, pathChannel := getTabFromPath(path)
	cookieServer, cookieChannel := parseTabCookie(r, path)

	for _, channel := range channels {
		if (channel.Server == pathServer && channel.Name == pathChannel) ||
			(channel.Server == cookieServer && channel.Name == cookieChannel) {
			// Userlist and messages for this channel gets embedded in the index page
			continue
		}

		h.state.sendJSON("users", Userlist{
			Server:  channel.Server,
			Channel: channel.Name,
			Users:   channelStore.GetUsers(channel.Server, channel.Name),
		})

		h.state.sendLastMessages(channel.Server, channel.Name, 50)
	}
}

func (h *wsHandler) connect(b []byte) {
	var data Server
	data.UnmarshalJSON(b)

	if _, ok := h.state.getIRC(data.Host); !ok {
		log.Println(h.addr, "[IRC] Add server", data.Host)

		connectIRC(data.Server, h.state, addrToIPBytes(h.addr))

		go h.state.user.AddServer(data.Server)
	} else {
		log.Println(h.addr, "[IRC]", data.Host, "already added")
	}
}

func (h *wsHandler) reconnect(b []byte) {
	var data ReconnectSettings
	data.UnmarshalJSON(b)

	if i, ok := h.state.getIRC(data.Server); ok && !i.Connected() {
		if i.TLS {
			i.TLSConfig.InsecureSkipVerify = data.SkipVerify
		}
		i.Reconnect()
	}
}

func (h *wsHandler) join(b []byte) {
	var data Join
	data.UnmarshalJSON(b)

	if i, ok := h.state.getIRC(data.Server); ok {
		i.Join(data.Channels...)
	}
}

func (h *wsHandler) part(b []byte) {
	var data Part
	data.UnmarshalJSON(b)

	if i, ok := h.state.getIRC(data.Server); ok {
		i.Part(data.Channels...)
	}
}

func (h *wsHandler) quit(b []byte) {
	var data Quit
	data.UnmarshalJSON(b)

	log.Println(h.addr, "[IRC] Remove server", data.Server)
	if i, ok := h.state.getIRC(data.Server); ok {
		h.state.deleteIRC(data.Server)
		i.Quit()
	}

	go h.state.user.RemoveServer(data.Server)
}

func (h *wsHandler) message(b []byte) {
	var data Message
	data.UnmarshalJSON(b)

	if i, ok := h.state.getIRC(data.Server); ok {
		i.Privmsg(data.To, data.Content)

		go h.state.user.LogMessage(betterguid.New(),
			data.Server, i.GetNick(), data.To, data.Content)
	}
}

func (h *wsHandler) nick(b []byte) {
	var data Nick
	data.UnmarshalJSON(b)

	if i, ok := h.state.getIRC(data.Server); ok {
		i.Nick(data.New)
	}
}

func (h *wsHandler) topic(b []byte) {
	var data Topic
	data.UnmarshalJSON(b)

	if i, ok := h.state.getIRC(data.Server); ok {
		i.Topic(data.Channel, data.Topic)
	}
}

func (h *wsHandler) invite(b []byte) {
	var data Invite
	data.UnmarshalJSON(b)

	if i, ok := h.state.getIRC(data.Server); ok {
		i.Invite(data.User, data.Channel)
	}
}

func (h *wsHandler) kick(b []byte) {
	var data Invite
	data.UnmarshalJSON(b)

	if i, ok := h.state.getIRC(data.Server); ok {
		i.Kick(data.Channel, data.User)
	}
}

func (h *wsHandler) whois(b []byte) {
	var data Whois
	data.UnmarshalJSON(b)

	if i, ok := h.state.getIRC(data.Server); ok {
		i.Whois(data.User)
	}
}

func (h *wsHandler) away(b []byte) {
	var data Away
	data.UnmarshalJSON(b)

	if i, ok := h.state.getIRC(data.Server); ok {
		i.Away(data.Message)
	}
}

func (h *wsHandler) raw(b []byte) {
	var data Raw
	data.UnmarshalJSON(b)

	if i, ok := h.state.getIRC(data.Server); ok {
		i.Write(data.Message)
	}
}

func (h *wsHandler) search(b []byte) {
	go func() {
		var data SearchRequest
		data.UnmarshalJSON(b)

		results, err := h.state.user.SearchMessages(data.Server, data.Channel, data.Phrase)
		if err != nil {
			log.Println(err)
			return
		}

		h.state.sendJSON("search", SearchResult{
			Server:  data.Server,
			Channel: data.Channel,
			Results: results,
		})
	}()
}

func (h *wsHandler) cert(b []byte) {
	var data ClientCert
	data.UnmarshalJSON(b)

	err := h.state.user.SetCertificate(data.Cert, data.Key)
	if err != nil {
		h.state.sendJSON("cert_fail", Error{Message: err.Error()})
		return
	}

	h.state.sendJSON("cert_success", nil)
}

func (h *wsHandler) fetchMessages(b []byte) {
	var data FetchMessages
	data.UnmarshalJSON(b)

	h.state.sendMessages(data.Server, data.Channel, 200, data.Next)
}

func (h *wsHandler) setServerName(b []byte) {
	var data ServerName
	data.UnmarshalJSON(b)

	if isValidServerName(data.Name) {
		h.state.user.SetServerName(data.Name, data.Server)
	}
}

func (h *wsHandler) initHandlers() {
	h.handlers = map[string]func([]byte){
		"connect":         h.connect,
		"reconnect":       h.reconnect,
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
