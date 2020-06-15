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

type Features struct {
	Network  string
	Features map[string]interface{}
}

type NetworkName struct {
	Network string
	Name    string
}

type ReconnectSettings struct {
	Network    string
	SkipVerify bool
}

type ConnectionUpdate struct {
	Network   string
	Connected bool
	Error     string
	ErrorType string
}

func newConnectionUpdate(network string, state irc.ConnectionState) ConnectionUpdate {
	status := ConnectionUpdate{
		Network:   network,
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
	Network string
	Old     string `json:"oldNick,omitempty"`
	New     string `json:"newNick,omitempty"`
}

type NickFail struct {
	Network string
}

type Join struct {
	Network  string
	User     string
	Channels []string
}

type Part struct {
	Network  string
	User     string
	Channel  string
	Channels []string
	Reason   string
}

type Mode struct {
	*irc.Mode
}

type Quit struct {
	Network string
	User    string
	Reason  string
}

type Message struct {
	ID      string
	Network string
	From    string
	To      string
	Content string
	Type    string
}

type Messages struct {
	Network  string
	To       string
	Messages []storage.Message
	Prepend  bool
	Next     string
}

type Topic struct {
	Network string
	Channel string
	Topic   string
	Nick    string
}

type Userlist struct {
	Network string
	Channel string
	Users   []string
}

type MOTD struct {
	Network string
	Title   string
	Content []string
}

type Invite struct {
	Network string
	Channel string
	User    string
}

type Kick struct {
	Network string
	Channel string
	Sender  string
	User    string
	Reason  string
}

type Whois struct {
	Network string
	User    string
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
	Network string
	Message string
}

type Raw struct {
	Network string
	Message string
}

type SearchRequest struct {
	Network string
	Channel string
	Phrase  string
}

type SearchResult struct {
	Network string
	Channel string
	Results []storage.Message
}

type ClientCert struct {
	Cert string
	Key  string
}

type FetchMessages struct {
	Network string
	Channel string
	Next    string
}

type Error struct {
	Network string
	Message string
}

type IRCError struct {
	Network string
	Target  string
	Message string
}

type ChannelSearch struct {
	Network string
	Q       string
	Start   int
}

type ChannelSearchResult struct {
	ChannelSearch
	Results []*storage.ChannelListItem
}

type ChannelForward struct {
	Network string
	Old     string
	New     string
}

type DCCSend struct {
	Network  string
	From     string
	Filename string
	URL      string
}

type Tab struct {
	storage.Tab
}
