package models

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/CheerChen/konachan-app/internal/service/konachan"
)

type KFile struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
	Ext  string `json:"ext"`
	Size int64
	Tags string `json:"tags"`

	Header string `json:"header"`
}

type KFiles []KFile

const FileNameLengthLimit = 200

func LoadFiles(path string) (pics KFiles) {
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		id := regexp.MustCompile(`[[:digit:]]+`).FindString(info.Name())
		if id == "" {
			return nil
		}
		var pic KFile
		pic.Id, err = strconv.ParseInt(id, 10, 64)
		if err != nil {
			return nil
		}
		pic.Ext = filepath.Ext(info.Name())
		pic.Name = path
		pic.Size = info.Size()

		if strings.HasSuffix(pic.Name, ".png") {
			pic.Header = "image/png"
		} else if strings.HasSuffix(pic.Name, ".jpg") {
			pic.Header = "image/jpeg"
		} else if strings.HasSuffix(pic.Name, ".gif") {
			pic.Header = "image/gif"
		} else {
			return nil
		}

		pics = append(pics, pic)

		return nil
	})
	if err != nil {
		fmt.Printf("walk error [%v]\n", err)
	}

	return
}

func GetFileById(id int64) (pic KFile, err error) {
	pics := LoadFiles(WallpaperPath)

	for _, p := range pics {
		if p.Id == id {
			return p, nil
		}
	}

	return pic, errors.New("record not found")
}

func DeleteFile(id int64) (err error) {
	pic, err := GetFileById(id)

	if err != nil {
		return
	}
	return os.Remove(pic.Name)

}

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

	log.Infof("downloading %v...", u)
	log.Infof("save to ./%s", dst)
	err = konachan.ParallelDownload(u, dst)
	if err != nil {
		log.Errorf("ParallelDownload err:", err)
	}
	return
}
