package models

import (
	"math"
	"strings"

	"github.com/CheerChen/konachan-app/internal/log"
)

func GetTfIdf() (map[string]float64, map[string]float64) {
	tags := GetRemoteTags()
	countMap := make(map[string]int)
	for _, tag := range tags {
		countMap[tag.Name] = tag.Count
	}
	log.Infof("tag count map: %d", len(countMap))

	lastId, err := GetLastId()
	if err != nil {
		log.Warnf("get lastid err: %s", err)
		lastId = 30 * 10000
	}
	log.Infof("post last id: %d", lastId)

	posts := new(Posts)
	pts, err := posts.FetchAllTags()
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
