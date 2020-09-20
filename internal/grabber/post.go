package grabber

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/CheerChen/konachan-app/internal/log"
	"github.com/CheerChen/konachan-app/internal/models"
	"github.com/CheerChen/konachan-app/internal/service/konachan"

	"github.com/cavaliercoder/grab"
)

const (
	postDstFmt = "%s/%s_%d.json"
	postLimit  = 100
	postUrlFmt = "%s/post.json?limit=%d&page=%d&tags=%s"
	tagDstFmt  = "%s/tag.json"
	tagLimit   = 10000
	tagUrlFmt  = "%s/tag.json?limit=%d&order=count"

	cacheTime = 1 * time.Hour
)

func GetPosts(query string, pageSize, page int) *models.Posts {
	posts := make(models.Posts, 0)
	endPage := page + (pageSize / postLimit)
	reqs := make([]*grab.Request, 0)
	for page < endPage {
		dst := fmt.Sprintf(postDstFmt, os.TempDir(), query, page)
		urlStr := fmt.Sprintf(postUrlFmt, hostname, postLimit, page, query)
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
			if (time.Now().Sub(info.ModTime())) > cacheTime {
				_ = os.Remove(req.Filename)
				realReqs = append(realReqs, req)
			}
		} else {
			log.Warnf("os.Stat err:", err)
			continue
		}
	}

	g := NewDownloadClient()
	g.SetProxy(proxyClient)
	g.BatchDownload(realReqs)

	idMap := make(map[int64]int)
	for _, req := range reqs {
		temp := make(models.Posts, 0)
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

func GetPostLastId() (id int64, err error) {
	dst := fmt.Sprintf(postDstFmt, os.TempDir(), "", 1)
	_, err = os.Stat(dst)
	if !os.IsNotExist(err) {
		temp := make(models.Posts, 0)
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

func GetPostByID(id int64) (*models.Post, error) {
	//kPosts := make([]konachan.Post, 0)
	post := new(models.Post)
	params := &konachan.PostListParams{
		Limit: 1,
		Page:  1,
		Tags:  fmt.Sprintf("id:%d", id),
	}
	kClient := konachan.NewClient(proxyClient)
	kPosts, _, err := kClient.Posts.List(params)
	if err != nil {
		log.Errorf("kClient.Posts.List: %s", err)
		return nil, err
	}
	if len(kPosts) > 0 {
		post.Post = &kPosts[0]
		return post, nil
	}

	return nil, ErrRecordNotFound
}
