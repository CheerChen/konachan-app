package kfile

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/cavaliercoder/grab"
	"golang.org/x/net/proxy"
)

func (pic KFile) BuildName() string {
	pic.Tags = strings.Replace(pic.Tags, "/", "_", -1)
	pic.Tags = strings.Replace(pic.Tags, ":", "_", -1)
	return fmt.Sprintf("Konachan.com - %d %s%s", pic.Id, pic.Tags, pic.Ext)
}

func DownloadFile(name string, u string) {

	filePath := FilePath + string(os.PathSeparator) + name

	// create client
	client := grab.NewClient()
	req, _ := grab.NewRequest(filePath, u)

	dialer, _ := proxy.SOCKS5("tcp", "127.0.0.1:10808",
		nil,
		&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 10 * time.Second,
		},
	)
	client.HTTPClient.Transport = &http.Transport{
		Dial: dialer.Dial,
		//Proxy:           http.ProxyURL(proxy),
		//TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	// start download
	fmt.Printf("Downloading %v...\n", req.URL())
	resp := client.Do(req)
	if resp.HTTPResponse != nil {
		fmt.Printf("  %v\n", resp.HTTPResponse.Status)
	}

	// start UI loop
	t := time.NewTicker(500 * time.Millisecond)
	defer t.Stop()

Loop:
	for {
		select {
		case <-t.C:
			fmt.Printf("  transferred %v / %v bytes (%.2f%%)\n",
				resp.BytesComplete(),
				resp.Size,
				100*resp.Progress())

		case <-resp.Done:
			// download is complete
			break Loop
		}
	}

	// check for errors
	if err := resp.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Download failed: %v\n", err)
		return
	}

	fmt.Printf("Download saved to ./%v \n", resp.Filename)

	// clean cache
	CleanFileCache()
	return
}

const FileNameLengthLimit = 200

func (pic *KFile) SlimTags() {

	tags := strings.Split(pic.Tags, " ")

	for len(pic.Tags) >= FileNameLengthLimit {
		tags = tags[:len(tags)-1]
		pic.Tags = strings.Join(tags, " ")
	}
}
