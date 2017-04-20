package storage

import (
	"os"
	"time"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/analysis/analyzer/keyword"
	"github.com/boltdb/bolt"
)

type Message struct {
	ID      string `json:"id" bleve:"-"`
	Server  string `json:"-" bleve:"server"`
	From    string `json:"from" bleve:"-"`
	To      string `json:"-" bleve:"to"`
	Content string `json:"content" bleve:"content"`
	Time    int64  `json:"time" bleve:"-"`
}

func (m Message) Type() string {
	return "message"
}

func (u *User) LogMessage(id, server, from, to, content string) error {
	message := Message{
		ID:      id,
		Server:  server,
		From:    from,
		To:      to,
		Content: content,
		Time:    time.Now().Unix(),
	}

	err := u.messageLog.Batch(func(tx *bolt.Tx) error {
		b, err := tx.Bucket(bucketMessages).CreateBucketIfNotExists([]byte(server + ":" + to))
		if err != nil {
			return err
		}

		data, err := message.Marshal(nil)
		if err != nil {
			return err
		}

		return b.Put([]byte(id), data)
	})

	if err != nil {
		return err
	}

	return u.messageIndex.Index(id, message)
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

func (u *User) GetMessages(server, channel string, count int, fromID string) ([]Message, error) {
	messages := make([]Message, count)

	u.messageLog.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketMessages).Bucket([]byte(server + ":" + channel))
		if b == nil {
			return nil
		}

		c := b.Cursor()
		c.Seek([]byte(fromID))

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

func (u *User) SearchMessages(server, channel, q string) ([]Message, error) {
	serverQuery := bleve.NewMatchQuery(server)
	serverQuery.SetField("server")
	channelQuery := bleve.NewMatchQuery(channel)
	channelQuery.SetField("to")
	contentQuery := bleve.NewMatchQuery(q)
	contentQuery.SetField("content")
	contentQuery.SetFuzziness(2)

	query := bleve.NewBooleanQuery()
	query.AddMust(serverQuery, channelQuery, contentQuery)

	search := bleve.NewSearchRequest(query)
	searchResults, err := u.messageIndex.Search(search)
	if err != nil {
		return nil, err
	}

	messages := []Message{}
	u.messageLog.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketMessages).Bucket([]byte(server + ":" + channel))

		for _, hit := range searchResults.Hits {
			message := Message{}
			message.Unmarshal(b.Get([]byte(hit.ID)))
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
		keywordMapping := bleve.NewTextFieldMapping()
		keywordMapping.Analyzer = keyword.Name
		keywordMapping.Store = false
		keywordMapping.IncludeTermVectors = false
		keywordMapping.IncludeInAll = false

		contentMapping := bleve.NewTextFieldMapping()
		contentMapping.Analyzer = "en"
		contentMapping.Store = false
		contentMapping.IncludeTermVectors = false
		contentMapping.IncludeInAll = false

		messageMapping := bleve.NewDocumentMapping()
		messageMapping.StructTagKey = "bleve"
		messageMapping.AddFieldMappingsAt("server", keywordMapping)
		messageMapping.AddFieldMappingsAt("to", keywordMapping)
		messageMapping.AddFieldMappingsAt("content", contentMapping)

		mapping := bleve.NewIndexMapping()
		mapping.AddDocumentMapping("message", messageMapping)

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
