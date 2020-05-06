package server

import (
	"crypto/x509"

	"github.com/mailru/easyjson"

	"github.com/khlieng/dispatch/pkg/irc"
	"github.com/khlieng/dispatch/storage"
)

type WSRequest struct {
	Type string
	Data easyjson.RawMessage
}

type WSResponse struct {
	Type string
	Data interface{}
}

type Server struct {
	*storage.Server
	Status   ConnectionUpdate
	Features map[string]interface{}
}

type Features struct {
	Server   string
	Features map[string]interface{}
}

type ServerName struct {
	Server string
	Name   string
}

type ReconnectSettings struct {
	Server     string
	SkipVerify bool
}

type ConnectionUpdate struct {
	Server    string
	Connected bool
	Error     string
	ErrorType string
}

func newConnectionUpdate(server string, state irc.ConnectionState) ConnectionUpdate {
	status := ConnectionUpdate{
		Server:    server,
		Connected: state.Connected,
	}
	if state.Error != nil {
		status.Error = state.Error.Error()
		if _, ok := state.Error.(x509.UnknownAuthorityError); ok {
			status.ErrorType = "verify"
		}
	}
	return status
}

type Nick struct {
	Server string
	Old    string `json:"oldNick,omitempty"`
	New    string `json:"newNick,omitempty"`
}

type NickFail struct {
	Server string
}

type Join struct {
	Server   string
	User     string
	Channels []string
}

type Part struct {
	Server   string
	User     string
	Channel  string
	Channels []string
	Reason   string
}

type Mode struct {
	Server  string
	Channel string
	User    string
	Add     string
	Remove  string
}

type Quit struct {
	Server string
	User   string
	Reason string
}

type Message struct {
	ID      string
	Server  string
	From    string
	To      string
	Content string
	Type    string
}

type Messages struct {
	Server   string
	To       string
	Messages []storage.Message
	Prepend  bool
	Next     string
}

type Topic struct {
	Server  string
	Channel string
	Topic   string
	Nick    string
}

type Userlist struct {
	Server  string
	Channel string
	Users   []string
}

type MOTD struct {
	Server  string
	Title   string
	Content []string
}

type Invite struct {
	Server  string
	Channel string
	User    string
}

type Kick struct {
	Server  string
	Channel string
	User    string
}

type Whois struct {
	Server string
	User   string
}

type WhoisReply struct {
	Nick     string
	Username string
	Host     string
	Realname string
	Server   string
	Channels []string
}

type Away struct {
	Server  string
	Message string
}

type Raw struct {
	Server  string
	Message string
}

type SearchRequest struct {
	Server  string
	Channel string
	Phrase  string
}

type SearchResult struct {
	Server  string
	Channel string
	Results []storage.Message
}

type ClientCert struct {
	Cert string
	Key  string
}

type FetchMessages struct {
	Server  string
	Channel string
	Next    string
}

type Error struct {
	Server  string
	Message string
}

type IRCError struct {
	Server  string
	Target  string
	Message string
}

type ChannelSearch struct {
	Server string
	Q      string
	Start  int
}

type ChannelSearchResult struct {
	ChannelSearch
	Results []*storage.ChannelListItem
}

type ChannelForward struct {
	Server string
	Old    string
	New    string
}

type Tab struct {
	storage.Tab
}
