package kfile

import (
	"fmt"
	"github.com/CheerChen/konachan-app/internal/service/konachan"
	"path"
	"regexp"
	"strings"

	"github.com/CheerChen/konachan-app/internal/log"
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

// DownloadFile ...
func DownloadFile(file *KFile, u string) {
	file.BuildName(u)
	log.Infof("building name %s...", file.Name)
	idxStr := fmt.Sprintf("%02d", file.Id/10000)
	err := EnsureDir(path.Join(WallpaperPath, idxStr))
	if err != nil {
		log.Errorf("EnsureDir err:", err)
		return
	}
	dst := path.Join(WallpaperPath, idxStr, file.Name)

	err = konachan.ParallelDownload(u, dst)
	if err != nil {
		log.Errorf("ParallelDownload err:", err)
	}
	return
}
