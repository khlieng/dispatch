package storage

import (
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
	To      string `json:"to,omitempty"`
	Content string `json:"content"`
	Time    int64  `json:"time"`
}

func (u *User) LogMessage(server, from, to, content string) error {
	message := Message{
		Server:  server,
		From:    from,
		To:      to,
		Content: content,
		Time:    time.Now().Unix(),
	}
	bucketKey := server + ":" + to

	err := u.messageLog.Batch(func(tx *bolt.Tx) error {
		b, err := tx.Bucket(bucketMessages).CreateBucketIfNotExists([]byte(bucketKey))
		if err != nil {
			return err
		}

		message.ID, _ = b.NextSequence()

		data, err := message.Marshal(nil)
		if err != nil {
			return err
		}

		return b.Put(idToBytes(message.ID), data)
	})

	if err != nil {
		return err
	}

	return u.messageIndex.Index(bucketKey+":"+strconv.FormatUint(message.ID, 10), message)
}

func (u *User) GetLastMessages(server, channel string, count int) ([]Message, error) {
	messages := make([]Message, count)

	u.messageLog.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketMessages).Bucket([]byte(server + ":" + channel))
		if b == nil {
			return nil
		}

		c := b.Cursor()

		for _, v := c.Last(); count > 0 && v != nil; _, v = c.Prev() {
			count--
			messages[count].Unmarshal(v)
		}

		return nil
	})

	if count == 0 {
		return messages, nil
	} else if count < len(messages) {
		return messages[count:], nil
	}

	return nil, nil
}

func (u *User) GetMessages(server, channel string, count int, fromID uint64) ([]Message, error) {
	messages := make([]Message, count)

	u.messageLog.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketMessages).Bucket([]byte(server + ":" + channel))
		if b == nil {
			return nil
		}

		c := b.Cursor()
		c.Seek(idToBytes(fromID))

		for k, v := c.Prev(); count > 0 && k != nil; k, v = c.Prev() {
			count--
			messages[count].Unmarshal(v)
		}

		return nil
	})

	if count == 0 {
		return messages, nil
	} else if count < len(messages) {
		return messages[count:], nil
	}

	return nil, nil
}

func (u *User) SearchMessages(server, channel, phrase string) ([]Message, error) {
	serverQuery := bleve.NewMatchQuery(server)
	serverQuery.SetField("server")
	channelQuery := bleve.NewMatchQuery(channel)
	channelQuery.SetField("to")
	contentQuery := bleve.NewFuzzyQuery(phrase)
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
			id, _ := strconv.ParseUint(hit.ID[idx+1:], 10, 64)

			message := Message{}
			message.Unmarshal(bc.Get(idToBytes(id)))
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

func (u *User) closeMessageLog() {
	u.messageLog.Close()
	u.messageIndex.Close()
}
