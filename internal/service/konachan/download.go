package konachan

import (
	"github.com/CheerChen/konachan-app/internal/client"
	"github.com/CheerChen/konachan-app/internal/log"
)

func ParallelDownload(u, filename string) error {
	log.Infof("downloading %v...", u)
	log.Infof("save to ./%s", filename)
	return myclient.Download(u, filename, client.DefaultProgress())
}
