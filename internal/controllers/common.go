package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/CheerChen/konakore/internal/logger"

	"github.com/julienschmidt/httprouter"
)

func GetPager(w http.ResponseWriter, ps httprouter.Params) (int, int, error) {
	limit, err := strconv.Atoi(ps.ByName("limit"))

	if err != nil {
		return 0, 0, err
	}

	page, err := strconv.Atoi(ps.ByName("page"))

	if err != nil {
		return 0, 0, err
	}

	if limit > 500 || limit <= 0 || page <= 0 {
		return 0, 0, errors.New(http.StatusText(http.StatusNotAcceptable))
	}

	return limit, page, nil
}

func GetQuery(name string, ps httprouter.Params) string {
	return ps.ByName(name)[1:]
}

type JsonResponse struct {
	// Reserved field to add some meta information to the API response
	Meta interface{} `json:"meta"`
	Data interface{} `json:"data"`
}

func cJson(w http.ResponseWriter, data interface{}, meta interface{}) {
	response := &JsonResponse{Data: &data, Meta: &meta}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Errorf("json encode: %s", err)
	}
}

var (
	log logger.Logger
)

func Log(l logger.Logger) {
	log = l
}
