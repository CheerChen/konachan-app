package controllers

import (
	"bytes"
	"image"
	"image/jpeg"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/disintegration/imaging"
	"github.com/julienschmidt/httprouter"

	"github.com/CheerChen/konachan-app/internal/kfile"
	"github.com/CheerChen/konachan-app/internal/models"
)

// 本地相册
func Album(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	limit, page, err := GetPager(w, ps)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
	}

	tag := ps.ByName("tag")[1:]
	posts, err := models.GetPosts(tag, limit, page)

	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	if len(posts) == 0 {
		http.Error(w, "no posts", http.StatusNotFound)
		return
	}

	tfIdf := models.GetTfIdf()
	marked := posts.MarkAndReduce(0.0, tfIdf)

	cJson(w, marked, map[string]int{
		"total": len(marked),
	})
	return

}

// 输出图片分布
func Dis(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	limit, err := strconv.Atoi(ps.ByName("limit"))

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}

	idMap, err := models.GetIdMap()

	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	disMap := make(map[int64]int)
	for id := range idMap {
		dis := id / int64(limit)
		if _, ok := disMap[dis]; !ok {
			disMap[dis] = 1
		} else {
			disMap[dis] += 1
		}
	}

	cJson(w, disMap, nil)
	return
}

// 检查数据一致
func Check(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	cJson(w, kfile.Check(), nil)
}

// 从本地删除
func Delete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.ParseInt(ps.ByName("id"), 10, 64)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}
	var post models.Post
	err = post.Find(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	err = post.Delete()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	err = kfile.Delete(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	_, _ = w.Write([]byte("<html><body><script>window.location.href=\"about:blank\";window.close();</script></body></html>"))

	return
}

// 输出缩略图
func Preview(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.ParseInt(ps.ByName("id"), 10, 64)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}

	pic, err := kfile.GetFileById(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if pic.Header == "" {
		http.Error(w, "file format error", http.StatusNotFound)
		return
	}

	file, err := os.Open(pic.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var post models.Post
	err = post.Find(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	dst := imaging.Resize(img, post.ActualPreviewWidth, post.ActualPreviewHeight, imaging.Lanczos)
	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, dst, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-type", pic.Header)
	w.Write(buf.Bytes())
}

// 输出全尺寸
func Sample(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.ParseInt(ps.ByName("id"), 10, 64)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}

	pic, err := kfile.GetFileById(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if pic.Header == "" {
		http.Error(w, "file format error", http.StatusNotFound)
		return
	}

	byte, err := ioutil.ReadFile(pic.Name)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-type", pic.Header)
	w.Write(byte)

}
