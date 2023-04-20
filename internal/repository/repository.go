package repository

import (
	"database/sql"
	"github.com/calvarado2004/go-movies-backend/internal/models"
)

// DatabaseRepo is a wrapper around the database connection pool.
type DatabaseRepo interface {
	AllMovies() ([]*models.Movie, error)
	Connection() *sql.DB
	GetUserByEmail(email string) (models.User, error)
	GetUserByID(id int) (models.User, error)
}
