package kpost

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"errors"

	"github.com/CheerChen/konachan-app/internal/memstore"
)

var (
	cache = memstore.NewMemoryStore(1 * time.Hour)
)

const (
	PostLimit        = 100
	PostByTagUrl     = "https://konachan.net/post.json?tags=id:%d"
	PostByTagNameUrl = "https://konachan.net/post.json?tags=%s&limit=%d&page=%d"
	TagLimit         = 10000
	TagUrl           = "https://konachan.net/tag.json?order=count&limit=%d"
	TagKey           = "tag_count"
)

func GetPosts(tags string, limit, page int) (posts KPosts) {
	if page <= 0 || limit <= 0 {
		return
	}
	if limit > PostLimit {
		limit = PostLimit
	}
	req, _ := http.NewRequest("GET", PostByTagNameUrl, nil)
	req.Header.Add("Accept", "application/json")
	q := req.URL.Query()
	q.Add("tags", fmt.Sprintf("%s", tags))
	q.Add("limit", fmt.Sprintf("%d", limit))
	q.Add("page", fmt.Sprintf("%d", page))
	req.URL.RawQuery = q.Encode()

	url := req.URL.String()
	//url := fmt.Sprintf(PostUrl, limit, page)
	//if len(tags) != 0 {
	//}
	getPosts := cache.Get(url)
	if getPosts == nil {
		body, err := proxyGet(url)
		if err != nil {
			fmt.Println(err)
			return
		}
		err = json.Unmarshal(body, &posts)
		if err != nil {
			fmt.Println(err)
		}
		cache.Set(url, posts)
	} else {
		posts = getPosts.(KPosts)
	}
	return
}

// 根据 post_id 获取
func GetPostByIdV2(postId int) (target KPost, err error) {
	url := fmt.Sprintf(PostByTagUrl, postId)
	body, err := proxyGet(url)
	if err != nil {
		fmt.Println(err)
		return
	}
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
	body, err := proxyGet(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = json.Unmarshal(body, &tags)
	if err != nil {
		fmt.Println(err)
	}
	return
}
