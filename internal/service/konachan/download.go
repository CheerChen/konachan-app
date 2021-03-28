package konachan

import (
	"github.com/CheerChen/konachan-app/internal/client"
)

func ParallelDownload(u, filename string) error {
	return myclient.Download(u, filename, client.DefaultProgress())
}
