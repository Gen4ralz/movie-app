package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gen4ralz/movie-app/movie-service/internal/controller/movie"
	metadataGateway "github.com/gen4ralz/movie-app/movie-service/internal/gateway/metadata/http"
	ratingGateway "github.com/gen4ralz/movie-app/movie-service/internal/gateway/rating/http"
	httphandler "github.com/gen4ralz/movie-app/movie-service/internal/handler/http"
	"github.com/gen4ralz/movie-app/pkg/discovery"
	"github.com/gen4ralz/movie-app/pkg/discovery/consul"
)

const serviceName = "movie"

func main() {
	var port int
	flag.IntVar(&port, "port", 8080, "API handler port")
	flag.Parse()
	log.Printf("Starting the movie service on port %d", port)
	registry, err := consul.NewRegistry("localhost:8500")
	if err != nil {
		panic(err)
	}
	ctx := context.Background()

	instanceID := discovery.GenerateInstanceID(serviceName)
	
	if err := registry.Register(ctx, instanceID, serviceName, fmt.Sprintf("localhost:%d", port)); err != nil {
		panic(err)
	}
	go func() {
		for {
			if err := registry.ReportHealthyState(instanceID, serviceName); err != nil {
				log.Println("Failed to report healthy state: " + err.Error())
			}
			time.Sleep(1 * time.Second)
		}
	}()
	defer registry.Deregister(ctx, instanceID, serviceName)

	metadataGateway := metadataGateway.New(registry)
	ratingGateway := ratingGateway.New(registry)

	svc := movie.New(ratingGateway, metadataGateway)

	h := httphandler.New(svc)
	http.Handle("/movie", http.HandlerFunc(h.GetMovieDetails))

	err = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		panic(err)
	}
}