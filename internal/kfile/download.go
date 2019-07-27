package kfile

import (
	"fmt"
	"strings"
	"time"

	"github.com/cavaliercoder/grab"

	"github.com/CheerChen/konachan-app/internal/log"
)

func (pic KFile) BuildName() string {
	pic.Tags = strings.Replace(pic.Tags, "/", "_", -1)
	pic.Tags = strings.Replace(pic.Tags, ":", "_", -1)
	return fmt.Sprintf("Konachan.com - %d %s%s", pic.Id, pic.Tags, pic.Ext)
}

func DownloadFile(name string, u string) {

	//filePath := FilePath + string(os.PathSeparator) + name
	filePath := name

	// create client
	client := grab.NewClient()
	req, _ := grab.NewRequest(filePath, u)

	//dialer, _ := proxy.SOCKS5("tcp", "127.0.0.1:10808",
	//	nil,
	//	&net.Dialer{
	//		Timeout:   10 * time.Second,
	//		KeepAlive: 10 * time.Second,
	//	},
	//)
	//client.HTTPClient.Transport = &http.Transport{
	//	Dial: dialer.Dial,
	//	//Proxy:           http.ProxyURL(proxy),
	//	//TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	//}

	// start download
	log.Infof("Downloading %v...\n", req.URL())
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
			log.Infof("  transferred %v / %v bytes (%.2f%%)\n",
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
		log.Errorf("Download failed: %v\n", err)
		return
	}

	log.Infof("Download saved to ./%v \n", resp.Filename)

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
