package models

import (
	"strings"

	"github.com/CheerChen/konachan-app/internal/log"
)

// Mark
func (p *Post) Mark(tfIdf map[string]float64, avg float64) {
	// 相似度打分
	tags := strings.Split(p.Tags, " ")
	score := 0.0
	if len(tags) > 2 {
		for _, tag := range tags {
			if _, ok := tfIdf[tag]; ok {
				score = score + tfIdf[tag]
			}
		}
	}

	p.TfIDf = score / float64(len(tags))
	// 对限制内容降权
	var userScore float64
	if p.Rating == "e" {
		userScore = float64(p.Score) * 0.618
	}

	userScore = float64(p.Score) / avg
	//
	//if userScore < 1 {
	//	userScore = userScore - 1
	//}

	p.MyScore = (score + userScore) / float64(len(tags)+1)

	sorted := SortTagsByTfIdf(tags, tfIdf)
	p.Tags = strings.Join(sorted, " ")
}

// Mark
func (ps *Posts) Mark(tfIdf map[string]float64) error {
	idMap, err := (*ps).FetchAllId()
	if err != nil {
		log.Errorf("fetch all post id: %s", err)
		return err
	}

	var sum int
	var avg float64
	for _, post := range *ps {
		sum += post.Score
	}
	avg = float64(sum) / float64(len(*ps))
	for k := range *ps {

		(*ps)[k].Mark(tfIdf, avg)

		if _, ok := idMap[(*ps)[k].ID]; ok {
			(*ps)[k].IsFav = idMap[(*ps)[k].ID]
		}
	}

	return nil
}
