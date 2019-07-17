package parallel

import (
	"fmt"
	"sync"

	"github.com/CheerChen/konachan-app/internal/kpost"
)

func worker(pages <-chan int, c chan<- kpost.KPost) {
	for page := range pages {
		//
		var posts kpost.KPosts
		if offset != 0 {
			fmt.Println("start fetching offset", offset)
			posts = kpost.GetPostsByRange(offset, 100, page)
		} else {
			fmt.Println("start fetching page", page)
			posts = kpost.GetPostsByPage(100, page)
		}
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

var offset int

func Work(from, limit, p int) (result kpost.KPosts) {
	offset = from
	pageCh := headman(limit, p)
	resultCh := make(chan kpost.KPost)
	var wg sync.WaitGroup
	const numWorkers = 10
	wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go func() {
			worker(pageCh, resultCh)
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
	offset = 0
	return
}
