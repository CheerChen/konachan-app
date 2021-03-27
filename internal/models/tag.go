package models

import (
	"encoding/json"
	bolt "go.etcd.io/bbolt"
	"sort"
	"strings"
)

type OriginalTag struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Count     int    `json:"count"`
	Type      int    `json:"type"`
	Ambiguous bool   `json:"ambiguous"`
}

type Tag struct {
	OriginalTag

	TfIdf float64 `json:"tf_idf"`
	Idf   float64 `json:"idf"`
}

func (ps *Posts) FetchAllTags() (pts []string, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(ps.TableName()))
		_ = b.ForEach(func(_, v []byte) error {
			var p Post
			err := json.Unmarshal(v, &p)
			if err != nil {
				return err
			}

			pts = append(pts, p.Tags)
			return nil
		})
		return nil
	})
	return
}

// SortTagsByTfIdf
func (p *Post) SortTagsByTfIdf(tfIdf map[string]float64) (err error) {
	var tags []*Tag
	for _, tag := range strings.Split(p.Tags, " ") {
		if _, ok := tfIdf[tag]; !ok {
			tfIdf[tag] = 0.0
		}
		tags = append(tags, &Tag{OriginalTag: OriginalTag{Name: tag}, TfIdf: tfIdf[tag]})
	}

	sort.Slice(tags, func(i, j int) bool {
		return tags[i].TfIdf > tags[j].TfIdf
	})

	var parts []string
	for _, tag := range tags {
		parts = append(parts, tag.Name)
	}
	p.Tags = strings.Join(parts, " ")
	return nil
}
