package storage

import (
	"log"
	"os"
	"os/user"
	"path"

	"github.com/khlieng/name_pending/Godeps/_workspace/src/github.com/boltdb/bolt"
)

var (
	AppDir string

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

	AppDir = path.Join(currentUser.HomeDir, ".name_pending")

	os.Mkdir(AppDir, 0777)
	os.Mkdir(path.Join(AppDir, "logs"), 0777)
}

func Initialize() {
	var err error

	log.Println("Storing data at", AppDir)

	db, err = bolt.Open(path.Join(AppDir, "data.db"), 0600, nil)
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
	os.RemoveAll(path.Join(AppDir, "logs"))
	os.Remove(path.Join(AppDir, "data.db"))
}
