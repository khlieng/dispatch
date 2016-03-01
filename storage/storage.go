package storage

import (
	"encoding/binary"
	"log"

	"github.com/boltdb/bolt"
)

var (
	Path directory

	db *bolt.DB

	bucketUsers    = []byte("Users")
	bucketServers  = []byte("Servers")
	bucketChannels = []byte("Channels")
	bucketMessages = []byte("Messages")
)

func Initialize(dir string) {
	Path = directory(dir)
}

func Open() {
	var err error
	db, err = bolt.Open(Path.Database(), 0600, nil)
	if err != nil {
		log.Fatal("Could not open database:", err)
	}

	db.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists(bucketUsers)
		tx.CreateBucketIfNotExists(bucketServers)
		tx.CreateBucketIfNotExists(bucketChannels)

		return nil
	})
}

func Close() {
	db.Close()
}

func idToBytes(i uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, i)
	return b
}

func idFromBytes(b []byte) uint64 {
	return binary.BigEndian.Uint64(b)
}
