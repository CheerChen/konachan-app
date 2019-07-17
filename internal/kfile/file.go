package kfile

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/CheerChen/konachan-app/internal/memstore"
)

type KFiles []KFile
type KFile struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Ext  string `json:"ext"`
	Size int64
	Tags string `json:"tags"`
}

const FilePath = "E:\\Wallpaper"

var (
	cache = memstore.NewMemoryStore(10 * time.Minute)
)

func LoadFiles() (pics KFiles) {

	if cache.Get("local_pics") == nil {
		err := filepath.Walk(FilePath, func(path string, info os.FileInfo, err error) error {
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
			pic.Id, err = strconv.Atoi(id)
			if err != nil {
				return nil
			}
			pic.Ext = filepath.Ext(info.Name())
			pic.Name = path
			pic.Size = info.Size()

			tags := strings.Split(strings.Replace(info.Name(), pic.Ext, "", 1), " ")
			if len(tags) > 3 {
				//if len(tags) > 9 {
				//	pic.Tags = strings.Join(tags[3:9], " ")
				//} else {
				pic.Tags = strings.Join(tags[3:], " ")
				//}

			}

			pics = append(pics, pic)

			return nil
		})
		if err != nil {
			fmt.Printf("walk error [%v]\n", err)
		} else {
			cache.Set("local_pics", pics)
		}

	} else {
		pics = cache.Get("local_pics").(KFiles)
	}

	return
}

func CleanFileCache() {
	cache.Del("local_pics")
}
