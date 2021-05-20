package controllers

import (
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"

	"github.com/CheerChen/konachan-app/internal/models"
)

// 本地相册
func Album(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	limit, page, err := GetPager(w, ps)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
	}
	query := GetQuery("tag", ps)

	posts := models.Posts{}
	err = posts.FetchAll(query, limit, page)

	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	log.Infof("fetch posts: %d", len(posts))

	if len(posts) != 0 {
		tfIdf, idf := getTfIdf()
		err = posts.Mark(tfIdf, idf)
		if err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
		}
		_ = posts.MarkExist()
	}

	cJson(w, posts, map[string]int{
		"total": len(posts),
	})
	return

}

// 检查数据一致
func Check(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	cJson(w, models.Check(), nil)
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
	err = models.DeleteFile(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	cJson(w, "OK", nil)
	return
}

// 输出缩略图
//func Preview(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
//	id, err := strconv.ParseInt(ps.ByName("id"), 10, 64)
//
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusNotAcceptable)
//		return
//	}
//
//	pic, err := kfile.GetFileById(id)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusNotFound)
//		return
//	}
//	if pic.Header == "" {
//		http.Error(w, "file format error", http.StatusNotFound)
//		return
//	}
//
//	file, err := os.Open(pic.Name)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusNotFound)
//		return
//	}
//	defer file.Close()
//	img, format, err := image.Decode(file)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//
//	var post models.Post
//	err = post.Find(id)
//
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusNotFound)
//		return
//	}
//
//	buf := new(bytes.Buffer)
//	dst := imaging.Resize(img, post.ActualPreviewWidth, post.ActualPreviewHeight, imaging.Lanczos)
//	switch format {
//	case "gif":
//		err = gif.Encode(buf, dst, nil)
//	case "png":
//		err = png.Encode(buf, dst)
//	case "jpeg":
//		fallthrough
//	default:
//		err = jpeg.Encode(buf, dst, nil)
//	}
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//	w.Header().Set("Content-type", pic.Header)
//	w.Header().Set("Cache-control", "max-age=315360000")
//	w.Write(buf.Bytes())
//}

// 输出全尺寸
func Sample(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.ParseInt(ps.ByName("id"), 10, 64)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}

	pic, err := models.GetFileById(id)
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
	w.Header().Set("Cache-control", "max-age=315360000")
	w.Write(byte)
	byte = nil
}
