package models

import (
	"errors"
	"net/http"
	"time"

	"github.com/boltdb/bolt"

	"github.com/CheerChen/konachan-app/internal/cache"
	"github.com/CheerChen/konachan-app/internal/log"
)

var db *bolt.DB
var mem *cache.Cache
var ErrRecordNotFound = errors.New("record not found")
var proxyClient *http.Client

func init() {
	db = getDb()
	mem = cache.New(1*time.Hour, 2*time.Hour)
}

func SetClient(c *http.Client) {
	proxyClient = c
}

func getDb() *bolt.DB {
	db, err := bolt.Open("my.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
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

	return db
}
