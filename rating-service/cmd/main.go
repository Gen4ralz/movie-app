package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/gen4ralz/movie-app/gen"
	"github.com/gen4ralz/movie-app/pkg/discovery"
	"github.com/gen4ralz/movie-app/pkg/discovery/consul"
	"github.com/gen4ralz/movie-app/rating-service/internal/controller/rating"
	grpchandler "github.com/gen4ralz/movie-app/rating-service/internal/handler/grpc"
	"github.com/gen4ralz/movie-app/rating-service/internal/repository/mysql"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gopkg.in/yaml.v3"
)

const serviceName = "rating"

func main() {
	var port int

	fi, err := os.Open("base.yaml")
	if err != nil {
		panic(err)
	}
	defer fi.Close()

	var cfg serviceConfig
	
	err = yaml.NewDecoder(fi).Decode(&cfg)
	if err != nil {
		panic(err)
	}

	port = cfg.APIConfig.Port
	
	log.Printf("Starting the rating service on port %d", port)

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
			time.Sleep(1 *time.Second)
		}
	}()
	defer registry.Deregister(ctx, instanceID, serviceName)

	// Creates a new instance of a MySQL-based repository for storing rating.
	repo, err := mysql.New()
	if err != nil {
		panic(err)
	}

	// Creates a controller instance for the rating service, passing in the memory-based repository.
	ctrl := rating.New(repo, nil)

	// Creates a gRPC handler instance, initializing it with the rating controller.
	h := grpchandler.New(ctrl)

	// Creates a TCP listener on the specified port. 
	// If there's an error creating the listener, logs an error message and terminates the program.
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%v", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Creates a new gRPC server instance.
	srv := grpc.NewServer()

	// Registers reflection support on the gRPC server, 
	// which allows clients to dynamically discover the available gRPC services.
	reflection.Register(srv)

	// Registers the gRPC rating service generated from the gen package with the gRPC server. 
	// It uses the h (handler) instance to handle incoming gRPC requests.
	gen.RegisterRatingServiceServer(srv, h)

	// Starts the gRPC server to listen for incoming requests on the previously created listener (lis). 
	// If there's an error during server startup, it panics.
	err = srv.Serve(lis)
	if err != nil {
		log.Panic(err)
	}
}