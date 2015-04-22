package storage

import (
	"log"
	"os"
	"os/user"
	"path"

	"github.com/boltdb/bolt"
)

var db *bolt.DB

func init() {
	var err error
	currentUser, _ := user.Current()
	appDir := path.Join(currentUser.HomeDir, ".name_pending")

	os.Mkdir(appDir, 0777)

	db, err = bolt.Open(path.Join(appDir, "data.db"), 0600, nil)
	if err != nil {
		log.Fatal("Unable to open database file")
	}

	db.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists([]byte("Users"))
		tx.CreateBucketIfNotExists([]byte("Servers"))
		tx.CreateBucketIfNotExists([]byte("Channels"))
		return nil
	})
}

func Cleanup() {
	db.Close()
}
