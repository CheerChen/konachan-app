package models

import (
	"fmt"
	"github.com/CheerChen/konachan-app/internal/log"
	"github.com/CheerChen/konachan-app/internal/service/konachan"
)

func Fetch(params *konachan.PostListParams) (posts Posts) {
	kClient := konachan.NewClient(proxyClient)
	key, _ := kClient.Posts.ListUrlEncode(params)
	val, found := mem.Get(key)

	if !found {
		kPosts, _, err := kClient.Posts.List(params)
		if err != nil {
			log.Errorf("kClient.Posts.List: %s", err)
			return
		}
		for k := range kPosts {
			posts = append(posts, Post{
				Post: &kPosts[k],
			})
		}
		mem.SetDefault(key, posts)
	} else {
		posts = val.(Posts)
	}

	return
}

func FetchId(id int64) (post Post, err error) {
	params := &konachan.PostListParams{
		Limit: 1,
		Page:  1,
		Tags:  fmt.Sprintf("id:%d", id),
	}
	posts := Fetch(params)
	if len(posts) > 0 {
		return posts[0], nil
	}

	return post, ErrRecordNotFound
}

func GetRemoteTags() (tags Tags) {
	params := &konachan.TagListParams{
		Limit: 10000,
		Order: "count",
	}
	kClient := konachan.NewClient(proxyClient)
	key, _ := kClient.Tags.ListUrlEncode(params)
	val, found := mem.Get(key)

	if !found {
		kTags, _, err := kClient.Tags.List(params)
		if err != nil {
			log.Errorf("kClient.Tags.List: %s", err)
			return
		}

		for k := range kTags {
			tags = append(tags, Tag{
				Tag: &kTags[k],
			})
		}
		mem.SetDefault(key, tags)
	} else {
		tags = val.(Tags)
	}

	return
}

func FetchLastId() (id int64) {
	params := &konachan.PostListParams{
		Limit: 1,
		Page:  1,
		Tags:  "",
	}

	posts := Fetch(params)
	if len(posts) > 0 {
		return posts[0].ID
	}
	return
}
