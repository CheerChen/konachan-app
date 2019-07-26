package controllers

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/julienschmidt/httprouter"

	"github.com/CheerChen/konachan-app/internal/kfile"
	"github.com/CheerChen/konachan-app/internal/kpost"
)

func Delete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("id"))

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}

	err = kpost.DeletePost(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	pics := kfile.LoadFiles()
	if len(pics) == 0 {
		http.Error(w, "no pics", http.StatusNotFound)
		return
	}

	for _, pic := range pics {
		if pic.Id == id {
			os.Remove(pic.Name)
		}
	}

	// clean cache
	kfile.CleanFileCache()

	fmt.Fprintln(w, "<html><body><script>window.location.href=\"about:blank\";window.close();</script></body></html>")
	return
}
