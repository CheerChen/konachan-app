package konachan

import (
	"time"

	"github.com/CheerChen/konakore/internal/client"
	"github.com/CheerChen/konakore/internal/logger"
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
	lru      *LRUCache
	log      logger.Logger
)

func Set(c client.Client, host string, l logger.Logger) {
	hostname = host
	myclient = c
	lru = New(100)
	log = l
}
