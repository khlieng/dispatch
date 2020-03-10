package storage

import (
	"errors"
	"os"

	"github.com/khlieng/dispatch/pkg/session"
	"github.com/spf13/viper"
)

var DataPath directory
var ConfigPath directory

func Initialize() {
	if viper.GetString("dir") != DefaultDirectory() {
		DataPath = directory(viper.GetString("dir"))
		ConfigPath = directory(viper.GetString("dir"))
	} else {
		DataPath = directory(viper.GetString("data"))
		ConfigPath = directory(viper.GetString("conf"))
	}
	os.MkdirAll(DataPath.Root(), 0700)
	os.MkdirAll(ConfigPath.Root(), 0700)
}

var (
	ErrNotFound = errors.New("no item found")
)

type Store interface {
	GetUsers() ([]*User, error)
	SaveUser(user *User) error
	DeleteUser(user *User) error

	GetServer(user *User, host string) (*Server, error)
	GetServers(user *User) ([]*Server, error)
	SaveServer(user *User, server *Server) error
	RemoveServer(user *User, host string) error

	GetChannels(user *User) ([]*Channel, error)
	AddChannel(user *User, channel *Channel) error
	RemoveChannel(user *User, server, channel string) error
}

type SessionStore interface {
	GetSessions() ([]*session.Session, error)
	SaveSession(session *session.Session) error
	DeleteSession(key string) error
}

type MessageStore interface {
	LogMessage(message *Message) error
	GetMessages(server, channel string, count int, fromID string) ([]Message, bool, error)
	GetMessagesByID(server, channel string, ids []string) ([]Message, error)
	Close()
}

type MessageSearchProvider interface {
	SearchMessages(server, channel, q string) ([]string, error)
	Index(id string, message *Message) error
	Close()
}
