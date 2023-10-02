package movie

import (
	"context"
	"errors"

	metadataModel "github.com/gen4ralz/movie-app/metadata-service/pkg/model"
	"github.com/gen4ralz/movie-app/movie-service/internal/gateway"
	"github.com/gen4ralz/movie-app/movie-service/pkg/model"
	ratingModel "github.com/gen4ralz/movie-app/rating-service/pkg/model"
)

// ErrNotFound is returned when the movie metadata is not found
var ErrNotFound = errors.New("movie metadata not found")

type ratingGateway interface {
	GetAggregatedRating(ctx context.Context, recordID ratingModel.RecordID, recordType ratingModel.RecordType) (float64, error)
	PutRating(ctx context.Context, recordID ratingModel.RecordID, recordType ratingModel.RecordType, rating *ratingModel.Rating) error
}

type metadataGateway interface {
	Get(ctx context.Context, id string) (*metadataModel.Metadata, error)
}

// Controller defines a movie service controller.
type Controller struct {
	ratingGateway	ratingGateway
	metadataGateway	metadataGateway
}

// New creates a new movie service controller
func New(r ratingGateway, m metadataGateway) *Controller {
	return &Controller{
		ratingGateway: r,
		metadataGateway: m,
	}
}

// Get returns the movie details including the aggregated rating and movie metadata
func (c *Controller) Get(ctx context.Context, id string) (*model.MovieDetails, error) {
	metadata, err := c.metadataGateway.Get(ctx, id)
	if err != nil && errors.Is(err, gateway.ErrNotFound){
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	details := &model.MovieDetails{
		Metadata: *metadata,
	}

	rating, err := c.ratingGateway.GetAggregatedRating(ctx, ratingModel.RecordID(id), ratingModel.RecordTypeMovie)
	if err != nil && !errors.Is(err, gateway.ErrNotFound) {
		// Proceed in this case, it's ok not to have ratings yet.
	} else if err != nil {
		return nil ,err 
	} else {
		details.Rating = &rating
	}
	return details, nil
}