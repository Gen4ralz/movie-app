package memory

import (
	"context"

	"github.com/gen4ralz/movie-app/rating-service/internal/repository"
	"github.com/gen4ralz/movie-app/rating-service/pkg/model"
)

// Repository defines a rating repository.
type Repository struct {
	data	map[model.RecordType]map[model.RecordID][]model.Rating
}

// New creates a new memory repository.
func New() *Repository {
	return &Repository{
		map[model.RecordType]map[model.RecordID][]model.Rating{},
	}
}

// Get retrieves all ratings for a given record.
func (r *Repository) Get(ctx context.Context, recordID model.RecordID, recordType model.RecordType) ([]model.Rating, error) {
	_, ok := r.data[recordType]
	if !ok {
		return nil, repository.ErrNotFound
	}

	ratings, ok := r.data[recordType][recordID]
	if !ok || len(ratings) == 0 {
		return nil, repository.ErrNotFound
	}

	return r.data[recordType][recordID], nil
}

// Put adds a rating for a given record.
func (r *Repository) Put(ctx context.Context, recordID model.RecordID, recordType model.RecordType, rating *model.Rating) error {
	_, ok := r.data[recordType]
	if !ok {
		r.data[recordType] = map[model.RecordID][]model.Rating{}
	}

	r.data[recordType][recordID] = append(r.data[recordType][recordID], *rating)
	return nil
}

