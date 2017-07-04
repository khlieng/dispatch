package server

import (
	"crypto/x509"
	"encoding/json"

	"github.com/khlieng/dispatch/irc"
	"github.com/khlieng/dispatch/storage"
)

type WSRequest struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

type WSResponse struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type Server struct {
	storage.Server
	Status ConnectionUpdate `json:"status"`
}

type ServerName struct {
	Server string `json:"server"`
	Name   string `json:"name"`
}

type ReconnectSettings struct {
	Server     string `json:"server"`
	SkipVerify bool   `json:"skipVerify"`
}

type ConnectionUpdate struct {
	Server    string `json:"server"`
	Connected bool   `json:"connected"`
	Error     string `json:"error,omitempty"`
	ErrorType string `json:"errorType,omitempty"`
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
	Server string `json:"server"`
	Old    string `json:"oldNick"`
	New    string `json:"newNick"`
}

type NickFail struct {
	Server string `json:"server"`
}

type Join struct {
	Server   string   `json:"server"`
	User     string   `json:"user"`
	Channels []string `json:"channels"`
}

type Part struct {
	Server   string   `json:"server"`
	User     string   `json:"user"`
	Channel  string   `json:"channel,omitempty"`
	Channels []string `json:"channels,omitempty"`
	Reason   string   `json:"reason,omitempty"`
}

type Mode struct {
	Server  string `json:"server"`
	Channel string `json:"channel"`
	User    string `json:"user"`
	Add     string `json:"add"`
	Remove  string `json:"remove"`
}

type Quit struct {
	Server string `json:"server"`
	User   string `json:"user"`
	Reason string `json:"reason,omitempty"`
}

type Message struct {
	ID      string `json:"id,omitempty"`
	Server  string `json:"server,omitempty"`
	From    string `json:"from,omitempty"`
	To      string `json:"to,omitempty"`
	Content string `json:"content"`
	Type    string `json:"type,omitempty"`
}

type Messages struct {
	Server   string            `json:"server"`
	To       string            `json:"to"`
	Messages []storage.Message `json:"messages"`
	Prepend  bool              `json:"prepend,omitempty"`
	Next     string            `json:"next,omitempty"`
}

type Topic struct {
	Server  string `json:"server"`
	Channel string `json:"channel"`
	Topic   string `json:"topic,omitempty"`
	Nick    string `json:"nick,omitempty"`
}

type Userlist struct {
	Server  string   `json:"server"`
	Channel string   `json:"channel"`
	Users   []string `json:"users"`
}

type MOTD struct {
	Server  string   `json:"server"`
	Title   string   `json:"title"`
	Content []string `json:"content"`
}

type Invite struct {
	Server  string `json:"server"`
	Channel string `json:"channel"`
	User    string `json:"user"`
}

type Kick struct {
	Server  string `json:"server"`
	Channel string `json:"channel"`
	User    string `json:"user"`
}

type Whois struct {
	Server string `json:"server"`
	User   string `json:"user"`
}

type WhoisReply struct {
	Nick     string   `json:"nick"`
	Username string   `json:"username"`
	Host     string   `json:"host"`
	Realname string   `json:"realname"`
	Server   string   `json:"server"`
	Channels []string `json:"channels"`
}

type Away struct {
	Server  string `json:"server"`
	Message string `json:"message"`
}

type Raw struct {
	Server  string `json:"server"`
	Message string `json:"message"`
}

type SearchRequest struct {
	Server  string `json:"server"`
	Channel string `json:"channel"`
	Phrase  string `json:"phrase"`
}

type SearchResult struct {
	Server  string            `json:"server"`
	Channel string            `json:"channel"`
	Results []storage.Message `json:"results"`
}

type ClientCert struct {
	Cert []byte `json:"cert"`
	Key  []byte `json:"key"`
}

type FetchMessages struct {
	Server  string `json:"server"`
	Channel string `json:"channel"`
	Next    string `json:"next"`
}

type Error struct {
	Server  string `json:"server"`
	Message string `json:"message"`
}
