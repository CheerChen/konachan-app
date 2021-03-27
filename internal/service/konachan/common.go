package konachan

import (
	"github.com/CheerChen/konachan-app/internal/client"
	"time"
)

const (
	postDstFmt = "%s/%s_%d.json"
	postLimit  = 100
	postUrlFmt = "%s/post.json?limit=%d&page=%d&tags=%s"
	tagDstFmt  = "%s/tag.json"
	tagLimit   = 10000
	tagUrlFmt  = "%s/tag.json?limit=%d&order=count"

	cacheTime   = 1 * time.Hour
	hourLayout  = "2006-01-02 15"
	monthLayout = "2006-01"
)

var (
	hostname string
	myclient client.Client
	lastid   int64
	lru   *LRUCache
)

func SetHost(host string) {
	hostname = host
}

func SetClient(c client.Client) {
	myclient = c
	lru = New(100)
}