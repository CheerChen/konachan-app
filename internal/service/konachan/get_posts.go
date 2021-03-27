package konachan

import (
	"errors"
	"fmt"
	"github.com/CheerChen/konachan-app/internal/log"
	"github.com/CheerChen/konachan-app/internal/models"
	"time"
)

func GetPosts(query string, pageSize, page int) *models.Posts {
	merge := new(models.Posts)
	endPage := page + (pageSize / postLimit)
	timeTag := time.Now().Format(hourLayout)

	for page < endPage {
		posts := new(models.Posts)

		u := fmt.Sprintf(postUrlFmt, hostname, postLimit, page, query)

		v, ok := lru.Get(timeTag + u)
		if ok {
			posts = v.(*models.Posts)
		} else {
			log.Infof("NewRequest: %s", u)
			err := myclient.GetJSON(u, posts)
			if err != nil {
				log.Errorf("[GetPosts] GetJSON err:", err)
				continue
			}
			lru.Put(timeTag+u, posts)
		}

		if len(*posts) != 0 {
			if (*posts)[0].ID > lastid {
				lastid = (*posts)[0].ID
			}
		}

		*merge = append(*merge, *posts...)

		page += 1
	}

	return merge
}

func GetPostLastId() (id int64, err error) {
	return lastid, nil
}

func GetPostByID(id int64) (*models.Post, error) {
	posts := &models.Posts{}
	u := fmt.Sprintf(postUrlFmt, hostname, 1, 1, fmt.Sprintf("id:%d", id))
	log.Infof("NewRequest: %s", u)
	err := myclient.GetJSON(u, posts)

	if err != nil {
		log.Errorf("GetPostByID err:", err)
		return nil, err
	}

	if len(*posts) > 0 {

		return &(*posts)[0], nil
	}

	return nil, errors.New("record not found")
}
