package handlers

import (
	"net/http"

	gorillaHandlers "github.com/gorilla/handlers"
)

func CompressHandler(h http.Handler) http.Handler {
	return gorillaHandlers.CompressHandler(h)
}
