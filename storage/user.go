package storage

import (
	"crypto/tls"
	"os"
	"sync"
	"time"

	"github.com/khlieng/dispatch/pkg/irc"
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

	err = os.MkdirAll(Path.User(user.Username), 0700)
	if err != nil {
		return nil, err
	}
	err = os.Mkdir(Path.Downloads(user.Username), 0700)
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

	return user, nil
}

func LoadUsers(store Store) ([]*User, error) {
	users, err := store.Users()
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

		channels, err := user.Channels()
		if err != nil {
			return nil, err
		}

		for _, channel := range channels {
			messages, _, err := user.LastMessages(channel.Network, channel.Name, 1)
			if err == nil && len(messages) == 1 {
				if _, ok := user.lastMessages[channel.Network]; !ok {
					user.lastMessages[channel.Network] = map[string]*Message{}
				}

				user.lastMessages[channel.Network][channel.Name] = &messages[0]
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

func (u *User) ClientSettings() *ClientSettings {
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

func (u *User) NewNetwork(template *Network, client *irc.Client) *Network {
	if template == nil {
		template = &Network{}
	}

	template.user = u
	template.client = client
	template.channels = map[string]*Channel{}
	template.lock = &sync.Mutex{}

	return template
}

func (u *User) Network(address string) (*Network, error) {
	return u.store.Network(u, address)
}

func (u *User) Networks() ([]*Network, error) {
	return u.store.Networks(u)
}

func (u *User) SaveNetwork(network *Network) error {
	return u.store.SaveNetwork(u, network)
}

func (u *User) RemoveNetwork(address string) error {
	return u.store.RemoveNetwork(u, address)
}

func (u *User) SetNick(nick, address string) error {
	network, err := u.Network(address)
	if err != nil {
		return err
	}
	network.Nick = nick
	return u.SaveNetwork(network)
}

func (u *User) SetNetworkName(name, address string) error {
	network, err := u.Network(address)
	if err != nil {
		return err
	}
	network.Name = name
	return u.SaveNetwork(network)
}

func (u *User) Channels() ([]*Channel, error) {
	return u.store.Channels(u)
}

func (u *User) SaveChannel(channel *Channel) error {
	return u.store.SaveChannel(u, channel)
}

func (u *User) RemoveChannel(network, channel string) error {
	return u.store.RemoveChannel(u, network, channel)
}

func (u *User) HasChannel(network, channel string) bool {
	return u.store.HasChannel(u, network, channel)
}

type Tab struct {
	Network string
	Name    string
}

func (u *User) OpenDMs() ([]Tab, error) {
	return u.store.OpenDMs(u)
}

func (u *User) AddOpenDM(network, nick string) error {
	return u.store.AddOpenDM(u, network, nick)
}

func (u *User) RemoveOpenDM(network, nick string) error {
	return u.store.RemoveOpenDM(u, network, nick)
}

type Message struct {
	ID      string  `json:"-" bleve:"-"`
	Network string  `json:"-" bleve:"server"`
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

	u.setLastMessage(msg.Network, msg.To, msg)

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

func (u *User) LogEvent(network, name string, params []string, channels ...string) error {
	now := time.Now().Unix()
	event := Event{
		Type:   name,
		Params: params,
		Time:   now,
	}

	for _, channel := range channels {
		lastMessage := u.getLastMessage(network, channel)

		if lastMessage != nil && shouldCollapse(lastMessage, event) {
			lastMessage.Events = append(lastMessage.Events, event)
			u.setLastMessage(network, channel, lastMessage)

			err := u.messageLog.LogMessage(lastMessage)
			if err != nil {
				return err
			}
		} else {
			msg := &Message{
				ID:      betterguid.New(),
				Network: network,
				To:      channel,
				Time:    now,
				Events:  []Event{event},
			}
			u.setLastMessage(network, channel, msg)

			err := u.messageLog.LogMessage(msg)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

var collapsed = []string{"join", "part", "quit", "nick"}

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

func (u *User) getLastMessage(network, channel string) *Message {
	u.lock.Lock()
	defer u.lock.Unlock()

	if _, ok := u.lastMessages[network]; !ok {
		return nil
	}

	last := u.lastMessages[network][channel]
	if last != nil {
		msg := *last
		return &msg
	}
	return nil
}

func (u *User) setLastMessage(network, channel string, msg *Message) {
	u.lock.Lock()

	if _, ok := u.lastMessages[network]; !ok {
		u.lastMessages[network] = map[string]*Message{}
	}

	u.lastMessages[network][channel] = msg
	u.lock.Unlock()
}

func (u *User) Messages(network, channel string, count int, fromID string) ([]Message, bool, error) {
	return u.messageLog.Messages(network, channel, count, fromID)
}

func (u *User) LastMessages(network, channel string, count int) ([]Message, bool, error) {
	return u.Messages(network, channel, count, "")
}

func (u *User) SearchMessages(network, channel, q string) ([]Message, error) {
	ids, err := u.messageIndex.SearchMessages(network, channel, q)
	if err != nil {
		return nil, err
	}

	return u.messageLog.MessagesByID(network, channel, ids)
}
