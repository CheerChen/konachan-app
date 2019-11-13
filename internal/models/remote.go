package models

import (
	"encoding/json"
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
	if limit > PostLimit || limit < 1 {
		limit = PostLimit
	}
	if page < 1 {
		page = 1
	}
	req, _ := http.NewRequest("GET", PostByTagNameUrl, nil)
	req.Header.Add("Accept", "application/json")
	q := req.URL.Query()
	q.Add("tags", fmt.Sprintf("%s", tags))
	q.Add("limit", fmt.Sprintf("%d", limit))
	q.Add("page", fmt.Sprintf("%d", page))
	req.URL.RawQuery = q.Encode()

	url := req.URL.String()

	body, err := proxyGet(url)
	if err != nil {
		log.Errorf("http get: %s", err)
		return
	}

	err = json.Unmarshal(body, &ps)
	if err != nil {
		log.Errorf("json Unmarshal: %s", err)
	}

	return
}

func GetRemotePost(postId int64) (target Post, err error) {
	url := fmt.Sprintf(PostByTagUrl, postId)
	body, err := proxyGet(url)
	if err != nil {
		log.Errorf("http get: %s", err)
		return
	}
	var posts Posts
	err = json.Unmarshal(body, &posts)
	if err != nil {
		log.Errorf("json Unmarshal: %s", err)
		return
	}

	for _, post := range posts {
		if post.ID == postId {
			return post, nil
		}
	}
	return target, ErrRecordNotFound
}

func GetRemoteTags() (ts Tags) {
	url := fmt.Sprintf(TagUrl, TagLimit)
	getBytes := cc.Get(url)
	if getBytes == nil {
		body, err := proxyGet(url)
		if err != nil {
			log.Errorf("http get: %s", err)
			return
		}
		getBytes = body
		cc.Set(url, body)
	}

	err := json.Unmarshal(getBytes, &ts)
	if err != nil {
		log.Errorf("json Unmarshal: %s", err)
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
	log.Infof("fetch url %s", url)
	resp, err := client.Get(url)
	if err != nil {
		return b, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
