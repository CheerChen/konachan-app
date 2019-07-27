package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/CheerChen/konachan-app/internal/log"
)

const (
	PostLimit        = 100
	PostByTagUrl     = "https://konachan.net/post.json?tags=id:%d"
	PostByTagNameUrl = "https://konachan.net/post.json?tags=%s&limit=%d&page=%d"
	TagLimit         = 10000
	TagUrl           = "https://konachan.net/tag.json?order=count&limit=%d"
)

func GetRemotePosts(tags string, limit, page int) (ps Posts) {
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
	getBytes := cc.Get(url)
	if getBytes == nil {
		body, err := proxyGet(url)
		if err != nil {
			log.Errorf("Get remote url failed: %s", err)
			return
		}
		err = json.Unmarshal(body, &ps)
		if err != nil {
			log.Errorf("Get remote body json Unmarshal failed: %s", err)
			return
		}
		cc.Set(url, body)
	} else {
		err := json.Unmarshal(getBytes, &ps)
		if err != nil {
			log.Errorf("Cache json Unmarshal failed: %s", err)
		}
	}
	return
}

func GetRemotePost(postId int64) (target Post, err error) {
	url := fmt.Sprintf(PostByTagUrl, postId)
	body, err := proxyGet(url)
	if err != nil {
		log.Errorf("Get remote url failed: %s", err)
		return
	}
	var posts Posts
	err = json.Unmarshal(body, &posts)
	if err != nil {
		log.Errorf("Body json Unmarshal failed: %s", err)
		return
	}

	for _, post := range posts {
		if post.ID == postId {
			return post, nil
		}
	}
	return target, errors.New("post id not in page")
}

func GetGlobalTagCount() (total int, tagMap map[string]int) {
	tagMap = make(map[string]int)

	url := fmt.Sprintf(TagUrl, TagLimit)
	getBytes := cc.Get(url)
	if getBytes == nil {
		tags := Tags{}

		body, err := proxyGet(url)
		if err != nil {
			log.Errorf("Get remote url failed: %s", err)
			return
		}
		err = json.Unmarshal(body, &tags)
		if err != nil {
			log.Errorf("Get remote body json Unmarshal failed: %s", err)
			return
		}

		for _, tag := range tags {
			tagMap[tag.Name] = tag.Count
		}
		cc.Set(url, body)

	} else {
		err := json.Unmarshal(getBytes, &tagMap)
		if err != nil {
			log.Errorf("Cache json Unmarshal failed: %s", err)
		}
	}

	for _, count := range tagMap {
		total = total + count
	}

	return
}

func proxyGet(url string) (b []byte, err error) {
	client := &http.Client{}
	//dialer, _ := proxy.SOCKS5("tcp", "127.0.0.1:1080",
	//	nil,
	//	&net.Dialer{
	//		Timeout:   5 * time.Second,
	//		KeepAlive: 5 * time.Second,
	//	},
	//)
	//client.Transport = &http.Transport{
	//	Dial: dialer.Dial,
	//	//Proxy:           http.ProxyURL(proxy),
	//	//TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	//}
	log.Infof("Fetching url", url)
	resp, err := client.Get(url)
	if err != nil {
		return b, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}