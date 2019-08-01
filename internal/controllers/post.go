package controllers

import (
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"

	"github.com/CheerChen/konachan-app/internal/humanize"
	"github.com/CheerChen/konachan-app/internal/log"
	"github.com/CheerChen/konachan-app/internal/models"
)

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

func Remote(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	limit, page, err := GetPager(w, ps)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
	}

	posts := models.Work(ps.ByName("tag")[1:], limit, page)

	//posts := kpost.GetPosts(ps.ByName("tag"), limit, page)

	log.Infof("fetch posts: %d", len(posts))

	if len(posts) == 0 {
		http.Error(w, "no posts", http.StatusNotFound)
		return
	}

	tfIdf := models.GetTfIdf()
	reduced := posts.MarkAndReduce(0.0, tfIdf)

	cJson(w, reduced, map[string]int{
		"total":   len(posts),
		"reduced": len(reduced),
	})
	return
}
