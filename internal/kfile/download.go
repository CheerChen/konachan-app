package kfile

import (
	"fmt"
	"path"
	"regexp"
	"strings"

	"github.com/CheerChen/konachan-app/internal/grabber"
	"github.com/CheerChen/konachan-app/internal/log"

	"github.com/cavaliercoder/grab"
)

const FileNameLengthLimit = 200

func (pic *KFile) BuildName(u string) {
	// reduce tags
	tags := strings.Split(pic.Tags, " ")

	for len(pic.Tags) >= FileNameLengthLimit {
		tags = tags[:len(tags)-1]
		pic.Tags = strings.Join(tags, " ")
	}

	// replace special char
	var re = regexp.MustCompile(`[\\/:*?""<>|]`)
	pic.Tags = re.ReplaceAllString(pic.Tags, `$1`)

	if strings.Contains(u, "png") {
		pic.Ext = "png"
	}
	pic.Ext = "jpg"

	pic.Name = fmt.Sprintf("Konachan.com - %d %s.%s", pic.Id, pic.Tags, pic.Ext)
}

func DownloadFile(file *KFile, u string) {
	file.BuildName(u)
	log.Infof("building name %s...", file.Name)
	idxStr := fmt.Sprintf("%02d", file.Id/10000)
	dst := path.Join(WallpaperPath, idxStr, file.Name)
	g := grabber.NewDownloadClient()
	// create client
	req, _ := grab.NewRequest(dst, u)
	log.Infof("downloading %v...", req.URL())
	g.BatchDownload([]*grab.Request{req})

	log.Infof("save to ./%s", req.Filename)

	return
}
