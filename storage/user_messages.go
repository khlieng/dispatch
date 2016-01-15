package storage

import (
	"encoding/json"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/khlieng/dispatch/Godeps/_workspace/src/github.com/blevesearch/bleve"
	"github.com/khlieng/dispatch/Godeps/_workspace/src/github.com/boltdb/bolt"
)

type Message struct {
	ID      uint64 `json:"id"`
	Server  string `json:"server"`
	From    string `json:"from"`
	To      string `json:"to"`
	Content string `json:"content"`
	Time    int64  `json:"time"`
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
	}

	return nil, nil
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

func (u *User) openMessageLog() error {
	err := os.MkdirAll(Path.User(u.Username), 0700)
	if err != nil {
		return err
	}

	u.messageLog, err = bolt.Open(Path.Log(u.Username), 0600, nil)
	if err != nil {
		return err
	}

	u.messageLog.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists(bucketMessages)

		return nil
	})

	indexPath := Path.Index(u.Username)
	u.messageIndex, err = bleve.Open(indexPath)
	if err == bleve.ErrorIndexPathDoesNotExist {
		mapping := bleve.NewIndexMapping()
		u.messageIndex, err = bleve.New(indexPath, mapping)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	return nil
}
