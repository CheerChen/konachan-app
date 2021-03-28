package models

import (
	"math"
	"strings"
)

// Mark
func (ps *Posts) Mark(tfIdf, idf map[string]float64) error {
	avgMap := ps.AvgMap()
	for k := range *ps {
		(*ps)[k].Mark(tfIdf, idf, avgMap)
	}

	return nil
}

// 相似度打分
func (p *Post) Mark(tfIdf, idf, avgMap map[string]float64) {
	tags := strings.Split(p.Tags, " ")

	// version 1
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

	// version 2
	for _, tag := range tags {
		if t, ok := tfIdf[tag]; ok {
			p.TfIDf += t
		}
	}
	p.TfIDf = p.TfIDf / float64(len(tags))
	p.MyScore = p.TfIDf + math.Log(float64(p.Score+1)/avgMap[p.Rating])/float64(len(tags))

	_ = p.SortTagsByTfIdf(tfIdf)
}

func (ps *Posts) MarkExist() error {
	idMap, err := (*ps).FetchAllId()
	if err != nil {
		log.Errorf("fetch all post id: %s", err)
		return err
	}
	for k := range *ps {
		if _, ok := idMap[(*ps)[k].ID]; ok {
			(*ps)[k].IsFav = true
		}
	}

	return nil
}

// 根据分级打平均分
func (ps *Posts) AvgMap() map[string]float64 {
	avgMap := make(map[string]float64)
	sumMap := make(map[string]float64)
	lenMap := make(map[string]int)
	for _, post := range *ps {
		if _, ok := sumMap[post.Rating]; !ok {
			sumMap[post.Rating] = float64(post.Score)
		} else {
			sumMap[post.Rating] += float64(post.Score)
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
	log.Infof("created avgMap: %v", avgMap)

	return avgMap
}
