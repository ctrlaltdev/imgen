package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/ctrlaltdev/imgen/img"
	"github.com/ctrlaltdev/imgen/logger"
	"github.com/ctrlaltdev/imgen/utils"
	"github.com/gorilla/mux"
)

func ImageHandler(w http.ResponseWriter, r *http.Request) {
	var (
		format string
		width  int
		height int
	)
	vars := mux.Vars(r)

	logger.Debug("image handler",
		"vars", vars,
	)

	if vars["format"] != "" {
		format = vars["format"]
	} else {
		format = "svg"
	}

	if vars["width"] != "" && vars["height"] != "" {
		var err error

		width, err = strconv.Atoi(vars["width"])
		utils.CheckErr(err)

		height, err = strconv.Atoi(vars["height"])
		utils.CheckErr(err)
	} else {
		width = 1920
		height = 1080
	}

	logger.Debug("image handler",
		"format", format,
		"width", width,
		"height", height,
	)

	contentTypeHeader := "Content-Type"

	switch format {
	case "svg":
		CacheHeader(w)
		w.Header().Set(contentTypeHeader, "image/svg+xml")
		err := img.GenSVG(w, width, height)
		utils.HTTPCheckErr(w, err)
	case "png":
		CacheHeader(w)
		w.Header().Set(contentTypeHeader, "image/png")
		err := img.GenPNG(w, width, height)
		utils.HTTPCheckErr(w, err)
	case "jpg":
		CacheHeader(w)
		w.Header().Set(contentTypeHeader, "image/jpeg")
		err := img.GenJPG(w, width, height)
		utils.HTTPCheckErr(w, err)
	default:
		http.Error(w, fmt.Sprintf("%s is not a support format", format), http.StatusBadRequest)
		return
	}
}

type CatchAllHandler struct{}

func (h CatchAllHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
