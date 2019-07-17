package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/CheerChen/konachan-app/internal/kfile"
	"github.com/CheerChen/konachan-app/internal/kpost"
)

// 检查文件分辨率和数据库的同步

func Sync(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	pics := kfile.LoadFiles()
	if len(pics) == 0 {
		http.Error(w, "no pics", http.StatusNotFound)
		return
	}

	for _, pic := range pics {
		var post kpost.KPost
		err := post.Find(pic.Id)
		if err != nil {
			log.Println("getting post from remote", pic.Id)
			post, err = kpost.GetPostByIdV2(pic.Id)
			if err != nil {
				log.Println("fetch id error", err.Error())
				continue
			} else {
				err = post.Sync2DB()
				if err != nil {
					log.Println("sync to db error", err.Error())
				} else {
					fmt.Println("sync to db", post.ID)
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

// 检查文件和数据库的同步

func Sync2(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	pics := kfile.LoadFiles()
	if len(pics) == 0 {
		http.Error(w, "no pics", http.StatusNotFound)
		return
	}

	idMap, err := kpost.SelectAllIds2Map()
	if err != nil {
		http.Error(w, "no id map", http.StatusNotFound)
		return
	}
	for _, pic := range pics {
		if _, ok := idMap[pic.Id]; ok {
			idMap[pic.Id] = false
		}
	}

	cJson(w, idMap, nil)
}
