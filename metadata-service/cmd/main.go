package main

import (
	"log"
	"net/http"

	"github.com/gen4ralz/movie-app/metadata-service/internal/controller/metadata"
	httphandler "github.com/gen4ralz/movie-app/metadata-service/internal/handler/http"
	"github.com/gen4ralz/movie-app/metadata-service/internal/repository/memory"
)

func main() {
	log.Println("Starting the movie metadata service")

	repo := memory.New()

	ctrl := metadata.New(repo)

	h := httphandler.New(ctrl)

	http.Handle("/metadata", http.HandlerFunc(h.GetMetadata))

	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		panic(err)
	}
}