package kfile

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/cavaliercoder/grab"

	"github.com/CheerChen/konachan-app/internal/log"
)

func (pic *KFile) BuildName() {
	pic.Tags = strings.Replace(pic.Tags, "/", "_", -1)
	pic.Tags = strings.Replace(pic.Tags, ":", "_", -1)
	pic.Name = fmt.Sprintf("Konachan.com - %d %s%s", pic.Id, pic.Tags, pic.Ext)
}

const FileNameLengthLimit = 200

func (pic *KFile) SlimTags() {

	tags := strings.Split(pic.Tags, " ")

	for len(pic.Tags) >= FileNameLengthLimit {
		tags = tags[:len(tags)-1]
		pic.Tags = strings.Join(tags, " ")
	}
}

func DownloadFile(file *KFile, u string) {
	file.SlimTags()
	file.BuildName()

	filePath := AlbumPath + string(os.PathSeparator) + file.Name

	// create client
	client := grab.NewClient()
	req, _ := grab.NewRequest(filePath, u)

	//dialer, _ := proxy.SOCKS5("tcp", "127.0.0.1:1080", nil, proxy.Direct)
	//
	//client.HTTPClient.Transport = &http.Transport{
	//	DialContext: func(ctx context.Context, network, addr string) (conn net.Conn, e error) {
	//		c, e := dialer.Dial(network, addr)
	//		return c, e
	//	},
	//	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	//}

	// start download
	log.Infof("downloading %v...", req.URL())
	resp := client.Do(req)
	if resp.HTTPResponse != nil {
		log.Infof("HTTPResponse %v", resp.HTTPResponse.Status)
	}

	// start UI loop
	t := time.NewTicker(500 * time.Millisecond)
	defer t.Stop()

Loop:
	for {
		select {
		case <-t.C:
			log.Infof("  transferred %v / %v bytes (%.2f%%)",
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
		log.Errorf("download: %s", err)
		return
	}

	log.Infof("save to ./%s", resp.Filename)

	return
}
