package storage

import (
	"errors"
	"os"

	"github.com/khlieng/dispatch/pkg/session"
)

var (
	Path directory

	GetMessageStore          MessageStoreCreator
	GetMessageSearchProvider MessageSearchProviderCreator
)

func Initialize(root, dataRoot, configRoot string) {
	if root != DefaultDirectory() {
		Path.dataRoot = root
		Path.configRoot = root
	} else {
		Path.dataRoot = dataRoot
		Path.configRoot = configRoot
	}
	os.MkdirAll(Path.DataRoot(), 0700)
	os.MkdirAll(Path.ConfigRoot(), 0700)
}

var (
	ErrNotFound = errors.New("no item found")
)

type Store interface {
	Users() ([]*User, error)
	SaveUser(user *User) error
	DeleteUser(user *User) error

	Network(user *User, host string) (*Network, error)
	Networks(user *User) ([]*Network, error)
	SaveNetwork(user *User, network *Network) error
	RemoveNetwork(user *User, host string) error

	Channels(user *User) ([]*Channel, error)
	HasChannel(user *User, network, channel string) bool
	SaveChannel(user *User, channel *Channel) error
	RemoveChannel(user *User, network, channel string) error

	OpenDMs(user *User) ([]Tab, error)
	AddOpenDM(user *User, network, nick string) error
	RemoveOpenDM(user *User, network, nick string) error
}

type SessionStore interface {
	Sessions() ([]*session.Session, error)
	SaveSession(session *session.Session) error
	DeleteSession(key string) error
}

type MessageStore interface {
	LogMessage(message *Message) error
	LogMessages(messages []*Message) error
	Messages(network, channel string, count int, fromID string) ([]Message, bool, error)
	MessagesByID(network, channel string, ids []string) ([]Message, error)
	Close()
}

type MessageStoreCreator func(*User) (MessageStore, error)

type MessageSearchProvider interface {
	SearchMessages(network, channel, q string) ([]string, error)
	Index(id string, message *Message) error
	Close()
}

type MessageSearchProviderCreator func(*User) (MessageSearchProvider, error)
