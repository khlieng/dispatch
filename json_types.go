package main

import (
	"encoding/json"
)

type WSRequest struct {
	Type    string          `json:"type"`
	Request json.RawMessage `json:"request"`
}

type WSResponse struct {
	Type     string           `json:"type"`
	Response *json.RawMessage `json:"response"`
}

type Connect struct {
	Server   string `json:"server"`
	TLS      bool   `json:"tls"`
	Name     string `json:"name"`
	Nick     string `json:"nick"`
	Username string `json:"username"`
}

type Nick struct {
	Server string `json:"server"`
	Old    string `json:"old"`
	New    string `json:"new"`
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
	Server string `json:"server"`
	User   string `json:"user"`
	Reason string `json:"reason,omitempty"`
}

type Chat struct {
	Server  string `json:"server"`
	From    string `json:"from"`
	To      string `json:"to"`
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

type Error struct {
	Server  string `json:"server"`
	Message string `json:"message"`
}
