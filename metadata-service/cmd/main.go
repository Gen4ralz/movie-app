package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gen4ralz/movie-app/metadata-service/internal/controller/metadata"
	httphandler "github.com/gen4ralz/movie-app/metadata-service/internal/handler/http"
	"github.com/gen4ralz/movie-app/metadata-service/internal/repository/memory"
	"github.com/gen4ralz/movie-app/pkg/discovery"
	"github.com/gen4ralz/movie-app/pkg/discovery/consul"
)

const serviceName = "metadata"

func main() {
	var port int

	flag.IntVar(&port, "port", 8081, "API handler port")
	flag.Parse()


	log.Printf("Starting the metadata service on port %d", port)

	registry, err := consul.NewRegistry("localhost:8500")
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)

	err = registry.Register(ctx, instanceID, serviceName, fmt.Sprintf("localhost:%d", port))
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			err := registry.ReportHealthyState(instanceID, serviceName)
			if err != nil {
				log.Println("Failed to report healthy state:" + err.Error())
			}
			time.Sleep(1 * time.Second)
		}
	}()
	defer registry.Deregister(ctx, instanceID, serviceName)

	repo := memory.New()

	ctrl := metadata.New(repo)

	h := httphandler.New(ctrl)

	http.Handle("/metadata", http.HandlerFunc(h.GetMetadata))

	err = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		panic(err)
	}
}