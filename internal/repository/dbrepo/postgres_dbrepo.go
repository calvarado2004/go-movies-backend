package dbrepo

import (
	"context"
	"database/sql"
	"github.com/calvarado2004/go-movies-backend/internal/models"
	"time"
)

// PostgresDBRepo is a wrapper around the database connection pool.
type PostgresDBRepo struct {
	DB *sql.DB
}

// dbTimeout is the maximum amount of time a database operation can take.
const dbTimeout = time.Second * 5

// Connection returns the database connection pool.
func (m *PostgresDBRepo) Connection() *sql.DB {
	return m.DB
}

// AllMovies returns all movies from the database.
func (m *PostgresDBRepo) AllMovies() ([]*models.Movie, error) {

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var movies []*models.Movie

	query := `SELECT 
    	id, title, release_date, runtime, mpaa_rating, description, coalesce(image, ''), created_at, updated_at 
	FROM 
	    movies 
	ORDER BY 
	    title DESC`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			return

		}
	}(rows)

	for rows.Next() {
		movie := models.Movie{}
		err := rows.Scan(
			&movie.ID,
			&movie.Title,
			&movie.ReleaseDate,
			&movie.Runtime,
			&movie.MPAARating,
			&movie.Description,
			&movie.Image,
			&movie.CreatedAt,
			&movie.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		movies = append(movies, &movie)
	}

	return movies, nil
}
