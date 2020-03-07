package konachan

import (
	"github.com/dghubble/sling"
	"github.com/google/go-querystring/query"
	"net/http"
	"net/url"
)

type Tag struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Count     int    `json:"count"`
	Type      int    `json:"type"`
	Ambiguous bool   `json:"ambiguous"`
}

type TagService struct {
	sling *sling.Sling
}

func newTagService(sling *sling.Sling) *TagService {
	return &TagService{
		sling: sling,
	}
}

// limit How many tags to retrieve. Setting this to 0 will return every tag.
// order Can be date, count, or name.
// name The exact name of the tag.
// name_pattern Search for any tag that has this parameter in its name.
type TagListParams struct {
	Limit       int64  `url:"limit,omitempty"`
	Page        int64  `url:"page,omitempty"`
	Order       string `url:"order,omitempty"`
	Name        string `url:"name,omitempty"`
	NamePattern string `url:"name_pattern,omitempty"`
}

// Tag list
// https://konachan.com/help/api
func (s *TagService) List(params *TagListParams) ([]Tag, *http.Response, error) {
	tags := new([]Tag)
	apiError := new(APIError)
	resp, err := s.sling.New().Get("tag.json").QueryStruct(params).Receive(tags, apiError)
	return *tags, resp, err
}

func (s *TagService) ListUrlEncode(params *TagListParams) (string, error) {
	urlValues, err := url.ParseQuery("tag.json")
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
