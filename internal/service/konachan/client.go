package konachan

import (
	"github.com/dghubble/sling"
	"net/http"
)

const konachanAPI = "https://konachan.com/"

// Client is a Twitter client for making Twitter API requests.
type Client struct {
	sling *sling.Sling
	// Twitter API Services
	Posts *PostService
	Tags  *TagService
}

// NewClient returns a new Client.
func NewClient(httpClient *http.Client) *Client {
	base := sling.New().Client(httpClient).Base(konachanAPI)
	return &Client{
		sling: base,
		Posts: newPostService(base.New()),
		Tags:  newTagService(base.New()),
	}
}
