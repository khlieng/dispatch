package storage

import (
	"bytes"
	"encoding/json"

	"github.com/boltdb/bolt"
)

type Server struct {
	Address  string `json:"address"`
	Nick     string `json:"nick"`
	Username string `json:"username"`
	Realname string `json:"realname"`
}

type Channel struct {
	Server string   `json:"server"`
	Name   string   `json:"name"`
	Users  []string `json:"users"`
	Topic  string   `json:"topic,omitempty"`
}

type User struct {
	UUID string
}

func NewUser(uuid string) User {
	user := User{
		UUID: uuid,
	}

	go db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Users"))
		data, _ := json.Marshal(user)

		b.Put([]byte(uuid), data)

		return nil
	})

	return user
}

func LoadUsers() []User {
	var users []User

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Users"))

		b.ForEach(func(k, v []byte) error {
			users = append(users, User{string(k)})

			return nil
		})

		return nil
	})

	return users
}

func (u User) GetServers() []Server {
	var servers []Server

	db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte("Servers")).Cursor()
		prefix := []byte(u.UUID)

		for k, v := c.Seek(prefix); bytes.HasPrefix(k, prefix); k, v = c.Next() {
			var server Server
			json.Unmarshal(v, &server)
			servers = append(servers, server)
		}

		return nil
	})

	return servers
}

func (u User) GetChannels() []Channel {
	var channels []Channel

	db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte("Channels")).Cursor()

		prefix := []byte(u.UUID)

		for k, v := c.Seek(prefix); bytes.HasPrefix(k, prefix); k, v = c.Next() {
			var channel Channel
			json.Unmarshal(v, &channel)
			channels = append(channels, channel)
		}

		return nil
	})

	return channels
}

func (u User) AddServer(server Server) {
	go db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Servers"))
		data, _ := json.Marshal(server)

		b.Put([]byte(u.UUID+":"+server.Address), data)

		return nil
	})
}

func (u User) AddChannel(channel Channel) {
	go db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Channels"))
		data, _ := json.Marshal(channel)

		b.Put([]byte(u.UUID+":"+channel.Server+":"+channel.Name), data)

		return nil
	})
}

func (u User) RemoveServer(address string) {
	go db.Update(func(tx *bolt.Tx) error {
		serverID := []byte(u.UUID + ":" + address)

		tx.Bucket([]byte("Servers")).Delete(serverID)

		b := tx.Bucket([]byte("Channels"))
		c := b.Cursor()

		for k, _ := c.Seek(serverID); bytes.HasPrefix(k, serverID); k, _ = c.Next() {
			b.Delete(k)
		}

		return nil
	})
}

func (u User) RemoveChannel(server, channel string) {
	go db.Update(func(tx *bolt.Tx) error {
		tx.Bucket([]byte("Channels")).Delete([]byte(u.UUID + ":" + server + ":" + channel))

		return nil
	})
}
