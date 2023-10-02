package main

import (
	"log"
	"net/http"

	"github.com/gen4ralz/movie-app/movie-service/internal/controller/movie"
	metadataGateway "github.com/gen4ralz/movie-app/movie-service/internal/gateway/metadata/http"
	ratingGateway "github.com/gen4ralz/movie-app/movie-service/internal/gateway/rating/http"
	httphandler "github.com/gen4ralz/movie-app/movie-service/internal/handler/http"
)

func main() {
	log.Println("Starting the movie service")

	metadataGateway := metadataGateway.New("localhost:8081")
	ratingGateway := ratingGateway.New("localhost:8082")

	ctrl := movie.New(ratingGateway, metadataGateway)

	h := httphandler.New(ctrl)
	http.Handle("/movie", http.HandlerFunc(h.GetMovieDetails))

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}