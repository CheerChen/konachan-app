package models

import (
	"encoding/json"
	"errors"
	"sort"
	"strconv"
	"strings"

	"github.com/CheerChen/konachan-app/internal/log"
	bolt "go.etcd.io/bbolt"
)

type Post struct {
	OriginalPost

	TfIDf   float64 `json:"tf_idf"`
	MyScore float64 `json:"my_score"`
	IsFav   bool    `json:"is_fav"`
}

type Posts []Post

func (p *Post) TableName() string {
	return "post"
}

func (ps *Posts) TableName() string {
	return "post"
}

func (p *Post) Save() (err error) {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(p.TableName()))
		key := []byte(strconv.FormatInt(p.ID, 10))
		value, err := json.Marshal(p)
		if err != nil {
			return err
		}
		return b.Put(key, value)
	})
}

func (p *Post) Find(id int64) (err error) {
	return db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(p.TableName()))
		key := []byte(strconv.FormatInt(id, 10))
		value := b.Get(key)
		if value == nil {
			return errors.New("record not found")
		}
		return json.Unmarshal(value, p)
	})
}

func (p *Post) Delete() (err error) {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(p.TableName()))
		key := []byte(strconv.FormatInt(p.ID, 10))
		return b.Delete(key)
	})
}

func (ps *Posts) FetchAllId() (idMap map[int64]bool, err error) {
	idMap = make(map[int64]bool)
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(ps.TableName()))

		_ = b.ForEach(func(k, _ []byte) error {
			id, err := strconv.ParseInt(string(k), 10, 64)
			if err == nil {
				idMap[id] = true
			}
			return nil
		})
		return nil
	})
	return
}

func (ps *Posts) FetchAll(tag string, l, page int) (err error) {
	log.Infof("fetch album tag=%s, limit=%d, page=%d", tag, l, page)

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(ps.TableName()))

		_ = b.ForEach(func(_, v []byte) error {
			var p Post
			err := json.Unmarshal(v, &p)
			if err != nil {
				return err
			}
			if tag != "" && !strings.Contains(p.Tags, tag) {
				return nil
			}

			*ps = append(*ps, p)
			return nil
		})
		return nil
	})

	sort.Slice(*ps, func(i, j int) bool {
		return (*ps)[i].ID > (*ps)[j].ID
	})

	min, max, start, end := 0, len(*ps), (page-1)*100, (page-1)*100+l

	if start < min {
		start = min
	}
	if start > max {
		start = max
	}
	if end > max {
		end = max
	}
	if end < min {
		end = min
	}

	*ps = (*ps)[start:end]
	return nil
}
