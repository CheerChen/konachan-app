package kpost

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"time"

	"errors"

	"github.com/CheerChen/konachan-app/internal/memstore"
)

var (
	cache = memstore.NewMemoryStore(1 * time.Hour)
)

const (
	PostLimit = 100
	PostUrl   = "https://konachan.com/post.json?limit=%d&page=%d"
	PostKey   = "page_%d_limit_%d"
	//RangeUrl        = "https://konachan.com/post.json?tags=id:<=%d%20order:id_desc&limit=%d&page=%d"
	RangeKey        = "from_%d_page_%d_limit_%d"
	PostLatestIdKey = "latest_id"
	PostByTagUrl    = "https://konachan.com/post.json?tags=id:%d"

	PostUrlV2 = "https://konachan.com/post.json"
)

func GetPostsByPage(limit, page int) (posts []KPost) {
	if page <= 0 || limit <= 0 {
		return
	}
	if limit > PostLimit {
		limit = PostLimit
	}
	getPosts := cache.Get(fmt.Sprintf(PostKey, limit, page))
	if getPosts == nil {
		url := fmt.Sprintf(PostUrl, limit, page)
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		err = json.Unmarshal(body, &posts)
		if err != nil {
			fmt.Println(err)
		}
		cache.Set(fmt.Sprintf(PostKey, limit, page), posts)
	} else {
		posts = getPosts.([]KPost)
	}
	return
}

// 从大到小
func GetPostsByRange(from, limit, page int) (posts []KPost) {
	if from <= 0 {
		return
	}
	if page <= 0 || limit <= 0 {
		return
	}
	if limit > PostLimit {
		limit = PostLimit
	}
	rangeKey := fmt.Sprintf(RangeKey, from, limit, page)
	getPosts := cache.Get(rangeKey)
	if getPosts == nil {
		req, _ := http.NewRequest("GET", PostUrlV2, nil)
		req.Header.Add("Accept", "application/json")
		q := req.URL.Query()
		q.Add("tags", fmt.Sprintf("id:>=%d order:id", from))
		q.Add("limit", fmt.Sprintf("%d", limit))
		q.Add("page", fmt.Sprintf("%d", page))
		req.URL.RawQuery = q.Encode()

		fmt.Println("GetPostsByRange:req", req.URL.String())

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("GetPostsByRange:do", err)
			return
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		err = json.Unmarshal(body, &posts)
		if err != nil {
			fmt.Println("GetPostsByRange:json_unmarshal", err)
			fmt.Println(string(body))
			return
		}
		cache.Set(rangeKey, posts)
	} else {
		posts = getPosts.([]KPost)
	}
	return
}

func getLatestId() int {
	getId := cache.Get(PostLatestIdKey)
	if getId == nil {
		posts := GetPostsByPage(1, 1)

		cache.Set(PostLatestIdKey, posts[0].ID)
		return posts[0].ID
	}
	return getId.(int)
}

// 根据 post_id 获取
func GetPostByIdV2(postId int) (target KPost, err error) {
	url := fmt.Sprintf(PostByTagUrl, postId)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))

	var posts KPosts
	err = json.Unmarshal(body, &posts)
	if err != nil {
		fmt.Println(err)
	}

	for _, post := range posts {
		if post.ID == postId {
			return post, nil
		}
	}
	return target, errors.New("post id not in page")
}

func GetPostById(postId, limit int) (target KPost, err error) {

	if postId <= 0 {
		return target, errors.New("postId negative")
	}

	if limit <= 0 || limit > 100 {
		limit = PostLimit
	}

	latestId := getLatestId()

	// fix offset
	var page int
	page = (latestId-postId)/limit + 1
	posts := GetPostsByPage(limit, page)
	log.Println("try page", page)
	if len(posts) == 0 {
		return target, errors.New("no val return")
	}

	maxId := posts[0].ID
	minId := posts[limit-1].ID

	log.Println("maxId", maxId)

	density := float64(latestId-postId) / float64(latestId-maxId)
	log.Println("val density", density)

	var i int

	for !(maxId >= postId && postId >= minId) {
		if i > 5 {
			log.Println("loop over limit", i, page)
			break
		}
		i = i + 1
		if maxId < postId {
			diff := round(float64(postId-maxId) * density)
			if diff < limit {
				diff += limit
			}
			page = page - (diff / limit)
		}
		if postId < minId {
			diff := round(float64(maxId-postId) * density)
			if diff < limit {
				diff += limit
			}
			page = page + (diff / limit)
		}

		log.Println("try page", i, page)

		posts = GetPostsByPage(limit, page)
		if len(posts) == 0 {
			log.Println("try page failed", limit, page)
			break
		}
		maxId = posts[0].ID
		minId = posts[limit-1].ID

		log.Println("maxId", maxId)
		log.Println("minId", minId)
	}

	for _, post := range posts {
		if post.ID == postId {
			return post, nil
		}
	}

	return target, errors.New("post id not in page")
}

func round(x float64) int {
	return int(math.Floor(x + 0/5))
}

const (
	TagLimit = 5000
	TagUrl   = "https://konachan.com/tag.json?order=count&limit=%d"
	TagKey   = "tag_count"
)

func getGlobalTagCount() (int, map[string]int) {
	var globalTagCount map[string]int
	var globalTotal int

	getMap := cache.Get(TagKey)
	if getMap == nil {
		globalTagCount = make(map[string]int)
		remoteTags := getRemoteTag(TagLimit)

		for _, tag := range remoteTags {
			globalTagCount[tag.Name] = tag.Count
		}
		cache.Set(TagKey, globalTagCount)

	} else {
		globalTagCount = getMap.(map[string]int)
	}

	for _, count := range globalTagCount {
		globalTotal = globalTotal + count
	}

	return globalTotal, globalTagCount
}

func getRemoteTag(limit int) (tags Tags) {
	if limit <= 0 {
		return
	}

	url := fmt.Sprintf(TagUrl, limit)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &tags)
	if err != nil {
		fmt.Println(err)
	}
	return
}
