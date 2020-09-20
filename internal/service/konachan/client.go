package konachan

import (
	"net/http"

	"github.com/dghubble/sling"
)

const konachanAPI = "https://konachan.com/"

type Client struct {
	sling *sling.Sling
	Posts *PostService
	Tags  *TagService
}

func NewClient(httpClient *http.Client) *Client {
	base := sling.New().Client(httpClient).Base(konachanAPI)
	return &Client{
		sling: base,
		Posts: newPostService(base.New()),
		Tags:  newTagService(base.New()),
	}
}
