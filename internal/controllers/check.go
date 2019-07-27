package controllers

import (
	"github.com/CheerChen/konachan-app/internal/humanize"
	"github.com/CheerChen/konachan-app/internal/kfile"
	"github.com/CheerChen/konachan-app/internal/models"

	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
)

func Check(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	pics := kfile.LoadFiles()
	if len(pics) == 0 {
		http.Error(w, "no pics", http.StatusNotFound)
		return
	}

	checkList := make(map[int64]string)
	for _, pic := range pics {
		var post models.Post
		err := post.Find(pic.Id)
		if err != nil {
			log.Println("post id not in db", err.Error(), pic.Id)
			continue
		}

		width, height := getImageDimension(pic.Name)

		if width != post.Width || height != post.Height {
			//if width < 1920 && width < post.Width {
			checkList[post.ID] = fmt.Sprintf("pic size：%s \n", humanize.Bytes(uint64(pic.Size)))
			checkList[post.ID] += fmt.Sprintf("post size：%s \n", humanize.Bytes(uint64(post.FileSize)))
			checkList[post.ID] += fmt.Sprintf("pic：%d*%d \n", width, height)
			checkList[post.ID] += fmt.Sprintf("post：%d*%d \n", post.Width, post.Height)
			//var file kfile.KFile
			//file.Id = post.ID
			//file.Tags = post.Tags
			//file.Ext = post.GetFileExt()
			//file.SlimTags()
			//
			//url := kfile.DownloadHelper(post.FileURL)
			//go kfile.DownloadFile(file.BuildName(), url)
			//}
		}

	}

	cJson(w, checkList, nil)

}

func getImageDimension(imagePath string) (int, int) {
	file, err := os.Open(imagePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
	defer file.Close()
	img, _, err := image.DecodeConfig(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", err)
	}
	return img.Width, img.Height
}
