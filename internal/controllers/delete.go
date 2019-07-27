package controllers

import (
	"net/http"
	"strconv"

	"github.com/CheerChen/konachan-app/internal/models"

	"github.com/julienschmidt/httprouter"
)

func Delete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.ParseInt(ps.ByName("id"), 10, 64)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}
	var post models.Post
	err = post.Find(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	err = post.Delete()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	//pics := kfile.LoadFiles()
	//if len(pics) == 0 {
	//	http.Error(w, "no pics", http.StatusNotFound)
	//	return
	//}
	//
	//for _, pic := range pics {
	//	if pic.Id == id {
	//		os.Remove(pic.Name)
	//	}
	//}
	return
}
