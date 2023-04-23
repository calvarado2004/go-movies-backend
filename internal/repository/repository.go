package repository

import (
	"database/sql"
	"github.com/calvarado2004/go-movies-backend/internal/models"
)

// DatabaseRepo is a wrapper around the database connection pool.
type DatabaseRepo interface {
	AllMovies() ([]*models.Movie, error)
	OneMovie(id int) (*models.Movie, error)
	OneMovieForEdit(id int) (*models.Movie, []*models.Genre, error)
	Connection() *sql.DB
	GetUserByEmail(email string) (models.User, error)
	GetUserByID(id int) (models.User, error)
	AllGenresDB() ([]*models.Genre, error)
	InsertMovie(movie models.Movie) (int, error)
	UpdateMovieGenres(id int, genreIDs []int) error
}
