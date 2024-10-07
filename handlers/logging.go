package handlers

import (
	"net"
	"net/http"
	"time"

	"github.com/ctrlaltdev/imgen/logger"
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

		logger.Info("request",
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
