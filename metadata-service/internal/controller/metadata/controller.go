package metadata

import (
	"context"
	"errors"

	"github.com/gen4ralz/movie-app/metadata-service/internal/repository"
	"github.com/gen4ralz/movie-app/metadata-service/pkg/model"
)

// ErrNotFound is returned when a requested record is not // found.
var ErrNotFound = errors.New("not found")

type metadataRepository interface {
	Get(ctx context.Context, id string) (*model.Metadata, error)
}

// Controller defines a metadata service controller.
type Controller struct {
	repo	metadataRepository
}

// New creates a metadata service controller.
func New(repo metadataRepository) *Controller {
	return &Controller{
		repo: repo,
	}
}

// Get returns movies metadata by movie id.
func (c *Controller) Get(ctx context.Context, id string) (*model.Metadata, error) {
	res, err := c.repo.Get(ctx, id)
	if err != nil && errors.Is(err, repository.ErrNotFound) {
		return nil, ErrNotFound
	}

	return res, err
}