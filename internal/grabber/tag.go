package grabber

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/CheerChen/konachan-app/internal/log"
	"github.com/CheerChen/konachan-app/internal/models"

	"github.com/cavaliercoder/grab"
)

func GetTags() (tags []*models.Tag) {
	g := NewDownloadClient()
	g.SetProxy(proxyClient)
	urlStr := fmt.Sprintf(tagUrlFmt, hostname, tagLimit)
	dst := fmt.Sprintf(tagDstFmt, os.TempDir())
	req, _ := grab.NewRequest(dst, urlStr)

	info, err := os.Stat(req.Filename)
	if os.IsNotExist(err) {
		g.BatchDownload([]*grab.Request{req})
	} else if !os.IsNotExist(err) {
		if (time.Now().Sub(info.ModTime())) > cacheTime {
			_ = os.Remove(req.Filename)
			g.BatchDownload([]*grab.Request{req})
		}
	} else {
		log.Warnf("os.Stat err:", err)
		return
	}
	var b []byte
	b, err = ioutil.ReadFile(req.Filename)
	if err != nil {
		log.Warnf("ReadDir err:", err)
		return
	}
	err = json.Unmarshal(b, &tags)
	if err != nil {
		log.Warnf("Unmarshal err:", err)
		return
	}
	return
}
