package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/gen4ralz/movie-app/gen"
	"github.com/gen4ralz/movie-app/metadata-service/internal/controller/metadata"
	grpchandler "github.com/gen4ralz/movie-app/metadata-service/internal/handler/grpc"
	"github.com/gen4ralz/movie-app/metadata-service/internal/repository/mysql"
	"github.com/gen4ralz/movie-app/pkg/discovery"
	"github.com/gen4ralz/movie-app/pkg/discovery/consul"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gopkg.in/yaml.v3"
)

const serviceName = "metadata"

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

	log.Printf("Starting the metadata service on port %d", port)

	// Creates a new Consul service registry instance connected to the Consul server running on "localhost:8500." 
	// If there's an error, it panics (terminates the program with an error message).
	registry, err := consul.NewRegistry("localhost:8500")
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	// Generates a unique instance ID for the metadata service using the GenerateInstanceID function from the discovery package. 
	// This ID will be used for service registration with Consul.
	instanceID := discovery.GenerateInstanceID(serviceName)

	// Registers the metadata service with Consul using the generated instance ID, service name, 
	// and the address formed by combining "localhost" with the port value. 
	// If there's an error during registration, it panics.
	err = registry.Register(ctx, instanceID, serviceName, fmt.Sprintf("localhost:%d", port))
	if err != nil {
		panic(err)
	}

	// This code runs a goroutine (concurrent function) that periodically reports the healthy state of the service to Consul. 
	// It logs an error message if reporting fails and sleeps for one second between reports.
	go func() {
		for {
			err := registry.ReportHealthyState(instanceID, serviceName)
			if err != nil {
				log.Println("Failed to report healthy state:" + err.Error())
			}
			time.Sleep(1 * time.Second)
		}
	}()

	// Defers the deregistration of the service with Consul until the program exits.
	// This ensures that the service is deregistered properly even if the program terminates unexpectedly.
	defer registry.Deregister(ctx, instanceID, serviceName)

	// Creates a new instance of a MySQL repository for storing metadata.
	repo, err := mysql.New()
	if err != nil {
		panic(err)
	}

	// Creates a controller instance for the metadata service, passing in the memory-based repository.
	ctrl := metadata.New(repo)

	// Creates a gRPC handler instance, initializing it with the metadata controller.
	h := grpchandler.New(ctrl)

	// Creates a TCP listener on the specified port. 
	// If there's an error creating the listener, logs an error message and terminates the program.
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Fatalf("Failed to listen: %d", err)
	}
	
	// Creates a new gRPC server instance.
	srv := grpc.NewServer()

	// Registers reflection support on the gRPC server, 
	// which allows clients to dynamically discover the available gRPC services.
	reflection.Register(srv)

	// Registers the gRPC metadata service generated from the gen package with the gRPC server. 
	// It uses the h (handler) instance to handle incoming gRPC requests.
	gen.RegisterMetadataServiceServer(srv, h)

	// Starts the gRPC server to listen for incoming requests on the previously created listener (lis). 
	// If there's an error during server startup, it panics.
	if err := srv.Serve(lis); err != nil {
		panic(err)
	}
}