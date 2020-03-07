package models

import (
	"github.com/CheerChen/konachan-app/internal/service/konachan"
)

type Tag struct {
	*konachan.Tag

	TfIdf float64 `json:"tf_idf"`
	Idf   float64 `json:"idf"`
}

type Tags []Tag

//func (t *Tag) TableName() string {
//	return "tag"
//}
//
//func (ts *Tags) TableName() string {
//	return "tag"
//}
//
//func (t *Tag) Save() (err error) {
//	return db.Update(func(tx *bolt.Tx) error {
//		b := tx.Bucket([]byte(t.TableName()))
//		key := []byte(strconv.Itoa(t.ID))
//		value, err := json.Marshal(t)
//		if err != nil {
//			return err
//		}
//		return b.Put(key, value)
//	})
//}
//
//func (ts *Tags) FetchAll() (err error) {
//	err = db.View(func(tx *bolt.Tx) error {
//		b := tx.Bucket([]byte(ts.TableName()))
//
//		_ = b.ForEach(func(_, v []byte) error {
//			var t Tag
//			err := json.Unmarshal(v, &t)
//			if err != nil {
//				return err
//			}
//			*ts = append(*ts, t)
//			return nil
//		})
//		return nil
//	})
//
//	sort.Slice(*ts, func(i, j int) bool {
//		return (*ts)[i].TfIdf > (*ts)[j].TfIdf
//	})
//
//	return nil
//}