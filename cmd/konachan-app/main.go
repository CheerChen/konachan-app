package main

import (
	"context"
	"crypto/tls"
	"flag"
	slog "log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"runtime"

	"github.com/CheerChen/konachan-app/internal/asset"
	"github.com/CheerChen/konachan-app/internal/controllers"
	"github.com/CheerChen/konachan-app/internal/kfile"
	"github.com/CheerChen/konachan-app/internal/log"
	"github.com/CheerChen/konachan-app/internal/models"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	"golang.org/x/net/proxy"
)

var (
	scanPath string
	proxyUrl string
)

func init() {
	flag.StringVar(&scanPath, "p", "wallpaper", "wallpaper path, input an absolute path")
	flag.StringVar(&proxyUrl, "proxy", "", "set proxy url")
	flag.Parse()
}

func main() {
	if err := ensureDir(scanPath); err != nil {
		log.Fatal(err)
	}

	proxyClient := &http.Client{}
	if proxyUrl != "" {
		url, err := url.Parse(proxyUrl)
		if err != nil {
			log.Fatal(err)
		}
		dialer, err := proxy.FromURL(url, proxy.Direct)
		if err != nil {
			log.Fatal(err)
		}
		proxyClient.Transport = &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (conn net.Conn, e error) {
				c, e := dialer.Dial(network, addr)
				return c, e
			},
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}

	}

	models.SetClient(proxyClient)
	kfile.Sync(scanPath)

	router := httprouter.New()

	// static
	router.GET("/", Index)
	router.ServeFiles("/web/*filepath", asset.AssetFS())

	//
	router.GET("/remote/:limit/:page/*tag", controllers.Remote)
	router.GET("/post/:id", controllers.GetByIdV2)
	router.GET("/tag/tf_idf", controllers.GetTfIdf)
	router.GET("/download/:id", controllers.Download)

	router.GET("/album/:limit/:page/*tag", controllers.Album)
	router.GET("/check", controllers.Check)
	router.GET("/delete/:id", controllers.Delete)
	router.GET("/preview/:id", controllers.Preview)
	router.GET("/sample/:id", controllers.Sample)
	router.GET("/dist/:limit", controllers.Dis)

	handler := cors.Default().Handler(router)

	go Open("http://127.0.0.1:7079/")
	slog.Fatal(http.ListenAndServe(":7079", handler))
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	bytes, err := asset.Asset("web/static/index.html")
	if err != nil {
		log.Errorf("", err)
		http.NotFound(w, r)
		return
	}
	w.Write(bytes)
}

func Open(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}

func ensureDir(dirName string) error {
	err := os.Mkdir(dirName, os.ModeDir)
	if err == nil || os.IsExist(err) {
		return nil
	} else {
		return err
	}
}
