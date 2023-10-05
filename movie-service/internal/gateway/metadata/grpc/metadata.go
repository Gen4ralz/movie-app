package grpc

import (
	"context"

	"github.com/gen4ralz/movie-app/gen"
	"github.com/gen4ralz/movie-app/internal/grpcutil"
	"github.com/gen4ralz/movie-app/metadata-service/pkg/model"
	"github.com/gen4ralz/movie-app/pkg/discovery"
)

// Gateway defines a movie metadata gRPC gateway.
type Gateway struct {
	registry	discovery.Registry
}

// New creates a new gRPC gateway for a movie metadata service
func New(reg discovery.Registry) *Gateway {
	return &Gateway{
		registry: reg,
	}
}

// Get returns movie metadata by a movie id.
func (g *Gateway) Get(ctx context.Context, id string) (*model.Metadata, error) {
	conn, err := grpcutil.ServiceConnection(ctx, "metadata", g.registry)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := gen.NewMetadataServiceClient(conn)

	response, err := client.GetMetadata(ctx, &gen.GetMetadataRequest{
		MovieId: id,
	})
	if err != nil {
		return nil, err
	}

	return model.MetadataFromProto(response.Metadata), nil
}