package konachan

import (
	"errors"
	"fmt"
	"net/url"
	"time"
)

func GetPosts(query string, pageSize, page int) (merge []Post) {
	endPage := page + (pageSize / postLimit)
	timeTag := time.Now().Format(hourLayout)

	for page < endPage {
		var posts []Post

		u := fmt.Sprintf(postUrlFmt, hostname, postLimit, page, url.QueryEscape(query))

		v, ok := lru.Get(timeTag + u)
		if ok {
			posts = v.([]Post)
		} else {
			log.Infof("req: %s", u)
			err := myclient.GetJSON(u, &posts)
			if err != nil {
				log.Errorf("get json err:", err)
				break
			}
			lru.Put(timeTag+u, posts)
		}

		if len(posts) != 0 {
			if posts[0].ID > lastid {
				lastid = posts[0].ID
			}
		}

		merge = append(merge, posts...)
		page += 1
	}

	return merge
}

func GetPostLastId() (id int64, err error) {
	return lastid, nil
}

func GetPostByID(id int64) (*Post, error) {
	var posts []Post
	u := fmt.Sprintf(postUrlFmt, hostname, 1, 1, fmt.Sprintf("id:%d", id))
	log.Infof("req: %s", u)
	err := myclient.GetJSON(u, &posts)

	if err != nil {
		log.Errorf("get json err:", err)
		return nil, err
	}

	if len(posts) > 0 {
		return &(posts[0]), nil
	}

	return nil, errors.New("record not found")
}
