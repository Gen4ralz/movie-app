package grpc

import (
	"context"

	"github.com/gen4ralz/movie-app/gen"
	"github.com/gen4ralz/movie-app/internal/grpcutil"
	"github.com/gen4ralz/movie-app/pkg/discovery"
	"github.com/gen4ralz/movie-app/rating-service/pkg/model"
)

// Gateway defines an gRPC gateway for a rating service.
type Gateway struct {
	registry	discovery.Registry
}

// New creates a new gRPC gateway for a rating service.
func New(reg discovery.Registry) *Gateway {
	return &Gateway{
		registry: reg,
	}
}

// GetAggregatedRating returns the aggregated rating for a record 
// or ErrNotFound if there are no ratings for it.
func (g *Gateway) GetAggregatedRating(ctx context.Context, recordID model.RecordID, recordType model.RecordType) (float64, error) {
	conn, err := grpcutil.ServiceConnection(ctx, "rating", g.registry)
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	client := gen.NewRatingServiceClient(conn)

	response, err := client.GetAggregatedRating(ctx, &gen.GetAggregatedRatingRequest{
		RecordId: string(recordID),
		RecordType: string(recordType),
	})
	if err != nil {
		return 0, nil
	}

	return response.RatingValue, nil
}