package utils

import (
	"net/http"

	"go.uber.org/zap"
)

func CheckErr(err error, logger *zap.SugaredLogger) {
	if err != nil {
		logger.Fatal(err)
	}
}

func HTTPCheckErr(w http.ResponseWriter, err error) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func Max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

func Min(x, y int) int {
	if x > y {
		return y
	}
	return x
}
