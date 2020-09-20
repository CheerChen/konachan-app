package main

import (
	"flag"
	"net/http"

	"github.com/CheerChen/konachan-app/internal/asset"
	"github.com/CheerChen/konachan-app/internal/conf"
	"github.com/CheerChen/konachan-app/internal/controllers"
	"github.com/CheerChen/konachan-app/internal/grabber"
	"github.com/CheerChen/konachan-app/internal/kfile"
	"github.com/CheerChen/konachan-app/internal/log"
	"github.com/CheerChen/konachan-app/internal/models"

	"github.com/NYTimes/gziphandler"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
)

var (
	configFile string
)

func init() {
	flag.StringVar(&configFile, "f", "config", "specify configuration file")
	flag.Parse()
}

func main() {
	conf.OpenCfgfile(configFile)
	models.OpenDbfile(conf.Cfg.Dbfile)

	kfile.CheckPath(conf.Cfg.Download.Path)
	grabber.SetHost(conf.Cfg.Download.Host)
	grabber.SetProxy(conf.Cfg.Proxy.Enable, conf.Cfg.Proxy.Socket)

	router := httprouter.New()

	// static
	router.GET("/", Index)
	router.ServeFiles("/web/*filepath", asset.AssetFS())

	// api
	router.GET("/remote/:limit/:page/*tag", controllers.Remote)
	router.GET("/post/:id", controllers.GetByIdV2)
	router.GET("/tag/tf_idf", controllers.GetTfIdf)
	router.GET("/download/:id", controllers.Download)
	router.GET("/album/:limit/:page/*tag", controllers.Album)
	router.GET("/check", controllers.Check)
	router.GET("/delete/:id", controllers.Delete)
	router.GET("/sample/:id", controllers.Sample)

	handler := cors.Default().Handler(router)
	withGz := gziphandler.GzipHandler(handler)

	log.Infof("HTTP listening at: %s", conf.Cfg.Addr)
	log.Fatalf("%s", http.ListenAndServe(conf.Cfg.Addr, withGz))
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
