package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/ctrlaltdev/imgen/img"
	"github.com/ctrlaltdev/imgen/utils"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var (
	PORT int
)

type ResponseWriter struct {
	http.ResponseWriter
	Status   int
	BodySize int
}

func LoggingResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{w, http.StatusOK, 0}
}

func (rw *ResponseWriter) WriteHeader(code int) {
	rw.Status = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *ResponseWriter) Write(b []byte) (int, error) {
	size, err := rw.ResponseWriter.Write(b)
	rw.BodySize += size
	return size, err
}

func LogMiddleware(next http.Handler) http.Handler {
	logger := utils.CreateLogger()
	defer logger.Sync()
	sugar := logger.Sugar()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		lrw := LoggingResponseWriter(w)
		next.ServeHTTP(lrw, r)

		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			ip = r.RemoteAddr
		}
		if r.Header.Get("X-Forwarded-For") != "" {
			ip = r.Header.Get("X-Forwarded-For")
		}

		sugar.Infow("request",
			"host", r.Host,
			"ip", ip,
			"method", r.Method,
			"path", r.URL.Path,
			"proto", r.Proto,
			"query", r.URL.RawQuery,
			"referer", r.Referer(),
			"request_length", r.ContentLength,
			"response_length", lrw.BodySize,
			"response_time", time.Since(startTime),
			"status_code", lrw.Status,
			"url", r.RequestURI,
			"user_agent", r.UserAgent(),
		)
	})
}

func CacheHeader(w http.ResponseWriter) {
	w.Header().Set("Cache-Control", "private, max-age=3600, s-maxage=900, proxy-revalidate, no-transform")
}

func ImageHandler(w http.ResponseWriter, r *http.Request) {
	logger := utils.CreateLogger()
	defer logger.Sync()
	sugar := logger.Sugar()

	var (
		format string
		width  int
		height int
	)
	vars := mux.Vars(r)

	sugar.Debugw("image handler",
		"vars", vars,
	)

	if vars["format"] != "" {
		format = vars["format"]
	} else {
		format = "svg"
	}

	if vars["width"] != "" && vars["height"] != "" {
		width64, err := strconv.ParseInt(vars["width"], 10, 64)
		utils.CheckErr(err, sugar)
		height64, err := strconv.ParseInt(vars["height"], 10, 64)
		utils.CheckErr(err, sugar)

		width = int(width64)
		height = int(height64)
	} else {
		width = 1920
		height = 1080
	}

	sugar.Debugw("image handler",
		"format", format,
		"width", width,
		"height", height,
	)

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

type catchAllHandler struct{}

func (h catchAllHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) > 1 {
		mux.Vars(r)["format"] = pathParts[1]
	}
	if len(pathParts) > 2 {
		mux.Vars(r)["width"] = pathParts[2]
	}
	if len(pathParts) > 3 {
		mux.Vars(r)["height"] = pathParts[3]
	}

	ImageHandler(w, r)
}

func main() {
	logger := utils.CreateLogger()
	defer logger.Sync()
	sugar := logger.Sugar()

	portStr, portSet := os.LookupEnv("PORT")
	if portSet {
		port, err := strconv.ParseInt(portStr, 10, 64)
		utils.CheckErr(err, sugar)
		PORT = int(port)
	} else {
		PORT = 3000
	}

	r := mux.NewRouter()
	r.StrictSlash(true)

	r.Use(LogMiddleware)

	r.HandleFunc("/{format:[a-zA-Z]+}/", ImageHandler).Methods("GET", "HEAD")
	r.HandleFunc("/{width:[0-9]+}/{height:[0-9]+}/", ImageHandler).Methods("GET", "HEAD")
	r.HandleFunc("/{format:[a-zA-Z]+}/{width:[0-9]+}/{height:[0-9]+}/", ImageHandler).Methods("GET", "HEAD")
	r.PathPrefix("/").Handler(catchAllHandler{}).Methods("GET", "HEAD")

	srv := &http.Server{
		Handler:      handlers.CompressHandler(r),
		Addr:         fmt.Sprintf(":%d", PORT),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	sugar.Info(fmt.Sprintf("starting server on port %d", PORT))

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			sugar.Error(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c

	srv.Shutdown(context.Background())

	sugar.Info("stopping server")
	os.Exit(0)
}
