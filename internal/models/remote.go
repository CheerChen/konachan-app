package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/cavaliercoder/grab"

	"github.com/CheerChen/konachan-app/internal/grabber"
	"github.com/CheerChen/konachan-app/internal/log"
	"github.com/CheerChen/konachan-app/internal/service/konachan"
)

// grab 版本
func BatchFetchPosts(query string, pagesize, page int) *Posts {
	posts := make(Posts, 0)
	endPage := page + (pagesize / 100)
	reqs := make([]*grab.Request, 0)
	for page < endPage {
		dst := fmt.Sprintf("temp/%s_%d.json", query, page)
		urlStr := fmt.Sprintf("https://konachan.com/post.json?limit=100&page=%d&tags=%s", page, query)
		req, _ := grab.NewRequest(dst, urlStr)
		reqs = append(reqs, req)
		page += 1
	}

	realReqs := make([]*grab.Request, 0)
	for _, req := range reqs {
		info, err := os.Stat(req.Filename)
		if os.IsNotExist(err) {
			realReqs = append(realReqs, req)
		} else if !os.IsNotExist(err) {
			if (time.Now().Sub(info.ModTime())) > 1*time.Hour {
				_ = os.Remove(req.Filename)
				realReqs = append(realReqs, req)
			}
		} else {
			log.Warnf("os.Stat err:", err)
			continue
		}
	}

	g := grabber.NewDownloadClient()
	g.SetProxy(proxyClient)
	g.BatchDownload(realReqs)

	idMap := make(map[int64]int)
	for _, req := range reqs {
		temp := make(Posts, 0)
		content, err := ioutil.ReadFile(req.Filename)
		if err != nil {
			log.Warnf("ReadDir err:", err)
			continue
		}
		err = json.Unmarshal(content, &temp)
		if err != nil {
			log.Warnf("Unmarshal err:", err)
			continue
		}
		for _, t := range temp {
			if _, ok := idMap[t.ID]; ok {
				continue
			}
			idMap[t.ID] = 1
			posts = append(posts, t)
		}
	}

	return &posts
}

func FetchPostByID(id int64) (post *Post, err error) {
	kPosts := make([]konachan.Post, 0)
	params := &konachan.PostListParams{
		Limit: 1,
		Page:  1,
		Tags:  fmt.Sprintf("id:%d", id),
	}
	kClient := konachan.NewClient(proxyClient)
	kPosts, _, err = kClient.Posts.List(params)
	if err != nil {
		log.Errorf("kClient.Posts.List: %s", err)
		return
	}
	if len(kPosts) > 0 {
		return &Post{Post: &kPosts[0]}, nil
	}

	return nil, ErrRecordNotFound
}

func GetRemoteTags() (tags []*Tag) {
	g := grabber.NewDownloadClient()
	g.SetProxy(proxyClient)
	urlStr := "https://konachan.com/tag.json?limit=10000&order=count"
	req, _ := grab.NewRequest("temp/tag.json", urlStr)

	info, err := os.Stat(req.Filename)
	if os.IsNotExist(err) {
		g.BatchDownload([]*grab.Request{req})
	} else if !os.IsNotExist(err) {
		if (time.Now().Sub(info.ModTime())) > 1*time.Hour {
			_ = os.Remove(req.Filename)
			g.BatchDownload([]*grab.Request{req})
		}
	} else {
		log.Warnf("os.Stat err:", err)
		return
	}
	var b []byte
	b, err = ioutil.ReadFile(req.Filename)
	if err != nil {
		log.Warnf("ReadDir err:", err)
		return
	}
	err = json.Unmarshal(b, &tags)
	if err != nil {
		log.Warnf("Unmarshal err:", err)
		return
	}
	return
}

func GetLastId() (id int64, err error) {
	dst := fmt.Sprintf("%s_%d.json", "", 1)
	_, err = os.Stat(dst)
	if !os.IsNotExist(err) {
		temp := make(Posts, 0)
		var b []byte
		b, err = ioutil.ReadFile(dst)
		if err != nil {
			return
		}
		err = json.Unmarshal(b, &temp)
		if err != nil {
			return
		}
		return temp[0].ID, nil
	}
	return
}
