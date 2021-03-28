package controllers

import (
	"math"
	"net/http"
	"strings"

	"github.com/CheerChen/konachan-app/internal/models"
	"github.com/CheerChen/konachan-app/internal/service/konachan"

	"github.com/julienschmidt/httprouter"
)

func GetTfIdf(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	tfIdf, _ := getTfIdf()
	cJson(w, tfIdf, nil)
	return
}

func getTfIdf() (map[string]float64, map[string]float64) {

	lastId, err := konachan.GetPostLastId()
	if err != nil {
		log.Warnf("get lastid err: %s", err)
		lastId = 30 * 10000
	}
	log.Infof("post last id: %d", lastId)

	localPosts := make(models.Posts, 0)
	pts, err := localPosts.FetchAllTags()
	if err != nil {
		log.Fatalf("fetch all tags: %s", err)
	}
	tf1 := make(map[string]int)
	tf2 := make(map[string]int)

	for _, pt := range pts {
		tags := strings.Split(pt, " ")
		for _, tag := range tags {

			if _, ok := tf1[tag]; !ok {
				tf1[tag] = 1
			} else {
				tf1[tag] += 1
			}

			if _, ok := tf2[tag]; !ok {
				tf2[tag] = len(tags)
			} else {
				tf2[tag] += len(tags)
			}

		}
	}

	tfIdf := make(map[string]float64)
	idfMap := make(map[string]float64)

	countMap := konachan.GetTags()
	for tag, tf1 := range tf1 {
		if _, ok := countMap[tag]; !ok {
			countMap[tag] = 1
		}
		idf := math.Log(float64(lastId) / (float64(countMap[tag] + 1)))
		tf := float64(tf1) / float64(tf2[tag])
		tfIdf[tag] = tf * idf
		idfMap[tag] = idf
	}
	log.Infof("available tags: %d", len(tfIdf))

	return tfIdf, idfMap
}
