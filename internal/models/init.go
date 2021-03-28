package models

import (
	bolt "go.etcd.io/bbolt"

	"github.com/CheerChen/konachan-app/internal/logger"
)

var (
	db  *bolt.DB
	log logger.Logger
)

func Log(l logger.Logger) {
	log = l
}

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
