package utils

import (
	"net/http"

	"github.com/ctrlaltdev/imgen/logger"
)

func CheckErr(err error) {
	if err != nil {
		logger.Fatal("unexpected error", err)
	}
}

func HTTPCheckErr(w http.ResponseWriter, err error) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
