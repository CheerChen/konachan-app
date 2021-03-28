package models

import (
	"fmt"
	"image"
	"os"
)

func Check() (diff map[int64]string) {
	pics := LoadFiles(WallpaperPath)
	if len(pics) == 0 {
		log.Warnf("Wallpaper path empty!")
		return
	}

	diff = make(map[int64]string)
	for _, pic := range pics {
		var post Post
		err := post.Find(pic.Id)
		if err != nil {
			log.Warnf("Post ID(%d) not in db", pic.Id)
			continue
		}

		width, height, err := getImageDimension(pic.Name)
		if err != nil {
			log.Errorf("getImageDimension err: %v", err)
			continue
		}

		if width != post.Width || height != post.Height {
			//diff[post.ID] = fmt.Sprintf("pic size：%d \n", pic.Size/1024)
			//diff[post.ID] += fmt.Sprintf("post size：%d \n", post.FileSize/1024)
			diff[post.ID] = fmt.Sprintf("rate：%f %% \n", float64(pic.Size)/float64(post.FileSize)*100)
			diff[post.ID] += fmt.Sprintf("pic：%d*%d \n", width, height)
			diff[post.ID] += fmt.Sprintf("post：%d*%d \n", post.Width, post.Height)
		}

	}
	return
}

func getImageDimension(imagePath string) (int, int, error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return 0, 0, err

	}
	defer file.Close()
	img, _, err := image.DecodeConfig(file)
	if err != nil {
		return 0, 0, err
	}
	return img.Width, img.Height, nil
}
