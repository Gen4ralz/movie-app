package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/gen4ralz/movie-app/gen"
	"github.com/gen4ralz/movie-app/movie-service/internal/controller/movie"
	metadataGateway "github.com/gen4ralz/movie-app/movie-service/internal/gateway/metadata/http"
	ratingGateway "github.com/gen4ralz/movie-app/movie-service/internal/gateway/rating/http"
	grpchandler "github.com/gen4ralz/movie-app/movie-service/internal/handler/grpc"
	"github.com/gen4ralz/movie-app/pkg/discovery"
	"github.com/gen4ralz/movie-app/pkg/discovery/consul"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gopkg.in/yaml.v3"
)

const serviceName = "movie"

func main() {
	var port int

	fi, err := os.Open("base.yaml")
	if err != nil {
		panic(err)
	}

	var cfg serviceConfig
	
	err = yaml.NewDecoder(fi).Decode(&cfg)
	if err != nil {
		panic(err)
	}

	port = cfg.APIConfig.Port
	
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

	h := grpchandler.New(svc)

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	srv := grpc.NewServer()

	reflection.Register(srv)

	gen.RegisterMovieServiceServer(srv, h)

	err = srv.Serve(lis)
	if err != nil {
		panic(err)
	}
}