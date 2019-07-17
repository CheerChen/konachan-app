package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"

	"github.com/CheerChen/konachan-app/internal/controllers"
	"github.com/CheerChen/konachan-app/internal/kpost"
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Service Available!\n")
}

func main() {

	kpost.InitDB()

	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/post/:id", controllers.GetByIdV2)
	router.GET("/page/:id", controllers.GetById)
	router.GET("/tag/tf_idf", controllers.GetTfIdf)
	//router.GET("/tag/update", Update)
	router.GET("/hot/:limit/:page", controllers.Popular)
	router.GET("/hot_from/:limit/:page/:from", controllers.PopularByRange)
	router.GET("/download/:id", controllers.Download)
	router.GET("/sync", controllers.Sync)
	router.GET("/check", controllers.Check)
	router.GET("/check2", controllers.Sync2)

	router.GET("/delete/:id", controllers.Delete)

	router.GET("/preview/:id", controllers.Preview)
	router.GET("/sample/:id", controllers.Sample)
	router.GET("/album/:limit/:page", controllers.Album)
	router.GET("/album_prefix/:p", controllers.Prefix)

	router.GET("/search/:tag", controllers.Search)

	router.GET("/dist/:limit", controllers.Dis)

	handler := cors.Default().Handler(router)

	log.Fatal(http.ListenAndServe(":8080", handler))
}
