package storage

import (
	"log"
	"os"

	"github.com/khlieng/dispatch/Godeps/_workspace/src/github.com/boltdb/bolt"
)

var (
	appDir string
	db     *bolt.DB

	bucketUsers    = []byte("Users")
	bucketServers  = []byte("Servers")
	bucketChannels = []byte("Channels")
	bucketMessages = []byte("Messages")
)

func Initialize() {
	log.Println("Storing data at", Path.Root())

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

func Clear() {
	os.RemoveAll(Path.Logs())
	os.Remove(Path.Database())
}
