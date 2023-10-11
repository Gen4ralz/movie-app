package mysql

import (
	"context"
	"database/sql"

	"github.com/gen4ralz/movie-app/metadata-service/internal/repository"
	"github.com/gen4ralz/movie-app/metadata-service/pkg/model"
	_ "github.com/go-sql-driver/mysql"
)

// Repository defines a MySQL-based movie metadata repository.
type Repository struct {
	db *sql.DB
}

// New creates a new MySQL-based repository.
func New() (*Repository, error) {
	db, err := sql.Open("mysql", "root:password@/movie")
	if err != nil {
		return nil, err
	}

	return &Repository{
		db: db,
	}, nil
}

// Get retrieves movie metadata for by movie id.
func (r *Repository) Get(ctx context.Context, id string) (*model.Metadata, error) {
	var title, description, director string

	query := `SELECT title, description, director FROM movies WHERE id = ?`
	row := r.db.QueryRowContext(ctx, query, id)
	err := row.Scan(&title, &description, &director)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}

	return &model.Metadata{
		ID: id,
		Title: title,
		Description: description,
		Director: director,
	}, nil
}

// Put adds movie metadata for a given movie id.
func (r *Repository) Put(ctx context.Context, id string, metadata *model.Metadata) error {
	query := `INSERT INTO movies (id, title, description, director) VALUES (?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, id, metadata.Title, metadata.Description, metadata.Director)
	return err
}