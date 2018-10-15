package storage

import (
	"crypto/tls"
	"os"
	"sync"
	"time"
)

type User struct {
	ID       uint64
	IDBytes  []byte
	Username string

	store          Store
	messageLog     MessageStore
	messageIndex   MessageSearchProvider
	clientSettings *ClientSettings
	lastIP         []byte
	certificate    *tls.Certificate
	lock           sync.Mutex
}

func NewUser(store Store) (*User, error) {
	user := &User{
		store:          store,
		clientSettings: DefaultClientSettings(),
	}

	err := store.SaveUser(user)
	if err != nil {
		return nil, err
	}

	err = os.MkdirAll(Path.User(user.Username), 0700)
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
		user.loadCertificate()
	}

	return users, nil
}

func (u *User) SetMessageStore(store MessageStore) {
	u.messageLog = store
}

func (u *User) SetMessageSearchProvider(search MessageSearchProvider) {
	u.messageIndex = search
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
	Name     string
	Host     string
	Port     string
	TLS      bool
	Password string
	Nick     string
	Username string
	Realname string
}

func (u *User) GetServers() ([]*Server, error) {
	return u.store.GetServers(u)
}

func (u *User) AddServer(server *Server) error {
	return u.store.AddServer(u, server)
}

func (u *User) RemoveServer(address string) error {
	return u.store.RemoveServer(u, address)
}

func (u *User) SetNick(nick, address string) error {
	return u.store.SetNick(u, nick, address)
}

func (u *User) SetServerName(name, address string) error {
	return u.store.SetServerName(u, name, address)
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

type Message struct {
	ID      string `json:"-" bleve:"-"`
	Server  string `json:"-" bleve:"server"`
	From    string `bleve:"-"`
	To      string `json:"-" bleve:"to"`
	Content string `bleve:"content"`
	Time    int64  `bleve:"-"`
}

func (m Message) Type() string {
	return "message"
}

func (u *User) LogMessage(id, server, from, to, content string) error {
	message := &Message{
		ID:      id,
		Server:  server,
		From:    from,
		To:      to,
		Content: content,
		Time:    time.Now().Unix(),
	}

	err := u.messageLog.LogMessage(message)
	if err != nil {
		return err
	}
	return u.messageIndex.Index(id, message)
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
