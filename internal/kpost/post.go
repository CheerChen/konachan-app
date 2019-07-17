package kpost

import (
	"strings"

	"github.com/jmoiron/sqlx"

	"github.com/CheerChen/konachan-app/internal/dbstore"
)

type KPost struct {
	ID        int    `json:"id" db:"id"`
	Tags      string `json:"tags" db:"tags"`
	CreatedAt int    `json:"created_at" db:"created_at"`
	//CreatorID           int           `json:"creator_id"`
	//Author              string        `json:"author"`
	//Change              int           `json:"change"`
	Source string `json:"source" db:"source"`
	Score  int    `json:"score" db:"score"`
	Md5    string `json:"md5" db:"md5"`

	FileSize int64  `json:"file_size" db:"file_size"`
	FileURL  string `json:"file_url" db:"file_url"`
	//IsShownInIndex      bool          `json:"is_shown_in_index"`

	PreviewURL          string `json:"preview_url"`
	PreviewWidth        int    `json:"preview_width"`
	PreviewHeight       int    `json:"preview_height"`
	ActualPreviewWidth  int    `json:"actual_preview_width"`
	ActualPreviewHeight int    `json:"actual_preview_height"`

	SampleURL      string `json:"sample_url"`
	SampleWidth    int    `json:"sample_width"`
	SampleHeight   int    `json:"sample_height"`
	SampleFileSize int    `json:"sample_file_size"`

	JpegURL      string `json:"jpeg_url"`
	JpegWidth    int    `json:"jpeg_width"`
	JpegHeight   int    `json:"jpeg_height"`
	JpegFileSize int    `json:"jpeg_file_size"`
	Rating       string `json:"rating"  db:"rating"`
	//HasChildren         bool          `json:"has_children"`
	ParentID interface{} `json:"parent_id"`
	Status   string      `json:"status"`
	Width    int         `json:"width"  db:"width"`
	Height   int         `json:"height"  db:"height"`
	//IsHeld              bool          `json:"is_held"`
	//FramesPendingString string        `json:"frames_pending_string"`
	//FramesPending       []interface{} `json:"frames_pending"`
	//FramesString        string        `json:"frames_string"`
	//Frames              []interface{} `json:"frames"`
	//FlagDetail          interface{}   `json:"flag_detail"`

	// 非官方字段
	TfIDf       float64 `json:"tf_idf"`
	MyScore     float64 `json:"my_score"`
	URL         string  `json:"url"`
	DownloadUrl string  `json:"download_url"`
	Size        string  `json:"size"`
	Has         bool    `json:"has"`
}

type KPosts []KPost

func (post KPost) GetFileExt() string {
	if strings.Contains(post.FileURL, "png") {
		return ".png"
	}
	return ".jpg"
}

func InitDB() {
	db = dbstore.GetDB()
}

var db *sqlx.DB

func (post KPost) Sync2DB() (err error) {

	_, err = db.NamedExec("INSERT INTO posts (id,tags,created_at,source,score,md5,file_size,file_url,rating,width,height) VALUES (:id,:tags,:created_at,:source,:score,:md5,:file_size,:file_url,:rating,:width,:height)",
		map[string]interface{}{
			"id":         post.ID,
			"tags":       post.Tags,
			"created_at": post.CreatedAt,
			"source":     post.Source,
			"score":      post.Score,
			"md5":        post.Md5,
			"file_size":  post.FileSize,
			"file_url":   post.FileURL,
			"rating":     post.Rating,
			"width":      post.Width,
			"height":     post.Height,
		})
	return
}

func (post *KPost) Find(id int) (err error) {
	err = db.Get(post, "SELECT * FROM posts where id = ?", id)

	return
}

func DeletePost(id int) (err error) {
	_, err = db.Exec("DELETE FROM posts where id = ?", id)

	return
}

func SelectAllTags() (tags []string, err error) {
	err = db.Select(&tags, "SELECT tags FROM posts")

	return
}

func SelectAllIds2Map() (idMap map[int]bool, err error) {
	var ids []int
	err = db.Select(&ids, "SELECT id FROM posts")

	if err == nil {
		idMap = make(map[int]bool)
		for _, id := range ids {
			idMap[id] = true
		}
	}
	return
}

func SelectPostByPage(limit, p int) (result KPosts, err error) {
	err = db.Select(&result, "SELECT * FROM posts order by id desc limit ?,?", (p-1)*limit, limit)

	return result, err
}

func SelectPostByTag(tag string) (result KPosts, err error) {
	err = db.Select(&result, "SELECT * FROM posts where tags like '%"+tag+"%'")

	return result, err
}

func SelectPostByPrefix(p string) (result KPosts, err error) {
	err = db.Select(&result, "SELECT * FROM posts where id like '"+p+"%'")

	return result, err
}
