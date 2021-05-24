package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/ctrlaltdev/imgen/img"
	"github.com/ctrlaltdev/imgen/utils"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var (
	PORT int
)

func LogMiddleware(next http.Handler) http.Handler {
	return handlers.CombinedLoggingHandler(os.Stdout, next)
}

func CacheHeader(w http.ResponseWriter) {
	w.Header().Set("Cache-Control", "private, max-age=3600, s-maxage=900, proxy-revalidate, no-transform")
}

func ImageHandler(w http.ResponseWriter, r *http.Request) {
	var (
		format string
		width  int
		height int
	)
	vars := mux.Vars(r)

	if vars["format"] != "" {
		format = vars["format"]
	} else {
		format = "svg"
	}

	if vars["width"] != "" && vars["height"] != "" {
		width64, err := strconv.ParseInt(vars["width"], 10, 64)
		utils.CheckErr(err)
		height64, err := strconv.ParseInt(vars["height"], 10, 64)
		utils.CheckErr(err)

		width = int(width64)
		height = int(height64)
	} else {
		width = 1920
		height = 1080
	}

	switch format {
	case "svg":
		CacheHeader(w)
		w.Header().Set("Content-Type", "image/svg+xml")
		err := img.GenSVG(w, width, height)
		utils.HTTPCheckErr(w, err)
	case "png":
		CacheHeader(w)
		w.Header().Set("Content-Type", "image/png")
		err := img.GenPNG(w, width, height)
		utils.HTTPCheckErr(w, err)
	case "jpg":
		CacheHeader(w)
		w.Header().Set("Content-Type", "image/jpeg")
		err := img.GenJPG(w, width, height)
		utils.HTTPCheckErr(w, err)
	default:
		http.Error(w, fmt.Sprintf("%s is not a support format", format), http.StatusBadRequest)
		return
	}
}

func main() {
	portStr, portSet := os.LookupEnv("PORT")
	if portSet {
		port, err := strconv.ParseInt(portStr, 10, 64)
		utils.CheckErr(err)
		PORT = int(port)
	} else {
		PORT = 3000
	}

	r := mux.NewRouter()
	r.StrictSlash(true)

	r.Use(LogMiddleware)

	r.HandleFunc("/", ImageHandler).Methods("GET", "HEAD")
	r.HandleFunc("/{format:(?:svg|png)}/", ImageHandler).Methods("GET", "HEAD")
	r.HandleFunc("/{width:[0-9]+}/{height:[0-9]+}/", ImageHandler).Methods("GET", "HEAD")
	r.HandleFunc("/{format:(?:svg|png|jpg)}/{width:[0-9]+}/{height:[0-9]+}/", ImageHandler).Methods("GET", "HEAD")

	srv := &http.Server{
		Handler:      handlers.CompressHandler(r),
		Addr:         fmt.Sprintf(":%d", PORT),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	fmt.Printf("\n")
	log.Println(fmt.Sprintf("starting server on port %d", PORT))
	fmt.Printf("\n")

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			fmt.Printf("\n")
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c

	srv.Shutdown(context.Background())

	fmt.Printf("\n")
	log.Println("stopping server")
	fmt.Printf("\n")
	os.Exit(0)
}
