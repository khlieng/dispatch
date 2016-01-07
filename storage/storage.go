package storage

import (
	"log"
	"os"

	"github.com/khlieng/dispatch/Godeps/_workspace/src/github.com/boltdb/bolt"
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

	err := os.MkdirAll(Path.Logs(), 0700)
	if err != nil {
		log.Fatal(err)
	}
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
