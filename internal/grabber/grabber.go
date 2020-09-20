package grabber

import (
	"net/http"
	"strings"
	"time"

	"github.com/CheerChen/konachan-app/internal/log"

	"github.com/cavaliercoder/grab"
)

type DownloadClient struct {
	*grab.Client
}

const Workers = 12
const Interval = 200 * time.Millisecond

func NewDownloadClient() *DownloadClient {
	grabClient := grab.NewClient()
	grabClient.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.88 Safari/537.36"
	return &DownloadClient{grabClient}
}

func (g *DownloadClient) SetProxy(client *http.Client) {
	g.Client.HTTPClient = client
}

func (g *DownloadClient) BatchDownload(reqs []*grab.Request) {
	respch := g.Client.DoBatch(Workers, reqs...)

	// start a ticker to update progress every 200ms
	t := time.NewTicker(Interval)
	retries := make([]*grab.Request, 0)

	// monitor downloads
	completed := 0
	inProgress := 0
	responses := make([]*grab.Response, 0)
	for completed < len(reqs) {
		select {
		case resp := <-respch:
			// a new response has been received and has started downloading
			// (nil is received once, when the channel is closed by grab)
			if resp != nil {
				responses = append(responses, resp)
			}

		case <-t.C:
			// clear lines
			//if inProgress > 0 {
			//	fmt.Printf("\n")
			//	fmt.Printf("\033[%dA\033[K", inProgress)
			//}

			// update completed downloads
			for i, resp := range responses {
				if resp != nil && resp.IsComplete() {
					// print final result
					if resp.Err() != nil {
						log.Errorf("Error downloading: %s", resp.Err())
						//fmt.Printf("Error downloading %s: %v\n", resp.Request.URL(), resp.Err())
						if strings.Contains(resp.Err().Error(), "EOF") {
							retry, _ := grab.NewRequest(resp.Filename, resp.Request.URL().String())
							retries = append(retries, retry)
						}
					} else {
						//	fmt.Printf("Finished %s %d / %d bytes (%d%%)\n", resp.Filename, resp.BytesComplete(), resp.Size, int(100*resp.Progress()))
						log.Infof("save to ./%s", resp.Filename)
					}

					// mark completed
					responses[i] = nil
					completed++
				}
			}

			// update downloads in progress
			inProgress = 0
			for _, resp := range responses {
				if resp != nil {
					inProgress++
					if resp.Size < 0 {
						resp.Size = 0
					}
					log.Infof("  transferred %v / %v bytes (%.2f%%)",
						resp.BytesComplete(),
						resp.Size,
						100*resp.Progress())
					//fmt.Printf("Downloading %s %d / %d bytes (%d%%)\033[K\n", resp.Filename, resp.BytesComplete(), resp.Size, int(100*resp.Progress()))
				}
			}
		}
	}

	t.Stop()

	//fmt.Printf("%d files successfully downloaded.\n", completed-len(retries))
	if len(retries) > 0 {
		//fmt.Printf("%d files retry downloading.\n", len(retries))
		g.BatchDownload(retries)
	}
}
