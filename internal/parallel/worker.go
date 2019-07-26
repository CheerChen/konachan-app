package parallel

import (
	"fmt"
	"sync"

	"github.com/CheerChen/konachan-app/internal/kpost"
)

func worker(pages <-chan int, c chan<- kpost.KPost, tags string) {
	for page := range pages {
		//
		var posts kpost.KPosts
		fmt.Println("start fetching page", page)
		posts = kpost.GetPosts(tags, 100, page)
		for _, post := range posts {
			c <- post
		}
	}
}

func headman(limit, p int) <-chan int {
	pages := make(chan int)
	go func() {
		defer close(pages)
		endPage := p + limit/100
		for p < endPage {
			pages <- p
			p += 1
		}
	}()
	return pages
}

func Work(tags string, limit, p int) (result kpost.KPosts) {
	pageCh := headman(limit, p)
	resultCh := make(chan kpost.KPost)
	var wg sync.WaitGroup
	const numWorkers = 10
	wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go func() {
			worker(pageCh, resultCh, tags)
			wg.Done()
		}()
	}
	go func() {
		wg.Wait()
		close(resultCh)
	}()

	for r := range resultCh {
		result = append(result, r)
	}
	return
}
