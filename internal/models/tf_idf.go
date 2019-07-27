package models

import (
	"math"
	"strings"

	"github.com/CheerChen/konachan-app/internal/log"
)

//var TfIdf map[string]float64

func GetTfIdf() map[string]float64 {
	//if len(TfIdf) != 0 {
	//	return TfIdf
	//}
	globalTotal, globalTagCount := GetGlobalTagCount()

	allTags, err := GetLocalTags()
	if err != nil {
		log.Fatalf("GetLocalTags failed", err)
	}
	tf1 := make(map[string]int)
	tf2 := make(map[string]int)

	for _, tags := range allTags {
		//if strings.Contains(file.Tags, "#") {
		//	file.Tags = strings.Replace(file.Tags, "#", "/", -1)
		//}
		ts := strings.Split(tags, " ")
		for _, tag := range ts {

			if _, ok := tf1[tag]; !ok {
				tf1[tag] = 1
			} else {
				tf1[tag] += 1
			}

			if _, ok := tf2[tag]; !ok {
				tf2[tag] = len(ts)
			} else {
				tf2[tag] += len(ts)
			}

		}
	}

	tfIdf := make(map[string]float64)

	for tag, tf1 := range tf1 {
		if _, ok := globalTagCount[tag]; !ok {
			globalTagCount[tag] = 1
		}
		idf := math.Log(float64(globalTotal) / (float64(globalTagCount[tag] + 1)))
		tf := float64(tf1) / float64(tf2[tag])
		tfIdf[tag] = tf * idf
	}

	// 降权
	tfIdf["nobody"] = 0.0
	tfIdf["all_male"] = 0.0

	log.Infof("available tags: %d", len(tfIdf))

	return tfIdf
}
