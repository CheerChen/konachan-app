package controllers

import (
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"

	"github.com/CheerChen/konakore/internal/models"
	"github.com/CheerChen/konakore/internal/service/konachan"
)

func GetByIdV2(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.ParseInt(ps.ByName("id"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}
	post := new(models.Post)
	kpost, err := konachan.GetPostByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	post.Make(kpost)
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
	posts := &models.Posts{}
	posts.Make(konachan.GetPosts(query, pageSize, page))

	log.Infof("fetch posts: %d", len(*posts))

	if len(*posts) != 0 {
		tfIdf, idf := getTfIdf()
		err = posts.Mark(tfIdf, idf)
		if err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
		}
		_ = posts.MarkExist()
	}

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

	post := new(models.Post)
	kpost, err := konachan.GetPostByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	post.Make(kpost)

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
		go models.DownloadFile(&models.KFile{Id: post.ID, Tags: post.Tags}, post.JpegURL)
	} else {
		go models.DownloadFile(&models.KFile{Id: post.ID, Tags: post.Tags}, post.FileURL)
	}

	cJson(w, "OK", nil)
	return
}
