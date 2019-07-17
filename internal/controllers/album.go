package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"

	"github.com/CheerChen/konachan-app/internal/kpost"
)

// 输出图片内容
func Album(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	limit, err := strconv.Atoi(ps.ByName("limit"))

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}

	page, err := strconv.Atoi(ps.ByName("page"))

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}

	if limit <= 0 || page <= 0 {
		http.Error(w, http.StatusText(http.StatusNotAcceptable), http.StatusNotAcceptable)
	}

	tfIdf := kpost.GetTfIdf()

	posts, err := kpost.SelectPostByPage(limit, page)

	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	log.Println("fetch posts:")
	log.Println(len(posts))

	if len(posts) == 0 {
		http.Error(w, "no posts", http.StatusNotFound)
		return
	}

	marked := posts.MarkNotReduce(tfIdf)
	//reduced.MarkDownloaded()

	cJson(w, marked, map[string]int{
		"total":   len(posts),
		"reduced": len(marked),
	})
	return

}

func Prefix(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	tfIdf := kpost.GetTfIdf()

	posts, err := kpost.SelectPostByPrefix(ps.ByName("p"))

	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	log.Println("fetch posts:")
	log.Println(len(posts))

	if len(posts) == 0 {
		http.Error(w, "no posts", http.StatusNotFound)
		return
	}

	marked := posts.MarkNotReduce(tfIdf)
	//reduced.MarkDownloaded()

	cJson(w, marked, map[string]int{
		"total":   len(posts),
		"reduced": len(marked),
	})
	return

}

// 输出图片内容
func Search(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	tfIdf := kpost.GetTfIdf()

	posts, err := kpost.SelectPostByTag(ps.ByName("tag"))

	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	log.Println("fetch posts:")
	log.Println(len(posts))

	if len(posts) == 0 {
		http.Error(w, "no posts", http.StatusNotFound)
		return
	}

	marked := posts.MarkNotReduce(tfIdf)
	//reduced.MarkDownloaded()

	cJson(w, marked, map[string]int{
		"total":   len(posts),
		"reduced": len(marked),
	})
	return

}

// 输出图片分布
func Dis(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	limit, err := strconv.Atoi(ps.ByName("limit"))

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}

	idMap, err := kpost.SelectAllIds2Map()

	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	disMap := make(map[int]int)
	for id := range idMap {
		dis := id / limit
		if _, ok := disMap[dis]; !ok {
			disMap[dis] = 1
		} else {
			disMap[dis] += 1
		}
	}

	cJson(w, disMap, nil)
	return
}
