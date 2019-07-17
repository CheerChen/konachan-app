package kpost

import (
	"log"
	"math"
	"sort"
	"strings"
)

//var TfIdf map[string]float64

func GetTfIdf() map[string]float64 {
	//if len(TfIdf) != 0 {
	//	return TfIdf
	//}
	globalTotal, globalTagCount := getGlobalTagCount()

	allTags, err := SelectAllTags()
	if err != nil {
		log.Fatal("init db tags failed", err)
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

	log.Println("available tags:")
	log.Println(len(tfIdf))

	return tfIdf
}

func SortByTfIdf(unsorted []string, tfIdf map[string]float64) (sorted []string) {
	var tags Tags
	for _, item := range unsorted {
		if _, ok := tfIdf[item]; !ok {
			tfIdf[item] = 0.0
		}
		tags = append(tags, Tag{Name: item, TfIdf: tfIdf[item]})
	}

	sort.Slice(tags, func(i, j int) bool {
		return tags[i].TfIdf > tags[j].TfIdf
	})

	for _, tag := range tags {
		sorted = append(sorted, tag.Name)
	}
	return sorted
}
