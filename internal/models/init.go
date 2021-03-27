package models

import (
	"github.com/CheerChen/konachan-app/internal/log"
	bolt "go.etcd.io/bbolt"
)

var (
	db *bolt.DB
)

func OpenDbfile(f string) {
	var err error
	db, err = bolt.Open(f, 0600, nil)
	if err != nil {
		log.Fatalf("open DB: %s", err)
	}

	err = db.Update(func(tx *bolt.Tx) (err error) {
		_, err = tx.CreateBucketIfNotExists([]byte("post"))
		if err != nil {
			return
		}
		return
	})

	if err != nil {
		log.Fatalf("create bucket: %s", err)
	}
}
