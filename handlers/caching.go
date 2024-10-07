package handlers

import "net/http"

func CacheHeader(w http.ResponseWriter) {
	w.Header().Set("Cache-Control", "private, max-age=3600, s-maxage=900, proxy-revalidate, no-transform")
}
