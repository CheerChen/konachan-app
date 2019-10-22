package models

import (
	"errors"
	"time"

	"github.com/boltdb/bolt"

	"github.com/CheerChen/konachan-app/internal/cache"
	"github.com/CheerChen/konachan-app/internal/log"
)

var db *bolt.DB
var cc *cache.Handler
var ErrRecordNotFound = errors.New("record not found")

func Init() {
	db = GetDB()
	cc = &cache.Handler{"./cache", 60}
}

func GetDB() *bolt.DB {
	db, err := bolt.Open("my.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatalf("open DB: %s", err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("post"))
		return err
	})

	if err != nil {
		log.Fatalf("create bucket: %s", err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("post_tag"))
		return err
	})

	if err != nil {
		log.Fatalf("create bucket: %s", err)
	}

	return db
}
