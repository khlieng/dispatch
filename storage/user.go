package storage

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"strconv"
	"sync"

	"github.com/khlieng/dispatch/Godeps/_workspace/src/github.com/blevesearch/bleve"
	"github.com/khlieng/dispatch/Godeps/_workspace/src/github.com/boltdb/bolt"
)

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
	Server string   `json:"server"`
	Name   string   `json:"name"`
	Users  []string `json:"users,omitempty"`
	Topic  string   `json:"topic,omitempty"`
}

func NewUser() (*User, error) {
	user := &User{}

	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketUsers)

		user.ID, _ = b.NextSequence()
		user.Username = strconv.FormatUint(user.ID, 10)

		data, err := json.Marshal(user)
		if err != nil {
			return err
		}

		user.id = idToBytes(user.ID)
		return b.Put(user.id, data)
	})

	if err != nil {
		return nil, err
	}

	err = user.openMessageLog()
	if err != nil {
		return nil, err
	}

	return user, nil
}

func LoadUsers() []*User {
	var users []*User

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketUsers)

		b.ForEach(func(k, _ []byte) error {
			id := idFromBytes(k)
			user := &User{
				ID:       id,
				Username: strconv.FormatUint(id, 10),
				id:       make([]byte, 8),
			}
			copy(user.id, k)

			users = append(users, user)

			return nil
		})

		return nil
	})

	for _, user := range users {
		user.openMessageLog()
		user.loadCertificate()
	}

	return users
}

func (u *User) GetServers() []Server {
	var servers []Server

	db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket(bucketServers).Cursor()

		for k, v := c.Seek(u.id); bytes.HasPrefix(k, u.id); k, v = c.Next() {
			var server Server
			json.Unmarshal(v, &server)
			servers = append(servers, server)
		}

		return nil
	})

	return servers
}

func (u *User) GetChannels() []Channel {
	var channels []Channel

	db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket(bucketChannels).Cursor()

		for k, v := c.Seek(u.id); bytes.HasPrefix(k, u.id); k, v = c.Next() {
			var channel Channel
			json.Unmarshal(v, &channel)
			channels = append(channels, channel)
		}

		return nil
	})

	return channels
}

func (u *User) AddServer(server Server) {
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketServers)
		data, _ := json.Marshal(server)

		b.Put(u.serverID(server.Host), data)

		return nil
	})
}

func (u *User) AddChannel(channel Channel) {
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketChannels)
		data, _ := json.Marshal(channel)

		b.Put(u.channelID(channel.Server, channel.Name), data)

		return nil
	})
}

func (u *User) SetNick(nick, address string) {
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketServers)
		id := u.serverID(address)
		var server Server

		json.Unmarshal(b.Get(id), &server)
		server.Nick = nick

		data, _ := json.Marshal(server)
		b.Put(id, data)

		return nil
	})
}

func (u *User) RemoveServer(address string) {
	db.Update(func(tx *bolt.Tx) error {
		serverID := u.serverID(address)
		tx.Bucket(bucketServers).Delete(serverID)

		b := tx.Bucket(bucketChannels)
		c := b.Cursor()

		for k, _ := c.Seek(serverID); bytes.HasPrefix(k, serverID); k, _ = c.Next() {
			b.Delete(k)
		}

		return nil
	})
}

func (u *User) RemoveChannel(server, channel string) {
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketChannels)
		id := u.channelID(server, channel)

		b.Delete(id)

		return nil
	})
}

func (u *User) Close() {
	u.messageLog.Close()
	u.messageIndex.Close()
}

func (u *User) serverID(address string) []byte {
	id := make([]byte, 8+len(address))
	copy(id, u.id)
	copy(id[8:], address)
	return id
}

func (u *User) channelID(server, channel string) []byte {
	id := make([]byte, 8+len(server)+1+len(channel))
	copy(id, u.id)
	copy(id[8:], server)
	copy(id[8+len(server)+1:], channel)
	return id
}
