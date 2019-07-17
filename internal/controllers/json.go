package controllers

import (
	"encoding/json"
	"net/http"
)

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
		panic(err)
	}
}
