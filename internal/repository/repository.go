package repository

import "github.com/calvarado2004/go-movies-backend/internal/models"

// DatabaseRepo is a wrapper around the database connection pool.
type DatabaseRepo interface {
	AllMovies() ([]*models.Movie, error)
}