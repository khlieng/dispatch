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
	Nick     string `json:"nick"`
	Username string `json:"username"`
}

type Join struct {
	Server   string   `json:"server"`
	User     string   `json:"user"`
	Channels []string `json:"channels"`
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
	Server  string `json:"server"`
	Title   string `json:"title"`
	Content string `json:"content"`
}
