package boltdb

import (
	"bytes"
	"encoding/binary"
	"strconv"

	"github.com/boltdb/bolt"

	"github.com/khlieng/dispatch/pkg/session"
	"github.com/khlieng/dispatch/storage"
)

var (
	bucketUsers    = []byte("Users")
	bucketServers  = []byte("Servers")
	bucketChannels = []byte("Channels")
	bucketMessages = []byte("Messages")
	bucketSessions = []byte("Sessions")
)

// BoltStore implements storage.Store, storage.MessageStore and storage.SessionStore
type BoltStore struct {
	db *bolt.DB
}

func New(path string) (*BoltStore, error) {
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return nil, err
	}

	db.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists(bucketUsers)
		tx.CreateBucketIfNotExists(bucketServers)
		tx.CreateBucketIfNotExists(bucketChannels)
		tx.CreateBucketIfNotExists(bucketMessages)
		tx.CreateBucketIfNotExists(bucketSessions)
		return nil
	})

	return &BoltStore{
		db,
	}, nil
}

func (s *BoltStore) Close() {
	s.db.Close()
}

func (s *BoltStore) GetUsers() ([]storage.User, error) {
	var users []storage.User

	s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketUsers)

		return b.ForEach(func(k, _ []byte) error {
			id := idFromBytes(k)
			user := storage.User{
				ID:       id,
				IDBytes:  make([]byte, 8),
				Username: strconv.FormatUint(id, 10),
			}
			copy(user.IDBytes, k)

			users = append(users, user)

			return nil
		})
	})

	return users, nil
}

func (s *BoltStore) SaveUser(user *storage.User) error {
	return s.db.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketUsers)

		user.ID, _ = b.NextSequence()
		user.Username = strconv.FormatUint(user.ID, 10)

		data, err := user.Marshal(nil)
		if err != nil {
			return err
		}

		user.IDBytes = idToBytes(user.ID)
		return b.Put(user.IDBytes, data)
	})
}

func (s *BoltStore) DeleteUser(user *storage.User) error {
	return s.db.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketServers)
		c := b.Cursor()

		for k, _ := c.Seek(user.IDBytes); bytes.HasPrefix(k, user.IDBytes); k, _ = c.Next() {
			b.Delete(k)
		}

		b = tx.Bucket(bucketChannels)
		c = b.Cursor()

		for k, _ := c.Seek(user.IDBytes); bytes.HasPrefix(k, user.IDBytes); k, _ = c.Next() {
			b.Delete(k)
		}

		return tx.Bucket(bucketUsers).Delete(user.IDBytes)
	})
}

func (s *BoltStore) GetServers(user *storage.User) ([]storage.Server, error) {
	var servers []storage.Server

	s.db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket(bucketServers).Cursor()

		for k, v := c.Seek(user.IDBytes); bytes.HasPrefix(k, user.IDBytes); k, v = c.Next() {
			server := storage.Server{}
			server.Unmarshal(v)
			servers = append(servers, server)
		}

		return nil
	})

	return servers, nil
}

func (s *BoltStore) AddServer(user *storage.User, server *storage.Server) error {
	return s.db.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketServers)
		data, _ := server.Marshal(nil)

		return b.Put(serverID(user, server.Host), data)
	})
}

func (s *BoltStore) RemoveServer(user *storage.User, address string) error {
	return s.db.Batch(func(tx *bolt.Tx) error {
		serverID := serverID(user, address)
		tx.Bucket(bucketServers).Delete(serverID)

		b := tx.Bucket(bucketChannels)
		c := b.Cursor()

		for k, _ := c.Seek(serverID); bytes.HasPrefix(k, serverID); k, _ = c.Next() {
			b.Delete(k)
		}

		return nil
	})
}

func (s *BoltStore) SetNick(user *storage.User, nick, address string) error {
	return s.db.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketServers)
		id := serverID(user, address)

		server := storage.Server{}
		v := b.Get(id)
		if v != nil {
			server.Unmarshal(v)
			server.Nick = nick

			data, _ := server.Marshal(nil)
			return b.Put(id, data)
		}

		return nil
	})
}

