package models

import (
	"github.com/CheerChen/konachan-app/internal/service/konachan"
	"sync"

	"github.com/CheerChen/konachan-app/internal/log"
)

func fetchPosts(pages <-chan int, c chan<- Post, tags string) {
	for page := range pages {
		log.Infof("fetching page %d...", page)

		params := &konachan.PostListParams{
			Limit: 100,
			Page:  int64(page),
			Tags:  tags,
		}
		posts := Fetch(params)

		for _, post := range posts {
			c <- post
		}
	}
}

func headman(pageSize, p int) <-chan int {
	pages := make(chan int)
	go func() {
		defer close(pages)
		startPage := p
		endPage := p + (pageSize / 100)
		for startPage < endPage {
			pages <- startPage
			startPage += 1
		}
	}()
	return pages
}

func Work(tags string, ps, p int) (result Posts) {
	pageCh := headman(ps, p)
	resultCh := make(chan Post)
	var wg sync.WaitGroup
	const numWorkers = 10
	wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go func() {
			fetchPosts(pageCh, resultCh, tags)
			wg.Done()
		}()
	}
	go func() {
		wg.Wait()
		close(resultCh)
	}()

	idMap := make(map[int64]int)
	for r := range resultCh {
		if _, ok := idMap[r.ID]; ok {
			continue
		}
		idMap[r.ID] = 1
		result = append(result, r)
	}
	return
}
