package storage

import (
	"bytes"
	"os"
	"strconv"

	"github.com/khlieng/dispatch/Godeps/_workspace/src/github.com/boltdb/bolt"
)

func NewUser() (*User, error) {
	user := &User{}

	err := db.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketUsers)

		user.ID, _ = b.NextSequence()
		user.Username = strconv.FormatUint(user.ID, 10)

		data, err := user.MarshalMsg(nil)
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
			server := Server{}
			server.UnmarshalMsg(v)
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
			channel := Channel{}
			channel.UnmarshalMsg(v)
			channels = append(channels, channel)
		}

		return nil
	})

	return channels
}

func (u *User) AddServer(server Server) {
	db.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketServers)
		data, _ := server.MarshalMsg(nil)

		b.Put(u.serverID(server.Host), data)

		return nil
	})
}

func (u *User) AddChannel(channel Channel) {
	db.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketChannels)
		data, _ := channel.MarshalMsg(nil)

		b.Put(u.channelID(channel.Server, channel.Name), data)

		return nil
	})
}

func (u *User) SetNick(nick, address string) {
	db.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketServers)
		id := u.serverID(address)

		server := Server{}
		server.UnmarshalMsg(b.Get(id))
		server.Nick = nick

		data, _ := server.MarshalMsg(nil)
		b.Put(id, data)

		return nil
	})
}

func (u *User) RemoveServer(address string) {
	db.Batch(func(tx *bolt.Tx) error {
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
	db.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketChannels)
		id := u.channelID(server, channel)

		b.Delete(id)

		return nil
	})
}

func (u *User) Remove() {
	db.Batch(func(tx *bolt.Tx) error {
		return tx.Bucket(bucketUsers).Delete(u.id)
	})
	u.closeMessageLog()
	os.RemoveAll(Path.User(u.Username))
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
