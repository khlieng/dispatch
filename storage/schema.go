package storage

import (
	"crypto/tls"
	"sync"

	"github.com/khlieng/dispatch/Godeps/_workspace/src/github.com/blevesearch/bleve"
	"github.com/khlieng/dispatch/Godeps/_workspace/src/github.com/boltdb/bolt"
)

//go:generate msgp

type User struct {
	ID       uint64
	Username string

	id           []byte
	messageLog   *bolt.DB
	messageIndex bleve.Index
	certificate  *tls.Certificate
	lock         sync.Mutex
}

type Server struct {
	Name     string `json:"name"`
	Host     string `json:"host"`
	Port     string `json:"port,omitempty"`
	TLS      bool   `json:"tls"`
	Password string `json:"password,omitempty"`
	Nick     string `json:"nick"`
	Username string `json:"username"`
	Realname string `json:"realname"`
}

type Channel struct {
	Server string `json:"server"`
	Name   string `json:"name"`
	Topic  string `json:"topic,omitempty"`
}

type Message struct {
	ID      uint64 `json:"id"`
	Server  string `json:"server"`
	From    string `json:"from"`
	To      string `json:"to,omitempty"`
	Content string `json:"content"`
	Time    int64  `json:"time"`
}
