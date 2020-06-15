package server

import (
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/khlieng/dispatch/storage"
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
	if r.Header.Get("X-Forwarded-For") != "" {
		ip := net.ParseIP(r.Header.Get("X-Forwarded-For"))
		if ip != nil {
			h.addr.(*net.TCPAddr).IP = ip
		}
	}

	h.state.setWS(h.addr.String(), h.ws)
	h.state.user.SetLastIP(addrToIPBytes(h.addr))
	if r.TLS != nil {
		h.state.Set("scheme", "https")
	} else {
		h.state.Set("scheme", "http")
	}
	h.state.Set("host", r.Host)

	log.Println(h.addr, "[State] User ID:", h.state.user.ID, "|",
		h.state.numIRC(), "IRC connections |",
		h.state.numWS(), "WebSocket connections")

	go h.sendData(r)
}

func (h *wsHandler) sendData(r *http.Request) {
	tab, err := tabFromRequest(r)
	if err != nil {
		log.Println(err)
	}

	h.state.lock.Lock()
	for _, network := range h.state.networks {
		for _, channel := range network.ChannelNames() {
			if network.Host == tab.Network && channel == tab.Name {
				// Userlist and messages for this channel gets embedded in the index page
				continue
			}

			if users := network.Client().ChannelUsers(channel); len(users) > 0 {
				h.state.sendJSON("users", Userlist{
					Network: network.Host,
					Channel: channel,
					Users:   users,
				})
			}

			h.state.sendLastMessages(network.Host, channel, 50)
		}

	}
	h.state.lock.Unlock()

	openDMs, err := h.state.user.OpenDMs()
	if err != nil {
		log.Println(err)
	}

	for _, openDM := range openDMs {
		if openDM.Network == tab.Network && openDM.Name == tab.Name {
			continue
		}

		h.state.sendLastMessages(openDM.Network, openDM.Name, 50)
	}
}

func (h *wsHandler) connect(b []byte) {
	var network storage.Network
	network.UnmarshalJSON(b)

	network.Host = strings.ToLower(network.Host)

	if _, ok := h.state.network(network.Host); !ok {
		log.Println(h.addr, "[IRC] Add server", network.Host)

		connectIRC(&network, h.state, addrToIPBytes(h.addr))

		go network.Save()
	} else {
		log.Println(h.addr, "[IRC]", network.Host, "already added")
	}
}

func (h *wsHandler) reconnect(b []byte) {
	var data ReconnectSettings
	data.UnmarshalJSON(b)

	if i, ok := h.state.client(data.Network); ok && !i.Connected() {
		if i.Config.TLS {
			i.Config.TLSConfig.InsecureSkipVerify = data.SkipVerify
		}
		i.Reconnect()
	}
}

func (h *wsHandler) join(b []byte) {
	var data Join
	data.UnmarshalJSON(b)

	if i, ok := h.state.client(data.Network); ok {
		i.Join(data.Channels...)
	}
}

func (h *wsHandler) part(b []byte) {
	var data Part
	data.UnmarshalJSON(b)

	if i, ok := h.state.client(data.Network); ok {
		i.Part(data.Channels...)
	}

	go func() {
		if network, ok := h.state.network(data.Network); ok {
			network.RemoveChannels(data.Channels...)
		}

		for _, channel := range data.Channels {
			h.state.user.RemoveChannel(data.Network, channel)
		}
	}()
}

func (h *wsHandler) quit(b []byte) {
	var data Quit
	data.UnmarshalJSON(b)

	log.Println(h.addr, "[IRC] Remove server", data.Network)
	if i, ok := h.state.client(data.Network); ok {
		h.state.deleteNetwork(data.Network)
		i.Quit()
	}

	go h.state.user.RemoveNetwork(data.Network)
}

func (h *wsHandler) message(b []byte) {
	var data Message
	data.UnmarshalJSON(b)

	if i, ok := h.state.client(data.Network); ok {
		i.Privmsg(data.To, data.Content)

		go h.state.user.LogMessage(&storage.Message{
			Network: data.Network,
			From:    i.GetNick(),
			To:      data.To,
			Content: data.Content,
		})
	}
}

func (h *wsHandler) nick(b []byte) {
	var data Nick
	data.UnmarshalJSON(b)

	if i, ok := h.state.client(data.Network); ok {
		i.Nick(data.New)
	}
}

func (h *wsHandler) topic(b []byte) {
	var data Topic
	data.UnmarshalJSON(b)

	if i, ok := h.state.client(data.Network); ok {
		i.Topic(data.Channel, data.Topic)
	}
}

