package kfile

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strconv"
	"sync"

	"github.com/CheerChen/konachan-app/internal/log"
	"github.com/CheerChen/konachan-app/internal/models"
)

var ErrRecordNotFound = errors.New("record not found")

var AlbumPath = ""

func Reduce(p string) {
	err := ensureDir("temp")
	if err != nil {
		log.Warnf("ensureDir temp err:", err)
		return
	}
	files, err := ioutil.ReadDir(p)
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
		err = ensureDir(path.Join(p, idxStr))
		if err != nil {
			log.Warnf("ensureDir err:", err)
			continue
		}
		err = os.Rename(
			path.Join(p, f.Name()),
			path.Join(p, idxStr, f.Name()),
		)
		if err != nil {
			log.Warnf("Rename err:", err)
			continue
		}
		log.Infof("Move file: %s", path.Join(p, idxStr, f.Name()))
	}
}

func ensureDir(dirName string) error {
	err := os.Mkdir(dirName, os.ModeDir)
	if err == nil || os.IsExist(err) {
		return nil
	} else {
		return err
	}
}

func Sync(path string) {
	AlbumPath = path

	pics := LoadFiles(path)
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
		post, err = models.FetchPostByID(id)
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
