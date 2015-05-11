package storage

import (
	"log"
	"os"
	"os/user"
	"path"

	"github.com/khlieng/name_pending/Godeps/_workspace/src/github.com/boltdb/bolt"
)

var (
	appDir string

	db *bolt.DB

	bucketUsers    = []byte("Users")
	bucketServers  = []byte("Servers")
	bucketChannels = []byte("Channels")
	bucketMessages = []byte("Messages")
)

func init() {
	currentUser, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	appDir = path.Join(currentUser.HomeDir, ".name_pending")
}

func Initialize() {
	var err error

	log.Println("Storing data at", appDir)

	os.Mkdir(appDir, 0777)
	os.Mkdir(path.Join(appDir, "logs"), 0777)

	db, err = bolt.Open(path.Join(appDir, "data.db"), 0600, nil)
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

func Clear() {
	os.RemoveAll(appDir)
}
