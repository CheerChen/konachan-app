package kfile

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type KFiles []KFile
type KFile struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
	Ext  string `json:"ext"`
	Size int64
	Tags string `json:"tags"`

	Header string `json:"header"`
}

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
	pics := LoadFiles(AlbumPath)

	for _, p := range pics {
		if p.Id == id {
			return p, nil
		}
	}

	return pic, ErrRecordNotFound
}

func Delete(id int64) (err error) {
	pic, err := GetFileById(id)

	if err != nil {
		return
	}
	return os.Remove(pic.Name)

}
