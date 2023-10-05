package grpc

import (
	"context"
	"errors"

	"github.com/gen4ralz/movie-app/gen"
	"github.com/gen4ralz/movie-app/rating-service/internal/controller/rating"
	"github.com/gen4ralz/movie-app/rating-service/pkg/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Handler defines a gRPC rating API handler.
type Handle struct {
	gen.UnimplementedRatingServiceServer
	ctrl	*rating.Controller
}

// New creates a new rating gRPC handler.
func New(ctrl *rating.Controller) *Handle {
	return &Handle{
		ctrl: ctrl,
	}
}

// GetAggregatedRating returns the aggregated rating for a record.
func (h *Handle) GetAggregatedRating(ctx context.Context, req *gen.GetAggregatedRatingRequest) (*gen.GetAggregatedRatingResponse, error) {
	if req == nil || req.RecordId == "" || req.RecordType == "" {
		return nil, status.Error(codes.InvalidArgument, "nil req or empty id")
	}

	v, err := h.ctrl.GetAggregatedRating(ctx, model.RecordID(req.RecordId), model.RecordType(req.RecordType))
	if err != nil && errors.Is(err, rating.ErrNotFound) {
		return nil, status.Errorf(codes.NotFound, err.Error())
	} else if err != nil {
		return nil, status.Errorf(codes.NotFound, err.Error())
	}

	return &gen.GetAggregatedRatingResponse{
		RatingValue: v,
	}, nil
}

// PutRating writes a rating for a given record.
func (h *Handle) PutRating(ctx context.Context, req *gen.PutRatingRequest) (*gen.PutRatingResponse, error) {
	if req == nil || req.RecordId == "" || req.UserId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "nil req or empty user id or record id")
	}

	err := h.ctrl.PutRating(ctx, model.RecordID(req.RecordId), model.RecordType(req.RecordType), &model.Rating{
		UserID: model.UserID(req.UserId),
		Value: model.RatingValue(req.RatingValue),
	})
	if err != nil {
		return nil, err
	}

	return &gen.PutRatingResponse{}, nil
}