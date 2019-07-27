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

		tags := strings.Split(strings.Replace(info.Name(), pic.Ext, "", 1), " ")
		if len(tags) > 3 {
			pic.Tags = strings.Join(tags[3:], " ")
			pics = append(pics, pic)
		}

		return nil
	})
	if err != nil {
		fmt.Printf("walk error [%v]\n", err)
	}

	return
}
