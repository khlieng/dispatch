package storage

import (
	"log"
	"os"
	"path"

	"github.com/khlieng/name_pending/Godeps/_workspace/src/github.com/boltdb/bolt"
)

var (
	appDir string
	db     *bolt.DB

	bucketUsers    = []byte("Users")
	bucketServers  = []byte("Servers")
	bucketChannels = []byte("Channels")
	bucketMessages = []byte("Messages")
)

func Initialize(dir string) {
	var err error
	appDir = dir

	log.Println("Storing data at", dir)

	db, err = bolt.Open(path.Join(dir, "data.db"), 0600, nil)
	if err != nil {
		log.Fatal("Could not open database file")
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

func Clear(dir string) {
	os.RemoveAll(path.Join(dir, "logs"))
	os.Remove(path.Join(dir, "data.db"))
}
