package main

import (
	"log"
	"net/http"

	"github.com/gen4ralz/movie-app/rating-service/internal/controller/rating"
	httphandler "github.com/gen4ralz/movie-app/rating-service/internal/handler/http"
	"github.com/gen4ralz/movie-app/rating-service/internal/repository/memory"
)

func main() {
	log.Println("Starting the rating service")

	repo := memory.New()

	ctrl := rating.New(repo)

	h := httphandler.New(ctrl)

	http.Handle("/rating", http.HandlerFunc(h.Handle))

	err := http.ListenAndServe(":8082", nil)
	if err != nil {
		panic(err)
	}
}