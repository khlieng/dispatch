package storage

import (
	"github.com/khlieng/dispatch/pkg/session"
)

var Path directory

func Initialize(dir string) {
	Path = directory(dir)
}

type Store interface {
	GetUsers() ([]User, error)
	SaveUser(*User) error
	DeleteUser(*User) error

	GetServers(*User) ([]Server, error)
	AddServer(*User, *Server) error
	RemoveServer(*User, string) error
	SetNick(*User, string, string) error
	SetServerName(*User, string, string) error

	GetChannels(*User) ([]Channel, error)
	AddChannel(*User, *Channel) error
	RemoveChannel(*User, string, string) error
}

type SessionStore interface {
	GetSessions() ([]session.Session, error)
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
