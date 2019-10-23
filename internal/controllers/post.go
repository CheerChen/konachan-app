package controllers

import (
	"net/http"
	"sort"
	"strconv"

	"github.com/julienschmidt/httprouter"

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

	post.Mark(tfIdf, map[string]float64{})

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
	err = posts.Mark(tfIdf)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
	}

	sort.Slice(posts, func(i, j int) bool {
		return posts[i].MyScore > posts[j].MyScore
	})

	cJson(w, posts, map[string]int{
		"total":   len(posts),
		"reduced": len(posts),
	})
	return
}
