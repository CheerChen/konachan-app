package controllers

import (
	"bytes"
	"image"
	"image/jpeg"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/julienschmidt/httprouter"

	"github.com/CheerChen/konachan-app/internal/kfile"
)

// 输出图片内容
func Preview(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.ParseInt(ps.ByName("id"), 10, 64)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}

	pics := kfile.LoadFiles()
	if len(pics) == 0 {
		http.Error(w, "no pics", http.StatusNotFound)
		return
	}
	var pic kfile.KFile
	for _, pic = range pics {
		if pic.Id == id {
			break
		}
	}
	var header string
	if strings.HasSuffix(pic.Name, ".png") {
		header = "image/png"
	} else if strings.HasSuffix(pic.Name, ".jpg") {
		header = "image/jpeg"
	} else if strings.HasSuffix(pic.Name, ".gif") {
		header = "image/gif"
	} else {
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

	resized := imaging.Resize(img, 100, 0, imaging.NearestNeighbor)
	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, resized, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-type", header)
	w.Write(buf.Bytes())
}

func Sample(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.ParseInt(ps.ByName("id"), 10, 64)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}

	pics := kfile.LoadFiles()
	if len(pics) == 0 {
		http.Error(w, "no pics", http.StatusNotFound)
		return
	}
	var pic kfile.KFile
	for _, pic = range pics {
		if pic.Id == id {
			break
		}
	}
	var header string
	if strings.HasSuffix(pic.Name, ".png") {
		header = "image/png"
	} else if strings.HasSuffix(pic.Name, ".jpg") {
		header = "image/jpeg"
	} else if strings.HasSuffix(pic.Name, ".gif") {
		header = "image/gif"
	} else {
		http.Error(w, "file format error", http.StatusNotFound)
		return
	}

	byte, err := ioutil.ReadFile(pic.Name)
	log.Println("read file:")
	log.Println(pic.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-type", header)
	w.Write(byte)

}
