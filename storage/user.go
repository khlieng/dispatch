package storage

import (
	"crypto/tls"
	"os"
	"sync"
	"time"

	"github.com/kjk/betterguid"
)

type User struct {
	ID       uint64
	IDBytes  []byte
	Username string

	store          Store
	messageLog     MessageStore
	messageIndex   MessageSearchProvider
	lastMessages   map[string]map[string]*Message
	clientSettings *ClientSettings
	lastIP         []byte
	certificate    *tls.Certificate
	lock           sync.Mutex
}

func NewUser(store Store) (*User, error) {
	user := &User{
		store:          store,
		clientSettings: DefaultClientSettings(),
		lastMessages:   map[string]map[string]*Message{},
	}

	err := store.SaveUser(user)
	if err != nil {
		return nil, err
	}

	user.messageLog, err = GetMessageStore(user)
	if err != nil {
		return nil, err
	}
	user.messageIndex, err = GetMessageSearchProvider(user)
	if err != nil {
		return nil, err
	}

	err = os.MkdirAll(Path.User(user.Username), 0700)
	if err != nil {
		return nil, err
	}
	err = os.Mkdir(Path.Downloads(user.Username), 0700)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func LoadUsers(store Store) ([]*User, error) {
	users, err := store.GetUsers()
	if err != nil {
		return nil, err
	}

	for _, user := range users {
		user.store = store
		user.messageLog, err = GetMessageStore(user)
		if err != nil {
			return nil, err
		}
		user.messageIndex, err = GetMessageSearchProvider(user)
		if err != nil {
			return nil, err
		}
		user.lastMessages = map[string]map[string]*Message{}
		user.loadCertificate()

		channels, err := user.GetChannels()
		if err != nil {
			return nil, err
		}

		for _, channel := range channels {
			messages, _, err := user.GetLastMessages(channel.Server, channel.Name, 1)
			if err == nil && len(messages) == 1 {
				user.lastMessages[channel.Server] = map[string]*Message{
					channel.Name: &messages[0],
				}
			}
		}
	}

	return users, nil
}

func (u *User) Remove() {
	u.store.DeleteUser(u)
	if u.messageLog != nil {
		u.messageLog.Close()
	}
	if u.messageIndex != nil {
		u.messageIndex.Close()
	}
	os.RemoveAll(Path.User(u.Username))
}

func (u *User) GetLastIP() []byte {
	u.lock.Lock()
	ip := u.lastIP
	u.lock.Unlock()
	return ip
}

func (u *User) SetLastIP(ip []byte) error {
	u.lock.Lock()
	u.lastIP = ip
	u.lock.Unlock()

	return u.store.SaveUser(u)
}

//easyjson:json
type ClientSettings struct {
	ColoredNicks bool
}

func DefaultClientSettings() *ClientSettings {
	return &ClientSettings{
		ColoredNicks: true,
	}
}

func (u *User) GetClientSettings() *ClientSettings {
	u.lock.Lock()
	settings := *u.clientSettings
	u.lock.Unlock()
	return &settings
}

func (u *User) SetClientSettings(settings *ClientSettings) error {
	u.lock.Lock()
	u.clientSettings = settings
	u.lock.Unlock()

	return u.store.SaveUser(u)
}

func (u *User) UnmarshalClientSettingsJSON(b []byte) error {
	u.lock.Lock()
	err := u.clientSettings.UnmarshalJSON(b)
	u.lock.Unlock()

	if err != nil {
		return err
	}

	return u.store.SaveUser(u)
}

type Server struct {
	Name           string
	Host           string
	Port           string
	TLS            bool
	ServerPassword string
	Nick           string
	Username       string
	Realname       string
	Account        string
	Password       string
}

func (u *User) GetServer(address string) (*Server, error) {
	return u.store.GetServer(u, address)
}

func (u *User) GetServers() ([]*Server, error) {
	return u.store.GetServers(u)
}

func (u *User) AddServer(server *Server) error {
	return u.store.SaveServer(u, server)
}

func (u *User) RemoveServer(address string) error {
	return u.store.RemoveServer(u, address)
}

func (u *User) SetNick(nick, address string) error {
	server, err := u.GetServer(address)
	if err != nil {
		return err
	}
	server.Nick = nick
	return u.AddServer(server)
}

func (u *User) SetServerName(name, address string) error {
	server, err := u.GetServer(address)
	if err != nil {
		return err
	}
	server.Name = name
	return u.AddServer(server)
}

type Channel struct {
	Server string
	Name   string
	Topic  string
}

func (u *User) GetChannels() ([]*Channel, error) {
	return u.store.GetChannels(u)
}

func (u *User) AddChannel(channel *Channel) error {
	return u.store.AddChannel(u, channel)
}

func (u *User) RemoveChannel(server, channel string) error {
	return u.store.RemoveChannel(u, server, channel)
}

type Tab struct {
	Server string
	Name   string
}

func (u *User) GetOpenDMs() ([]Tab, error) {
	return u.store.GetOpenDMs(u)
}

func (u *User) AddOpenDM(server, nick string) error {
	return u.store.AddOpenDM(u, server, nick)
}

func (u *User) RemoveOpenDM(server, nick string) error {
	return u.store.RemoveOpenDM(u, server, nick)
}

type Message struct {
	ID      string  `json:"-" bleve:"-"`
	Server  string  `json:"-" bleve:"server"`
	From    string  `bleve:"-"`
	To      string  `json:"-" bleve:"to"`
	Content string  `bleve:"content"`
	Time    int64   `bleve:"-"`
	Events  []Event `bleve:"-"`
}

func (m Message) Type() string {
	return "message"
}

func (u *User) LogMessage(msg *Message) error {
	if msg.Time == 0 {
		msg.Time = time.Now().Unix()
	}

	if msg.ID == "" {
		msg.ID = betterguid.New()
	}

	if msg.To == "" {
		msg.To = msg.From
	}

	u.setLastMessage(msg.Server, msg.To, msg)

	err := u.messageLog.LogMessage(msg)
	if err != nil {
		return err
	}
	return u.messageIndex.Index(msg.ID, msg)
}

type Event struct {
	Type   string
	Params []string
	Time   int64
}

func (u *User) LogEvent(server, name string, params []string, channels ...string) error {
	now := time.Now().Unix()
	event := Event{
		Type:   name,
		Params: params,
		Time:   now,
	}

	for _, channel := range channels {
		lastMessage := u.getLastMessage(server, channel)

		if lastMessage != nil && shouldCollapse(lastMessage, event) {
			lastMessage.Events = append(lastMessage.Events, event)
			u.setLastMessage(server, channel, lastMessage)

			err := u.messageLog.LogMessage(lastMessage)
			if err != nil {
				return err
			}
		} else {
			msg := &Message{
				ID:     betterguid.New(),
				Server: server,
				To:     channel,
				Time:   now,
				Events: []Event{event},
			}
			u.setLastMessage(server, channel, msg)

			err := u.messageLog.LogMessage(msg)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

var collapsed = []string{"join", "part", "quit"}

func shouldCollapse(msg *Message, event Event) bool {
	matches := 0
	if len(msg.Events) > 0 {
		for _, collapseType := range collapsed {
			if msg.Events[0].Type == collapseType {
				matches++
			}
			if event.Type == collapseType {
				matches++
			}
		}
	}
	return matches == 2
}

func (u *User) getLastMessage(server, channel string) *Message {
	u.lock.Lock()
	defer u.lock.Unlock()

	if _, ok := u.lastMessages[server]; !ok {
		return nil
	}

	last := u.lastMessages[server][channel]
	if last != nil {
		msg := *last
		return &msg
	}
	return nil
}

func (u *User) setLastMessage(server, channel string, msg *Message) {
	u.lock.Lock()

	if _, ok := u.lastMessages[server]; !ok {
		u.lastMessages[server] = map[string]*Message{}
	}

	u.lastMessages[server][channel] = msg
	u.lock.Unlock()
}

func (u *User) GetMessages(server, channel string, count int, fromID string) ([]Message, bool, error) {
	return u.messageLog.GetMessages(server, channel, count, fromID)
}

func (u *User) GetLastMessages(server, channel string, count int) ([]Message, bool, error) {
	return u.GetMessages(server, channel, count, "")
}

func (u *User) SearchMessages(server, channel, q string) ([]Message, error) {
	ids, err := u.messageIndex.SearchMessages(server, channel, q)
	if err != nil {
		return nil, err
	}

	return u.messageLog.GetMessagesByID(server, channel, ids)
}