func (h *wsHandler) invite(b []byte) {
	var data Invite
	data.UnmarshalJSON(b)

	if i, ok := h.state.client(data.Network); ok {
		i.Invite(data.User, data.Channel)
	}
}

func (h *wsHandler) kick(b []byte) {
	var data Invite
	data.UnmarshalJSON(b)

	if i, ok := h.state.client(data.Network); ok {
		i.Kick(data.Channel, data.User)
	}
}

func (h *wsHandler) whois(b []byte) {
	var data Whois
	data.UnmarshalJSON(b)

	if i, ok := h.state.client(data.Network); ok {
		i.Whois(data.User)
	}
}

func (h *wsHandler) away(b []byte) {
	var data Away
	data.UnmarshalJSON(b)

	if i, ok := h.state.client(data.Network); ok {
		i.Away(data.Message)
	}
}

func (h *wsHandler) raw(b []byte) {
	var data Raw
	data.UnmarshalJSON(b)

	if i, ok := h.state.client(data.Network); ok {
		i.Write(data.Message)
	}
}

func (h *wsHandler) search(b []byte) {
	go func() {
		var data SearchRequest
		data.UnmarshalJSON(b)

		results, err := h.state.user.SearchMessages(data.Network, data.Channel, data.Phrase)
		if err != nil {
			log.Println(err)
			return
		}

		h.state.sendJSON("search", SearchResult{
			Network: data.Network,
			Channel: data.Channel,
			Results: results,
		})
	}()
}

func (h *wsHandler) cert(b []byte) {
	var data ClientCert
	data.UnmarshalJSON(b)

	err := h.state.user.SetCertificate([]byte(data.Cert), []byte(data.Key))
	if err != nil {
		h.state.sendJSON("cert_fail", Error{Message: err.Error()})
		return
	}

	h.state.sendJSON("cert_success", nil)
}

func (h *wsHandler) fetchMessages(b []byte) {
	var data FetchMessages
	data.UnmarshalJSON(b)

	h.state.sendMessages(data.Network, data.Channel, 200, data.Next)
}

func (h *wsHandler) setNetworkName(b []byte) {
	var data NetworkName
	data.UnmarshalJSON(b)

	if isValidNetworkName(data.Name) {
		if network, ok := h.state.network(data.Network); ok {
			network.SetName(data.Name)
			go network.Save()
		}
	}
}

func (h *wsHandler) setSettings(b []byte) {
	err := h.state.user.UnmarshalClientSettingsJSON(b)
	if err != nil {
		log.Println(err)
	}
}

func (h *wsHandler) channelSearch(b []byte) {
	var data ChannelSearch
	data.UnmarshalJSON(b)

	index, needsUpdate := channelIndexes.Get(data.Network)
	if index != nil {
		n := 10
		if data.Start > 0 {
			n = 50
		}

		h.state.sendJSON("channel_search", ChannelSearchResult{
			ChannelSearch: data,
			Results:       index.SearchN(data.Q, data.Start, n),
		})
	}

	if i, ok := h.state.client(data.Network); ok && needsUpdate {
		h.state.Set("update_chanlist_"+data.Network, true)
		i.List()
	}
}

func (h *wsHandler) openDM(b []byte) {
	var data Tab
	data.UnmarshalJSON(b)

	h.state.sendLastMessages(data.Network, data.Name, 50)
	h.state.user.AddOpenDM(data.Network, data.Name)
}

func (h *wsHandler) closeDM(b []byte) {
	var data Tab
	data.UnmarshalJSON(b)

	h.state.user.RemoveOpenDM(data.Network, data.Name)
}

func (h *wsHandler) initHandlers() {
	h.handlers = map[string]func([]byte){
		"connect":          h.connect,
		"reconnect":        h.reconnect,
		"join":             h.join,
		"part":             h.part,
		"quit":             h.quit,
		"message":          h.message,
		"nick":             h.nick,
		"topic":            h.topic,
		"invite":           h.invite,
		"kick":             h.kick,
		"whois":            h.whois,
		"away":             h.away,
		"raw":              h.raw,
		"search":           h.search,
		"cert":             h.cert,
		"fetch_messages":   h.fetchMessages,
		"set_network_name": h.setNetworkName,
		"settings_set":     h.setSettings,
		"channel_search":   h.channelSearch,
		"open_dm":          h.openDM,
		"close_dm":         h.closeDM,
	}
}

func isValidNetworkName(name string) bool {
	return strings.TrimSpace(name) != ""
}
