package models

import "strings"

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
	ParentID            interface{} `json:"parent_id"`
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

func (p *Post) Save() (err error) {
	return db.Save(p).Error
}

func (p *Post) Find(id int64) (err error) {
	return db.First(p, id).Error
}

func (p *Post) Delete() (err error) {
	return db.Delete(p).Error
}

func (ps *Posts) SelectTags() (err error) {
	return db.Select("tags").Find(&ps).Error
}

func (ps *Posts) SelectId() (err error) {
	return db.Select("id").Find(&ps).Error
}

func GetLocalTags() (tags []string, err error) {
	var posts Posts

	err = posts.SelectTags()
	if err != nil {
		return
	}

	for _, post := range posts {
		tags = append(tags, post.Tags)
	}
	return
}

func GetIdMap() (idMap map[int64]bool, err error) {
	idMap = make(map[int64]bool)
	var posts Posts

	err = posts.SelectId()
	if err != nil {
		return
	}

	if err == nil {
		for _, post := range posts {
			idMap[post.ID] = true
		}
	}
	return
}

func GetPostsByPage(l, p int) (ps Posts, err error) {
	return ps, db.Limit(l).Offset((p - 1) * l).Find(&ps).Error
}

func GetPostsByTag(tag string) (ps Posts, err error) {
	return ps, db.Where("tags LIKE ?", "%"+tag+"%").Find(&ps).Error
}

func (p Post) GetFileExt() string {
	if strings.Contains(p.FileURL, "png") {
		return ".png"
	}
	return ".jpg"
}
