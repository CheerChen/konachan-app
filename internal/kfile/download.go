package kfile

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func DownloadHelper(fileUrl string) string {
	return "http://konachan.wjcodes.com/konachan_download.php?url=" + fileUrl
}

func (pic KFile) BuildName() string {
	pic.Tags = strings.Replace(pic.Tags, "/", "_", -1)
	pic.Tags = strings.Replace(pic.Tags, ":", "_", -1)
	return fmt.Sprintf("Konachan.com - %d %s%s", pic.Id, pic.Tags, pic.Ext)
}

func DownloadFile(name string, u string) error {

	filePath := FilePath + string(os.PathSeparator) + name
	// Create the file
	out, err := os.Create(filePath)
	if err != nil {
		log.Println(err)
		return err
	}
	defer out.Close()

	// Get the data

	//proxyUrl, err := url.Parse("http://127.0.0.1:1080")
	//myClient := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}
	myClient := &http.Client{}

	resp, err := myClient.Get(u)
	if err != nil {
		log.Println(err)
		return err
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

const FileNameLengthLimit = 200

func (pic *KFile) SlimTags() {

	tags := strings.Split(pic.Tags, " ")

	for len(pic.Tags) >= FileNameLengthLimit {
		tags = tags[:len(tags)-1]
		pic.Tags = strings.Join(tags, " ")
	}
}
