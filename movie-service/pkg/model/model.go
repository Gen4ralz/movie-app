package model

import "github.com/gen4ralz/movie-app/metadata-service/pkg/model"

// MovieDetails includes movie metadata its aggregated rating.
type MovieDetails struct {
	Rating		*float64			`json:"rating,omitempty"`
	Metadata	model.Metadata	`json:"metadata"`
}