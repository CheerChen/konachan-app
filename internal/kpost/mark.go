package kpost

import (
	"fmt"
	"log"
	"reflect"
	"sort"
	"strings"
)

// 给post打分
func (post *KPost) Mark(tfIdf map[string]float64, avg float64) {
	// 相似度打分
	tags := strings.Split(post.Tags, " ")
	score := 0.0
	if len(tags) > 2 {
		for _, tag := range tags {
			if _, ok := tfIdf[tag]; ok {
				score = score + tfIdf[tag]
			}
		}
	}

	if score > 0 {

		post.TfIDf = score / float64(len(tags))

		// 对限制内容降权
		var userScore float64
		if post.Rating == "e" {
			//	userScore = float64(post.Score) * (1 - 0.618) / 100
			//} else if post.Rating == "q" {
			userScore = float64(post.Score) * 0.618 / avg
		} else {
			userScore = float64(post.Score) / avg
		}
		if userScore < 1 {
			userScore = userScore - 1
		}

		post.MyScore = (score + userScore) / float64(len(tags)+1)

		sorted := SortByTfIdf(tags, tfIdf)
		post.Tags = strings.Join(sorted, " ")
	}
}

// 打分后按分数筛选排序
func (posts KPosts) MarkAndReduce(baseline float64, tfIdf map[string]float64) (reduced KPosts) {
	idMap, err := SelectAllIds2Map()
	if err != nil {
		log.Println(idMap)
		log.Println(err)
	}

	var sum int
	var avg float64
	for _, post := range posts {
		sum += post.Score
	}
	avg = float64(sum) / float64(len(posts))
	for _, post := range posts {

		post.Mark(tfIdf, avg)

		if post.MyScore > baseline {
			if _, ok := idMap[post.ID]; ok {
				post.Has = idMap[post.ID]
			}
			post.URL = fmt.Sprintf("https://konachan.com/post/show/%d", post.ID)
			reduced = append(reduced, post)
		}
	}
	sort.Slice(reduced, func(i, j int) bool {
		return reduced[i].MyScore > reduced[j].MyScore
	})
	return
}

func (posts KPosts) MarkNotReduce(tfIdf map[string]float64) (marked KPosts) {
	var sum int
	var avg float64
	for _, post := range posts {
		sum += post.Score
	}
	avg = float64(sum) / float64(len(posts))
	for _, post := range posts {

		post.Mark(tfIdf, avg)
		post.URL = fmt.Sprintf("https://konachan.com/post/show/%d", post.ID)

		marked = append(marked, post)
	}
	sort.Slice(marked, func(i, j int) bool {
		return marked[i].MyScore > marked[j].MyScore
	})
	return
}

func (posts KPosts) FilterTags() (remain KPosts) {
	tags := []string{"moonknives"}

	myfilter := func(post KPost) bool {
		combine := SimpleIntersect(strings.Split(post.Tags, " "), tags)
		return len(combine.([]interface{})) == 0
	}

	return choose(posts, myfilter)
}

func (posts KPosts) FilterDeleted() (remain KPosts) {
	myfilter := func(post KPost) bool {
		return post.Status != "deleted"
	}

	return choose(posts, myfilter)
}

func choose(ss KPosts, test func(post KPost) bool) (ret KPosts) {
	for _, s := range ss {
		if test(s) {
			ret = append(ret, s)
		}
	}
	return
}

func SimpleIntersect(a interface{}, b interface{}) interface{} {
	set := make([]interface{}, 0)
	av := reflect.ValueOf(a)

	for i := 0; i < av.Len(); i++ {
		el := av.Index(i).Interface()
		if contains(b, el) {
			set = append(set, el)
		}
	}

	return set
}

func contains(a interface{}, e interface{}) bool {
	v := reflect.ValueOf(a)

	for i := 0; i < v.Len(); i++ {
		if v.Index(i).Interface() == e {
			return true
		}
	}
	return false
}
