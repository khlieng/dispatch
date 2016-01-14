package storage

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/khlieng/dispatch/Godeps/_workspace/src/github.com/blevesearch/bleve"
	"github.com/khlieng/dispatch/Godeps/_workspace/src/github.com/boltdb/bolt"
)

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

type Message struct {
	ID      uint64 `json:"id"`
	Server  string `json:"server"`
	From    string `json:"from"`
	To      string `json:"to"`
	Content string `json:"content"`
	Time    int64  `json:"time"`
}

type User struct {
	UUID string

	messageLog   *bolt.DB
	messageIndex bleve.Index
	certificate  *tls.Certificate
	lock         sync.Mutex
}

func NewUser(uuid string) *User {
	user := &User{
		UUID: uuid,
	}

	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Users"))
		data, _ := json.Marshal(user)

		b.Put([]byte(uuid), data)

		return nil
	})

	user.openMessageLog()

	return user
}

func LoadUsers() []*User {
	var users []*User

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Users"))

		b.ForEach(func(k, v []byte) error {
			user := User{UUID: string(k)}
			user.openMessageLog()
			user.loadCertificate()

			users = append(users, &user)

			return nil
		})

		return nil
	})

	return users
}

func (u *User) GetServers() []Server {
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

func (u *User) GetChannels() []Channel {
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

func (u *User) AddServer(server Server) {
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Servers"))
		data, _ := json.Marshal(server)

		b.Put([]byte(u.UUID+":"+server.Host), data)

		return nil
	})
}

func (u *User) AddChannel(channel Channel) {
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Channels"))
		data, _ := json.Marshal(channel)

		b.Put([]byte(u.UUID+":"+channel.Server+":"+channel.Name), data)

		return nil
	})
}

func (u *User) SetNick(nick, address string) {
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Servers"))
		id := []byte(u.UUID + ":" + address)
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

func (u *User) RemoveChannel(server, channel string) {
	db.Update(func(tx *bolt.Tx) error {
		tx.Bucket([]byte("Channels")).Delete([]byte(u.UUID + ":" + server + ":" + channel))

		return nil
	})
}

func (u *User) LogMessage(server, from, to, content string) {
	bucketKey := server + ":" + to
	var id uint64
	var idStr string
	var message Message

	u.messageLog.Update(func(tx *bolt.Tx) error {
		b, _ := tx.Bucket(bucketMessages).CreateBucketIfNotExists([]byte(bucketKey))
		id, _ = b.NextSequence()
		idStr = strconv.FormatUint(id, 10)

		message = Message{
			ID:      id,
			Content: content,
			Server:  server,
			From:    from,
			To:      to,
			Time:    time.Now().Unix(),
		}

		data, _ := json.Marshal(message)
		b.Put([]byte(idStr), data)

		return nil
	})

	u.messageIndex.Index(bucketKey+":"+idStr, message)
}

func (u *User) GetLastMessages(server, channel string, count int) ([]Message, error) {
	messages := make([]Message, count)

	u.messageLog.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketMessages).Bucket([]byte(server + ":" + channel))
		if b == nil {
			return nil
		}

		c := b.Cursor()

		for k, v := c.Last(); count > 0 && k != nil; k, v = c.Prev() {
			count--
			json.Unmarshal(v, &messages[count])
		}

		return nil
	})

	if count < len(messages) {
		return messages[count:], nil
	} else {
		return nil, nil
	}
}

func (u *User) GetMessages(server, channel string, count int, fromID uint64) ([]Message, error) {
	messages := make([]Message, count)

	u.messageLog.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketMessages).Bucket([]byte(server + ":" + channel))
		if b == nil {
			return nil
		}

		c := b.Cursor()
		c.Seek([]byte(strconv.FormatUint(fromID, 10)))

		for k, v := c.Prev(); count > 0 && k != nil; k, v = c.Prev() {
			count--
			json.Unmarshal(v, &messages[count])
		}

		return nil
	})

	if count < len(messages) {
		return messages[count:], nil
	} else {
		return nil, nil
	}
}

func (u *User) SearchMessages(server, channel, phrase string) ([]Message, error) {
	serverQuery := bleve.NewMatchQuery(server)
	serverQuery.SetField("server")
	channelQuery := bleve.NewMatchQuery(channel)
	channelQuery.SetField("to")
	contentQuery := bleve.NewMatchQuery(phrase)
	contentQuery.SetField("content")

	query := bleve.NewBooleanQuery([]bleve.Query{serverQuery, channelQuery, contentQuery}, nil, nil)

	search := bleve.NewSearchRequest(query)
	searchResults, err := u.messageIndex.Search(search)
	if err != nil {
		return nil, err
	}

	messages := []Message{}
	u.messageLog.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketMessages)

		for _, hit := range searchResults.Hits {
			idx := strings.LastIndex(hit.ID, ":")
			bc := b.Bucket([]byte(hit.ID[:idx]))
			var message Message

			json.Unmarshal(bc.Get([]byte(hit.ID[idx+1:])), &message)
			messages = append(messages, message)
		}

		return nil
	})

	return messages, nil
}

func (u *User) Close() {
	u.messageLog.Close()
	u.messageIndex.Close()
}

func (u *User) openMessageLog() {
	var err error

	u.messageLog, err = bolt.Open(Path.Log(u.UUID), 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	u.messageLog.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists(bucketMessages)

		return nil
	})

	indexPath := Path.Index(u.UUID)
	u.messageIndex, err = bleve.Open(indexPath)
	if err == bleve.ErrorIndexPathDoesNotExist {
		mapping := bleve.NewIndexMapping()
		u.messageIndex, err = bleve.New(indexPath, mapping)
		if err != nil {
			log.Fatal(err)
		}
	} else if err != nil {
		log.Fatal(err)
	}
}
