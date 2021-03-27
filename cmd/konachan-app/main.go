package main

import (
	"embed"
	"flag"
	"net/http"

	"github.com/NYTimes/gziphandler"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"

	"github.com/CheerChen/konachan-app/internal/client"
	"github.com/CheerChen/konachan-app/internal/conf"
	"github.com/CheerChen/konachan-app/internal/controllers"
	"github.com/CheerChen/konachan-app/internal/kfile"
	"github.com/CheerChen/konachan-app/internal/log"
	"github.com/CheerChen/konachan-app/internal/models"
	"github.com/CheerChen/konachan-app/internal/service/konachan"
)

var (
	configFile string
	//go:embed static/*
	f embed.FS
)

func init() {
	flag.StringVar(&configFile, "c", "", "specify configuration file")
	flag.Parse()

	conf.OpenCfgfile(configFile)
	models.OpenDbfile(conf.Cfg.Dbfile)

	kfile.CheckPath(conf.Cfg.Download.Path)

	myclient := client.New()
	if conf.Cfg.Proxy.Enable {
		err := myclient.SetProxyUrl(conf.Cfg.Proxy.Socket)
		if err != nil {
			log.Fatalf("Error load client proxy, %s", err)
		}
	}
	konachan.SetHost(conf.Cfg.Download.Host)
	konachan.SetClient(myclient)
}

func main() {
	router := httprouter.New()

	// static
	router.GET("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		bytes, err := f.ReadFile("static/index.html")
		if err != nil {
			http.NotFound(w, r)
			return
		}
		w.Write(bytes)
	})
	router.ServeFiles("/web/*filepath", http.FS(f))

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
