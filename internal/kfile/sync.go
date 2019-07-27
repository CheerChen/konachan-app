package kfile

import (
	"fmt"

	"github.com/CheerChen/konachan-app/internal/log"
	"github.com/CheerChen/konachan-app/internal/models"
)

func Sync(path string) {
	pics := LoadFiles(path)
	if len(pics) == 0 {
		log.Warnf("Pic files Not Found!")
		return
	}

	for _, pic := range pics {
		var post models.Post
		err := post.Find(pic.Id)
		if err != nil {
			post, err = models.GetRemotePost(pic.Id)
			if err != nil {
				log.Infof("fetch id error", err.Error())
				continue
			} else {
				err = post.Save()
				if err != nil {
					log.Infof("sync to db error", err.Error())
				} else {
					log.Infof("sync to db", post.ID)
				}
			}
			//posts := kpost.GetPostsByTags(strings.Replace(pic.Tags, " ", "+", -1))
			//if len(posts) == 0 {
			//	log.Println("fetch empty posts")
			//	continue
			//} else {
			//	for _, post := range posts {
			//		if post.ID == pic.Id {
			//			err = post.Sync2DB()
			//			if err != nil {
			//				log.Println("sync to db error", err.Error())
			//			} else {
			//				fmt.Println("sync to db", post.ID)
			//			}
			//		} else {
			//			log.Println("fetch wrong post", post.ID)
			//		}
			//	}
			//}
		} else {
			//log.Println("find post ", post.ID)
		}

	}

	fmt.Fprint(w, "done\n")

	return
}
