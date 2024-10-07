package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/ctrlaltdev/imgen/handlers"
	"github.com/ctrlaltdev/imgen/logger"
	"github.com/ctrlaltdev/imgen/utils"

	"github.com/gorilla/mux"
)

var (
	PORT int
)

func main() {
	portStr, portSet := os.LookupEnv("PORT")
	if portSet {
		var err error

		PORT, err = strconv.Atoi(portStr)
		utils.CheckErr(err)
	} else {
		PORT = 3000
	}

	r := mux.NewRouter()
	r.StrictSlash(true)

	r.Use(handlers.LogMiddleware)

	r.HandleFunc("/{format:[a-zA-Z]+}/", handlers.ImageHandler).Methods("GET", "HEAD")
	r.HandleFunc("/{width:[0-9]+}/{height:[0-9]+}/", handlers.ImageHandler).Methods("GET", "HEAD")
	r.HandleFunc("/{format:[a-zA-Z]+}/{width:[0-9]+}/{height:[0-9]+}/", handlers.ImageHandler).Methods("GET", "HEAD")
	r.PathPrefix("/").Handler(handlers.CatchAllHandler{}).Methods("GET", "HEAD")

	srv := &http.Server{
		Handler:      handlers.CompressHandler(r),
		Addr:         fmt.Sprintf(":%d", PORT),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	logger.Info(fmt.Sprintf("starting server on port %d", PORT))

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logger.Fatal("fail to start the server", err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c

	srv.Shutdown(context.Background())

	logger.Info("stopping server")
	os.Exit(0)
}
