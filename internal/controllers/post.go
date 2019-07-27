package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"

	"github.com/CheerChen/konachan-app/internal/humanize"
	"github.com/CheerChen/konachan-app/internal/models"
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

	tfIdf := models.GetTfIdf()

	posts := models.Work("", limit, page)

	log.Println("fetch posts:")
	log.Println(len(posts))

	if len(posts) == 0 {
		http.Error(w, "no posts", http.StatusNotFound)
		return
	}

	reduced := posts.MarkAndReduce(0.0, tfIdf)
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

func GetByIdV2(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.ParseInt(ps.ByName("id"), 10, 64)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}
	post, err := models.GetRemotePost(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	tfIdf := models.GetTfIdf()

	post.Mark(tfIdf, 100)
	post.Size = humanize.Bytes(uint64(post.FileSize))

	cJson(w, post, nil)
}

func Tag(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	//queryValues := r.URL.Query()

	tfIdf := models.GetTfIdf()

	posts := models.Work(ps.ByName("tag")[1:], limit, page)

	//posts := kpost.GetPosts(ps.ByName("tag"), limit, page)

	log.Println("fetch posts:")
	log.Println(len(posts))

	if len(posts) == 0 {
		http.Error(w, "no posts", http.StatusNotFound)
		return
	}

	reduced := posts.MarkAndReduce(0.0, tfIdf)

	cJson(w, reduced, map[string]int{
		"total":   len(posts),
		"reduced": len(reduced),
	})
	return
}
