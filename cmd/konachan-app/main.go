package main

import (
	"bufio"
	"flag"
	"fmt"
	slog "log"
	"net/http"
	"os"
	"os/exec"
	"runtime"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/CheerChen/konachan-app/internal/asset"
	"github.com/CheerChen/konachan-app/internal/controllers"
	"github.com/CheerChen/konachan-app/internal/kfile"
	"github.com/CheerChen/konachan-app/internal/log"
	"github.com/CheerChen/konachan-app/internal/models"
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	bytes, err := asset.Asset("web/static/index.html")
	if err != nil {
		log.Errorf("", err)
		http.NotFound(w, r)
		return
	}
	w.Write(bytes)
}

var (
	level = flag.Int("l", -1, "log level, -1:debug, 0:info, 1:warn, 2:error")
	path  = flag.String("p", "path/to/wallpaper", "wallpaper path, input an absolute path")
)

func main() {

	flag.Parse()

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

	if _, err := os.Stat(*path); err != nil {
		*path = GetUserInput()
	}
	go kfile.Sync(*path)

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

	go Open("http://localhost:7079/")
	slog.Fatal(http.ListenAndServe(":7079", handler))
}

func GetUserInput() (input string) {
	for {
		consoleReader := bufio.NewReader(os.Stdin)
		fmt.Print("Please input wallpaper path >")

		line, _, _ := consoleReader.ReadLine()
		if _, err := os.Stat(string(line)); err != nil {
			log.Errorf("Wallpaper path err: %v", err)
			continue
		} else {
			input = string(line)
			break
		}
	}
	return
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
