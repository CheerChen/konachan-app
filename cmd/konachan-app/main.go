package main

import (
	"flag"
	slog "log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/CheerChen/konachan-app/internal/controllers"
	"github.com/CheerChen/konachan-app/internal/log"
	"github.com/CheerChen/konachan-app/internal/models"
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	http.ServeFile(w, r, "static/index.html")
}

var (
	level = flag.Int("l", -1, "log level, -1:debug, 0:info, 1:warn, 2:error")
)

func main() {

	lcf := zap.NewDevelopmentConfig()
	lcf.Level.SetLevel(zapcore.Level(*level))
	lcf.Development = false
	lcf.Sampling = nil
	lcf.DisableStacktrace = true

	logger, err := lcf.Build(zap.AddCallerSkip(1))
	if err != nil {
		slog.Fatalln("logger err:", err.Error())
	}
	log.SetLogger(logger.Sugar())

	models.Init()

	router := httprouter.New()

	// static
	router.GET("/", Index)
	router.ServeFiles("/static/*filepath", http.Dir("static"))

	//
	router.GET("/post/:id", controllers.GetByIdV2)
	router.GET("/tag/tf_idf", controllers.GetTfIdf)
	//router.GET("/tag/update", Update)
	router.GET("/hot/:limit/:page", controllers.Popular)
	//router.GET("/hot_from/:limit/:page/:from", controllers.PopularByRange)
	router.GET("/download/:id", controllers.Download)
	router.GET("/check", controllers.Check)

	router.GET("/delete/:id", controllers.Delete)

	router.GET("/preview/:id", controllers.Preview)
	router.GET("/sample/:id", controllers.Sample)
	router.GET("/album/:limit/:page", controllers.Album)
	router.GET("/album_prefix/:p", controllers.Prefix)

	router.GET("/search/:tag", controllers.Search)
	router.GET("/remote/:limit/:page/*tag", controllers.Tag)

	router.GET("/dist/:limit", controllers.Dis)

	handler := cors.Default().Handler(router)

	slog.Fatal(http.ListenAndServe(":8080", handler))
}
