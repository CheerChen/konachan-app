package models

import (
	"encoding/json"
	"sort"
	"strconv"
	"strings"

	"github.com/CheerChen/konachan-app/internal/log"

	"github.com/boltdb/bolt"
)

type Post struct {
	ID                  int64       `json:"id"`
	Tags                string      `json:"tags"`
	CreatedAt           int         `json:"created_at"`
	CreatorID           int         `json:"creator_id"`
	Author              string      `json:"author"`
	Change              int         `json:"change"`
	Source              string      `json:"source"`
	Score               int         `json:"score"`
	Md5                 string      `json:"md5"`
	FileSize            int64       `json:"file_size"`
	FileURL             string      `json:"file_url"`
	IsShownInIndex      bool        `json:"is_shown_in_index"`
	PreviewURL          string      `json:"preview_url"`
	PreviewWidth        int         `json:"preview_width"`
	PreviewHeight       int         `json:"preview_height"`
	ActualPreviewWidth  int         `json:"actual_preview_width"`
	ActualPreviewHeight int         `json:"actual_preview_height"`
	SampleURL           string      `json:"sample_url"`
	SampleWidth         int         `json:"sample_width"`
	SampleHeight        int         `json:"sample_height"`
	SampleFileSize      int         `json:"sample_file_size"`
	JpegURL             string      `json:"jpeg_url"`
	JpegWidth           int         `json:"jpeg_width"`
	JpegHeight          int         `json:"jpeg_height"`
	JpegFileSize        int         `json:"jpeg_file_size"`
	Rating              string      `json:"rating"  db:"rating"`
	HasChildren         bool        `json:"has_children"`
	ParentID            interface{} `json:"parent_id" gorm:"-"`
	Status              string      `json:"status"`
	Width               int         `json:"width"`
	Height              int         `json:"height"`
	//IsHeld              bool          `json:"is_held"`
	//FramesPendingString string        `json:"frames_pending_string"`
	//FramesPending       []interface{} `json:"frames_pending"`
	//FramesString        string        `json:"frames_string"`
	//Frames              []interface{} `json:"frames"`
	//FlagDetail          interface{}   `json:"flag_detail"`

	// 非官方字段
	TfIDf   float64 `json:"tf_idf"`
	MyScore float64 `json:"my_score"`
	IsFav   bool    `json:"is_fav"`
	Size    string  `json:"size"`
}

type Posts []Post

func (p Post) GetFileExt() string {
	if strings.Contains(p.FileURL, "png") {
		return ".png"
	}
	return ".jpg"
}

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
			return ErrRecordNotFound
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

	min, max, start, end := 0, len(*ps), (page-1)*l, page*l

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
	var tags Tags
	for _, tag := range strings.Split(p.Tags, " ") {
		if _, ok := tfIdf[tag]; !ok {
			tfIdf[tag] = 0.0
		}
		tags = append(tags, Tag{Name: tag, TfIdf: tfIdf[tag]})
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
