package controllers

import (
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"

	"github.com/CheerChen/konachan-app/internal/kfile"
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
	pageSize, page, err := GetPager(w, ps)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
	}

	posts := models.Work(ps.ByName("tag")[1:], pageSize, page)

	log.Infof("fetch posts: %d", len(posts))

	if len(posts) == 0 {
		http.Error(w, "no posts", http.StatusNotFound)
		return
	}

	idMap := make(map[int64]int)
	i := 0
	for _, post := range posts {
		if _, ok := idMap[post.ID]; ok {
			continue
		}
		idMap[post.ID] = 1
		posts[i] = post
		i++
	}
	posts = posts[0:i]

	log.Infof("dedup posts: %d", len(posts))

	tfIdf := models.GetTfIdf()
	err = posts.Mark(tfIdf)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
	}

	cJson(w, posts, map[string]int{
		"total": len(posts),
	})
	return
}

func Download(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	err = post.Find(post.ID)
	if err != nil {
		err = post.Save()
		if err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}
	}

	go kfile.DownloadFile(&kfile.KFile{Id: post.ID, Tags: post.Tags, Ext: post.GetFileExt()}, post.FileURL)

	cJson(w, "OK", nil)
	return
}
