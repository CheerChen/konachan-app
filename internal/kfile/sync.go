package kfile

import (
	"errors"
	"sync"

	"github.com/CheerChen/konachan-app/internal/log"
	"github.com/CheerChen/konachan-app/internal/models"
)

var ErrRecordNotFound = errors.New("record not found")

var AlbumPath = ""

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

	resultCh := make(chan models.Post)
	var wg sync.WaitGroup
	const numWorkers = 10
	wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go func() {
			worker(idCh, resultCh)
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

func worker(ids <-chan int64, c chan<- models.Post) {
	for id := range ids {
		var post models.Post
		err := post.Find(id)
		if err == nil {
			log.Infof("find ID(%d) in db", post.ID)
			continue
		}
		post, err = models.GetRemotePost(id)
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
