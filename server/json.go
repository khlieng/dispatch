package server

import (
	"encoding/json"

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

type Connect struct {
	Name     string `json:"name"`
	Server   string `json:"server"`
	TLS      bool   `json:"tls"`
	Password string `json:"password"`
	Nick     string `json:"nick"`
	Username string `json:"username"`
	Realname string `json:"realname"`
}

type Nick struct {
	Server   string   `json:"server"`
	Old      string   `json:"old"`
	New      string   `json:"new"`
	Channels []string `json:"channels"`
}

type Join struct {
	Server   string   `json:"server"`
	User     string   `json:"user"`
	Channels []string `json:"channels"`
}

type Part struct {
	Join
	Reason string `json:"reason,omitempty"`
}

type Mode struct {
	Server  string `json:"server"`
	Channel string `json:"channel"`
	User    string `json:"user"`
	Add     string `json:"add"`
	Remove  string `json:"remove"`
}

type Quit struct {
	Server   string   `json:"server"`
	User     string   `json:"user"`
	Reason   string   `json:"reason,omitempty"`
	Channels []string `json:"channels"`
}

type Chat struct {
	Server  string `json:"server"`
	From    string `json:"from"`
	To      string `json:"to,omitempty"`
	Message string `json:"message"`
}

type Topic struct {
	Server  string `json:"server"`
	Channel string `json:"channel"`
	Topic   string `json:"topic"`
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

type Error struct {
	Server  string `json:"server"`
	Message string `json:"message"`
}
