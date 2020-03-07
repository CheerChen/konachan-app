package konachan

import (
	"net/http"
	"net/url"

	"github.com/dghubble/sling"
	"github.com/google/go-querystring/query"
)

type Post struct {
	ID                  int64  `json:"id"`
	Tags                string `json:"tags"`
	CreatedAt           int    `json:"created_at"`
	CreatorID           int    `json:"creator_id"`
	Author              string `json:"author"`
	Change              int    `json:"change"`
	Source              string `json:"source"`
	Score               int    `json:"score"`
	Md5                 string `json:"md5"`
	FileSize            int64  `json:"file_size"`
	FileURL             string `json:"file_url"`
	IsShownInIndex      bool   `json:"is_shown_in_index"`
	PreviewURL          string `json:"preview_url"`
	PreviewWidth        int    `json:"preview_width"`
	PreviewHeight       int    `json:"preview_height"`
	ActualPreviewWidth  int    `json:"actual_preview_width"`
	ActualPreviewHeight int    `json:"actual_preview_height"`
	SampleURL           string `json:"sample_url"`
	SampleWidth         int    `json:"sample_width"`
	SampleHeight        int    `json:"sample_height"`
	SampleFileSize      int    `json:"sample_file_size"`
	JpegURL             string `json:"jpeg_url"`
	JpegWidth           int    `json:"jpeg_width"`
	JpegHeight          int    `json:"jpeg_height"`
	JpegFileSize        int    `json:"jpeg_file_size"`
	Rating              string `json:"rating"`
	Status              string `json:"status"`
	Width               int    `json:"width"`
	Height              int    `json:"height"`
	//HasChildren         bool        `json:"has_children"`
	//ParentID            interface{} `json:"parent_id"`
	//IsHeld              bool          `json:"is_held"`
	//FramesPendingString string        `json:"frames_pending_string"`
	//FramesPending       []interface{} `json:"frames_pending"`
	//FramesString        string        `json:"frames_string"`
	//Frames              []interface{} `json:"frames"`
	//FlagDetail          interface{}   `json:"flag_detail"`
}

type PostService struct {
	sling *sling.Sling
}

func newPostService(sling *sling.Sling) *PostService {
	return &PostService{
		sling: sling,
	}
}

type PostListParams struct {
	Limit int64  `url:"limit,omitempty"`
	Page  int64  `url:"page,omitempty"`
	Tags  string `url:"tags,omitempty"`
}

//
// https://konachan.com/help/api
func (s *PostService) List(params *PostListParams) ([]Post, *http.Response, error) {
	posts := new([]Post)
	apiError := new(APIError)
	resp, err := s.sling.New().Get("post.json").QueryStruct(params).Receive(posts, apiError)
	return *posts, resp, err
}

func (s *PostService) ListUrlEncode(params *PostListParams) (string, error) {
	urlValues, err := url.ParseQuery("post.json")
	if err != nil {
		return "", err
	}
	queryValues, err := query.Values(params)
	if err != nil {
		return "", err
	}
	for key, values := range queryValues {
		for _, value := range values {
			urlValues.Add(key, value)
		}
	}
	return urlValues.Encode(), nil
}
