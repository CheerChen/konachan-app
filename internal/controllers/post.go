package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"

	"github.com/CheerChen/konachan-app/internal/humanize"
	"github.com/CheerChen/konachan-app/internal/kpost"
	"github.com/CheerChen/konachan-app/internal/parallel"
)

func Popular(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	limit, err := strconv.Atoi(ps.ByName("limit"))

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}

	page, err := strconv.Atoi(ps.ByName("page"))

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}

	if limit <= 0 || page <= 0 {
		http.Error(w, http.StatusText(http.StatusNotAcceptable), http.StatusNotAcceptable)
	}

	tfIdf := kpost.GetTfIdf()

	posts := parallel.Work(0, limit, page)

	log.Println("fetch posts:")
	log.Println(len(posts))

	if len(posts) == 0 {
		http.Error(w, "no posts", http.StatusNotFound)
		return
	}
	posts = posts.FilterDeleted()
	posts = posts.FilterTags()

	reduced := posts.MarkAndReduce(0.418, tfIdf)
	//reduced.MarkDownloaded()
	//sort.Slice(reduced, func(i, j int) bool {
	//	return reduced[i].Score > reduced[j].Score
	//})

	cJson(w, reduced, map[string]int{
		"total":   len(posts),
		"reduced": len(reduced),
	})
	return
}

func PopularByRange(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	from, err := strconv.Atoi(ps.ByName("from"))

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}

	limit, err := strconv.Atoi(ps.ByName("limit"))

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}

	page, err := strconv.Atoi(ps.ByName("page"))

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}

	if limit <= 0 || page <= 0 {
		http.Error(w, http.StatusText(http.StatusNotAcceptable), http.StatusNotAcceptable)
	}

	tfIdf := kpost.GetTfIdf()

	posts := parallel.Work(from, limit, page)

	log.Println("fetch posts:")
	log.Println(len(posts))

	if len(posts) == 0 {
		http.Error(w, "no posts", http.StatusNotFound)
		return
	}
	posts = posts.FilterDeleted()
	posts = posts.FilterTags()

	reduced := posts.MarkAndReduce(0.418, tfIdf)
	//reduced.MarkDownloaded()

	cJson(w, reduced, map[string]int{
		"total":   len(posts),
		"reduced": len(reduced),
	})
	return
}

func GetByIdV2(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("id"))

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}
	post, err := kpost.GetPostByIdV2(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	tfIdf := kpost.GetTfIdf()

	post.Mark(tfIdf, 100)
	post.Size = humanize.Bytes(uint64(post.FileSize))

	post.URL = fmt.Sprintf("https://konachan.com/post/show/%d", post.ID)
	post.DownloadUrl = fmt.Sprintf("http://localhost:8080/download/%d", post.ID)
	cJson(w, post, nil)
}

func GetById(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}
	post, err := kpost.GetPostById(id, 100)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	tfIdf := kpost.GetTfIdf()

	post.Mark(tfIdf, 100)
	post.Size = humanize.Bytes(uint64(post.FileSize))

	post.URL = fmt.Sprintf("https://konachan.com/post/show/%d", post.ID)
	post.DownloadUrl = fmt.Sprintf("http://localhost:8080/download/%d", post.ID)
	cJson(w, post, nil)
}
