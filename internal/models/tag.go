package models

import (
	"sort"
)

type Tag struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Count     int    `json:"count"`
	Type      int    `json:"type"`
	Ambiguous bool   `json:"ambiguous"`

	// 某个tag对post的重要性越高，它的TF-IDF值就越大
	TfIdf float64 `json:"tf_idf"`
}

type Tags []Tag

// SortTagsByTfIdf
func SortTagsByTfIdf(unsorted []string, tfIdf map[string]float64) (sorted []string) {
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
