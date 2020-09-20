package kfile

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strconv"
	"sync"

	"github.com/CheerChen/konachan-app/internal/grabber"
	"github.com/CheerChen/konachan-app/internal/log"
	"github.com/CheerChen/konachan-app/internal/models"
)

var WallpaperPath = ""

func CheckPath(wp string) {
	if err := EnsureDir(wp); err != nil {
		log.Fatalf("Error reading path, %s", err)
	}
	if err := EnsureDir(os.TempDir()); err != nil {
		log.Fatalf("Error reading temp path, %s", err)
	}
	WallpaperPath = wp
	// 按 id 分布到文件夹
	Reduce()
	// 检查本地文件和数据库一致
	Sync()
}

func Reduce() {
	files, err := ioutil.ReadDir(WallpaperPath)
	if err != nil {
		log.Warnf("ReadDir err:", err)
		return
	}
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		log.Infof("Check file: %s", f.Name())
		TypeID, _ := regexp.Compile(`[0-9]+`)
		id, err := strconv.Atoi(TypeID.FindString(f.Name()))
		if err != nil {
			log.Warnf("strconv err:", err)
			continue
		}
		idx := id / 10000
		idxStr := fmt.Sprintf("%02d", idx)
		err = EnsureDir(path.Join(WallpaperPath, idxStr))
		if err != nil {
			log.Warnf("EnsureDir err:", err)
			continue
		}
		err = os.Rename(
			path.Join(WallpaperPath, f.Name()),
			path.Join(WallpaperPath, idxStr, f.Name()),
		)
		if err != nil {
			log.Warnf("Rename err:", err)
			continue
		}
		log.Infof("Move file: %s", path.Join(WallpaperPath, idxStr, f.Name()))
	}
}

func EnsureDir(dirName string) error {
	err := os.Mkdir(dirName, os.ModeDir)
	if err == nil || os.IsExist(err) {
		return nil
	} else {
		return err
	}
}

func Sync() {
	pics := LoadFiles(WallpaperPath)
	if len(pics) == 0 {
		log.Warnf("Wallpaper path empty!")
		return
	}

	idCh := make(chan int64)
	go func() {
		defer close(idCh)
		for _, pic := range pics {
			idCh <- pic.Id
		}
	}()

	resultCh := make(chan *models.Post)
	var wg sync.WaitGroup
	const numWorkers = 10
	wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go func() {
			syncPost(idCh, resultCh)
			wg.Done()
		}()
	}
	go func() {
		wg.Wait()
		close(resultCh)
	}()

	var downloaded int
	for r := range resultCh {
		downloaded++
		log.Infof("Sync SUCCESS to db ID(%d)", r.ID)
	}

	log.Infof("Sync complete")
	log.Infof("Synced: %d", downloaded)

	return
}

func syncPost(ids <-chan int64, c chan<- *models.Post) {
	for id := range ids {
		post := new(models.Post)
		err := post.Find(id)
		if err == nil {
			log.Infof("find ID(%d) in db", post.ID)
			continue
		}
		post, err = grabber.GetPostByID(id)
		if err != nil {
			log.Warnf("fetch ID(%d) from web: %s", id, err.Error())
			continue
		}
		err = post.Save()
		if err != nil {
			log.Errorf("save post ID(%d): %s", post.ID, err.Error())
			continue
		}

		c <- post
	}
}
