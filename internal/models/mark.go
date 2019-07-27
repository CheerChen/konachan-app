package models

import (
	"sort"
	"strings"

	"github.com/CheerChen/konachan-app/internal/log"
)

// 给post打分
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

	if score > 0 {
		p.TfIDf = score / float64(len(tags))
		// 对限制内容降权
		var userScore float64
		if p.Rating == "e" {
			userScore = float64(p.Score) * 0.618 / avg
		} else {
			userScore = float64(p.Score) / avg
		}
		if userScore < 1 {
			userScore = userScore - 1
		}

		p.MyScore = (score + userScore) / float64(len(tags)+1)

		sorted := SortTagsByTfIdf(tags, tfIdf)
		p.Tags = strings.Join(sorted, " ")
	}
}

// 打分后按分数筛选排序
func (ps Posts) MarkAndReduce(baseline float64, tfIdf map[string]float64) (reduced Posts) {
	idMap, err := GetIdMap()
	if err != nil {
		log.Errorf("Get Id Map failed: %s", err)
		return ps
	}

	var sum int
	var avg float64
	for _, post := range ps {
		sum += post.Score
	}
	avg = float64(sum) / float64(len(ps))
	for _, post := range ps {

		post.Mark(tfIdf, avg)

		if post.MyScore > baseline {
			if _, ok := idMap[post.ID]; ok {
				post.IsFav = idMap[post.ID]
			}
			reduced = append(reduced, post)
		}
	}
	sort.Slice(reduced, func(i, j int) bool {
		return reduced[i].MyScore > reduced[j].MyScore
	})
	return
}

//func (ps Posts) FilterTags() (remain Posts) {
//	tags := []string{"moonknives"}
//
//	myfilter := func(p Post) bool {
//		combine := SimpleIntersect(strings.Split(p.Tags, " "), tags)
//		return len(combine.([]interface{})) == 0
//	}
//
//	return choose(ps, myfilter)
//}
//
//func (posts KPosts) FilterDeleted() (remain KPosts) {
//	myfilter := func(post KPost) bool {
//		return post.Status != "deleted"
//	}
//
//	return choose(posts, myfilter)
//}

//func choose(ss KPosts, test func(post KPost) bool) (ret KPosts) {
//	for _, s := range ss {
//		if test(s) {
//			ret = append(ret, s)
//		}
//	}
//	return
//}
//
//func SimpleIntersect(a interface{}, b interface{}) interface{} {
//	set := make([]interface{}, 0)
//	av := reflect.ValueOf(a)
//
//	for i := 0; i < av.Len(); i++ {
//		el := av.Index(i).Interface()
//		if contains(b, el) {
//			set = append(set, el)
//		}
//	}
//
//	return set
//}

//func contains(a interface{}, e interface{}) bool {
//	v := reflect.ValueOf(a)
//
//	for i := 0; i < v.Len(); i++ {
//		if v.Index(i).Interface() == e {
//			return true
//		}
//	}
//	return false
//}

// 按权重给标签排序
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
