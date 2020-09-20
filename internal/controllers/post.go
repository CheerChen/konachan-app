package controllers

import (
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"

	"github.com/CheerChen/konachan-app/internal/grabber"
	"github.com/CheerChen/konachan-app/internal/kfile"
	"github.com/CheerChen/konachan-app/internal/log"
)

func GetByIdV2(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.ParseInt(ps.ByName("id"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}

	post, err := grabber.GetPostByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	tfIdf, idf := getTfIdf()
	post.Mark(tfIdf, idf, map[string]float64{})

	cJson(w, post, nil)
}

func Remote(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	pageSize, page, err := GetPager(w, ps)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
	}

	query := GetQuery("tag", ps)
	posts := grabber.GetPosts(query, pageSize, page)

	log.Infof("fetch posts: %d", len(*posts))

	if len(*posts) == 0 {
		http.Error(w, "no posts", http.StatusNotFound)
		return
	}

	tfIdf, idf := getTfIdf()
	err = posts.Mark(tfIdf, idf)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
	}
	_ = posts.MarkExist()

	cJson(w, posts, map[string]int{
		"total": len(*posts),
	})
	return
}

func Download(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.ParseInt(ps.ByName("id"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}

	post, err := grabber.GetPostByID(id)
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

	if post.JpegFileSize != 0 && post.FileSize > int64(post.JpegFileSize*10) {
		log.Warnf("Downloading from Jpeg URL: %s", post.JpegURL)
		go kfile.DownloadFile(&kfile.KFile{Id: post.ID, Tags: post.Tags}, post.JpegURL)
	} else {
		go kfile.DownloadFile(&kfile.KFile{Id: post.ID, Tags: post.Tags}, post.FileURL)
	}

	cJson(w, "OK", nil)
	return
}
