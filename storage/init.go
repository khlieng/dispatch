package storage

import (
	"github.com/boltdb/bolt"
)

var db *bolt.DB

func init() {
	db, _ = bolt.Open("data.db", 0600, nil)

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
