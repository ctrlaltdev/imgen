package utils

import (
	"net/http"
)

func CheckErr(err error) {
	if err != nil {
		panic(err)
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
