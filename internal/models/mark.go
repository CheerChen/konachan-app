package models

import (
	"strings"

	"github.com/CheerChen/konachan-app/internal/humanize"
	"github.com/CheerChen/konachan-app/internal/log"
)

// Mark
func (p *Post) Mark(tfIdf, avgMap map[string]float64) {
	// 相似度打分
	tags := strings.Split(p.Tags, " ")

	//var tfIdfSum float64
	//if len(tags) > 2 {
	//	for _, tag := range tags {
	//		if _, ok := tfIdf[tag]; ok {
	//			tfIdfSum = tfIdfSum + tfIdf[tag]
	//		}
	//	}
	//}
	//p.TfIDf = tfIdfSum / float64(len(tags))

	//p.MyScore = (tfIdfSum + float64(p.Score)/avgMap[p.Rating]) / float64(len(tags)+1)
	//p.MyScore = (tfIdfSum + math.Log(float64(p.Score+1)/avgMap[p.Rating])) / float64(len(tags)+1)
	for _, t := range tfIdf {
		p.MyScore += t
	}
	for _, tag := range tags {
		if t, ok := tfIdf[tag]; ok {
			p.TfIDf += t
		}
	}
	p.Score = float64(p.Score) / avgMap[p.Rating]
	p.MyScore = p.TfIDf + p.Score

	_ = p.SortTagsByTfIdf(tfIdf)

	p.Size = humanize.Bytes(uint64(p.FileSize))
}

// Mark
func (ps *Posts) Mark(tfIdf map[string]float64) error {
	idMap, err := (*ps).FetchAllId()
	if err != nil {
		log.Errorf("fetch all post id: %s", err)
		return err
	}

	// 根据分级打平均分
	avgMap := make(map[string]float64)
	sumMap := make(map[string]float64)
	lenMap := make(map[string]int)
	for _, post := range *ps {
		if _, ok := sumMap[post.Rating]; !ok {
			sumMap[post.Rating] = post.Score
		} else {
			sumMap[post.Rating] += post.Score
		}

		if _, ok := lenMap[post.Rating]; !ok {
			lenMap[post.Rating] = 1
		} else {
			lenMap[post.Rating] += 1
		}
	}

	for ranting, sum := range sumMap {
		if l, ok := lenMap[ranting]; ok {
			avgMap[ranting] = float64(sum) / float64(l)
		}
	}

	log.Infof("create lenMap: %v", lenMap)
	log.Infof("create avgMap: %v", avgMap)

	for k := range *ps {

		(*ps)[k].Mark(tfIdf, avgMap)

		if _, ok := idMap[(*ps)[k].ID]; ok {
			(*ps)[k].IsFav = idMap[(*ps)[k].ID]
		}
	}

	return nil
}