func (s *BoltStore) SetServerName(user *storage.User, name, address string) error {
	return s.db.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketServers)
		id := serverID(user, address)

		server := storage.Server{}
		v := b.Get(id)
		if v != nil {
			server.Unmarshal(v)
			server.Name = name

			data, _ := server.Marshal(nil)
			return b.Put(id, data)
		}

		return nil
	})
}

func (s *BoltStore) GetChannels(user *storage.User) ([]storage.Channel, error) {
	var channels []storage.Channel

	s.db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket(bucketChannels).Cursor()

		for k, v := c.Seek(user.IDBytes); bytes.HasPrefix(k, user.IDBytes); k, v = c.Next() {
			channel := storage.Channel{}
			channel.Unmarshal(v)
			channels = append(channels, channel)
		}

		return nil
	})

	return channels, nil
}

func (s *BoltStore) AddChannel(user *storage.User, channel *storage.Channel) error {
	return s.db.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketChannels)
		data, _ := channel.Marshal(nil)

		return b.Put(channelID(user, channel.Server, channel.Name), data)
	})
}

func (s *BoltStore) RemoveChannel(user *storage.User, server, channel string) error {
	return s.db.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketChannels)
		id := channelID(user, server, channel)

		return b.Delete(id)
	})
}

func (s *BoltStore) LogMessage(message *storage.Message) error {
	return s.db.Batch(func(tx *bolt.Tx) error {
		b, err := tx.Bucket(bucketMessages).CreateBucketIfNotExists([]byte(message.Server + ":" + message.To))
		if err != nil {
			return err
		}

		data, err := message.Marshal(nil)
		if err != nil {
			return err
		}

		return b.Put([]byte(message.ID), data)
	})
}

func (s *BoltStore) GetMessages(server, channel string, count int, fromID string) ([]storage.Message, bool, error) {
	messages := make([]storage.Message, count)
	hasMore := false

	s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketMessages).Bucket([]byte(server + ":" + channel))
		if b == nil {
			return nil
		}

		c := b.Cursor()

		if fromID != "" {
			c.Seek([]byte(fromID))

			for k, v := c.Prev(); count > 0 && k != nil; k, v = c.Prev() {
				count--
				messages[count].Unmarshal(v)
			}
		} else {
			for k, v := c.Last(); count > 0 && k != nil; k, v = c.Prev() {
				count--
				messages[count].Unmarshal(v)
			}
		}

		c.Next()
		k, _ := c.Prev()
		hasMore = k != nil

		return nil
	})

	if count == 0 {
		return messages, hasMore, nil
	} else if count < len(messages) {
		return messages[count:], hasMore, nil
	}

	return nil, false, nil
}

func (s *BoltStore) GetMessagesByID(server, channel string, ids []string) ([]storage.Message, error) {
	messages := make([]storage.Message, len(ids))

	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketMessages).Bucket([]byte(server + ":" + channel))

		for i, id := range ids {
			messages[i].Unmarshal(b.Get([]byte(id)))
		}
		return nil
	})
	return messages, err
}

func (s *BoltStore) GetSessions() ([]session.Session, error) {
	var sessions []session.Session

	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketSessions)

		return b.ForEach(func(_ []byte, v []byte) error {
			session := session.Session{}
			_, err := session.Unmarshal(v)
			sessions = append(sessions, session)
			return err
		})
	})

	if err != nil {
		return nil, err
	}

	return sessions, nil
}

func (s *BoltStore) SaveSession(session *session.Session) error {
	return s.db.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketSessions)

		data, err := session.Marshal(nil)
		if err != nil {
			return err
		}

		return b.Put([]byte(session.Key()), data)
	})
}

func (s *BoltStore) DeleteSession(key string) error {
	return s.db.Batch(func(tx *bolt.Tx) error {
		return tx.Bucket(bucketSessions).Delete([]byte(key))
	})
}

func serverID(user *storage.User, address string) []byte {
	id := make([]byte, 8+len(address))
	copy(id, user.IDBytes)
	copy(id[8:], address)
	return id
}

func channelID(user *storage.User, server, channel string) []byte {
	id := make([]byte, 8+len(server)+1+len(channel))
	copy(id, user.IDBytes)
	copy(id[8:], server)
	copy(id[8+len(server)+1:], channel)
	return id
}

func idToBytes(i uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, i)
	return b
}

func idFromBytes(b []byte) uint64 {
	return binary.BigEndian.Uint64(b)
}
