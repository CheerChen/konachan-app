package controllers

import (
	"github.com/julienschmidt/httprouter"

	"github.com/CheerChen/konachan-app/internal/kfile"
	"github.com/CheerChen/konachan-app/internal/kpost"

	"fmt"
	"net/http"
	"strconv"
)

func Download(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
	var dbPost kpost.KPost
	err = dbPost.Find(post.ID)
	if err != nil {
		err = post.Sync2DB()
		if err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}
	}

	var file kfile.KFile
	file.Id = post.ID
	file.Tags = post.Tags
	file.Ext = post.GetFileExt()
	file.SlimTags()

	//url := kfile.DownloadHelper(post.FileURL)
	go kfile.DownloadFile(file.BuildName(), post.FileURL)
	// auto close
	fmt.Fprintln(w, "<html><body><script>window.location.href=\"about:blank\";window.close();</script></body></html>")
	return
}
