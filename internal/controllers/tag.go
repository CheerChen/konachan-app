package controllers

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/CheerChen/konachan-app/internal/kpost"
)

func GetTfIdf(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	cJson(w, kpost.GetTfIdf(), nil)
	return
}
