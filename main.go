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
	"github.com/CheerChen/konachan-app/internal/logger"
	"github.com/CheerChen/konachan-app/internal/models"
	"github.com/CheerChen/konachan-app/internal/service/konachan"
)

var (
	c string
	//go:embed assets/*
	f   embed.FS
	log logger.Logger
)

func main() {
	log = logger.New()

	flag.StringVar(&c, "c", "", "specify configuration file")
	flag.Parse()

	conf, err := conf.NewLoader().LoadFile(c)
	if err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	controllers.Log(log)
	models.Log(log)
	models.OpenDbfile(conf.Dbfile)
	models.CheckPath("Wallpaper")

	myclient := client.New()
	if conf.Proxy.Enable {
		err := myclient.SetProxyUrl(conf.Proxy.Socket)
		if err != nil {
			log.Fatalf("Error load client proxy, %s", err)
		}
	}
	konachan.Set(myclient, conf.Download.Host, log)

	router := httprouter.New()

	// assets
	router.GET("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		bytes, err := f.ReadFile("assets/index.html")
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

	log.Infof("HTTP listening at: %s", conf.Addr)
	log.Fatalf("%s", http.ListenAndServe(conf.Addr, withGz))
}
