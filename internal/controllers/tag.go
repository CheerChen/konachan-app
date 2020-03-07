package controllers

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/CheerChen/konachan-app/internal/models"
)

func GetTfIdf(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	tfIdf, _ := models.GetTfIdf()
	cJson(w, tfIdf, nil)
	return
}
