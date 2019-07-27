package controllers

import (
	"github.com/julienschmidt/httprouter"

	"github.com/CheerChen/konachan-app/internal/kfile"
	"github.com/CheerChen/konachan-app/internal/models"

	"fmt"
	"net/http"
	"strconv"
)

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
