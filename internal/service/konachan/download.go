package konachan

import (
	"github.com/CheerChen/konakore/internal/client"
)

func ParallelDownload(u, filename string) error {
	return myclient.Download(u, filename, client.DefaultProgress())
}
