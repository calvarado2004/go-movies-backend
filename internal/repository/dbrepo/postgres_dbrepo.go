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

// OneMovie returns one movie from the database.
func (m *PostgresDBRepo) OneMovie(id int) (*models.Movie, error) {

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `SELECT id, title, release_date, runtime, mpaa_rating, description, coalesce(image, ''), created_at, updated_at FROM movies WHERE id = $1`

	row := m.DB.QueryRowContext(ctx, query, id)

	var movie models.Movie

	err := row.Scan(
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

	// get genres for this movie

	query = `SELECT g.id, g.genre from movies_genres mg, genres g where movie_id = $1 and g.id = mg.genre_id`

	rows, err := m.DB.QueryContext(ctx, query, id)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			return
		}
	}(rows)

	var genres []*models.Genre

	for rows.Next() {
		var genre models.Genre
		err := rows.Scan(
			&genre.ID,
			&genre.Genre,
		)
		if err != nil {
			return nil, err
		}
		genres = append(genres, &genre)
	}

	movie.Genres = genres

	return &movie, nil

}

// OneMovieForEdit returns one movie from the database and edit it.
func (m *PostgresDBRepo) OneMovieForEdit(id int) (*models.Movie, []*models.Genre, error) {

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `SELECT id, title, release_date, runtime, mpaa_rating, description, coalesce(image, ''), created_at, updated_at FROM movies WHERE id = $1`

	row := m.DB.QueryRowContext(ctx, query, id)

	var movie models.Movie

	err := row.Scan(
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
		return nil, nil, err
	}

	// get genres for this movie

	query = `SELECT g.id, g.genre from movies_genres mg, genres g where movie_id = $1 and g.id = mg.genre_id`

	rows, err := m.DB.QueryContext(ctx, query, id)
	if err != nil && err != sql.ErrNoRows {
		return nil, nil, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			return
		}
	}(rows)

	var genres []*models.Genre

	var genresArray []int

	for rows.Next() {
		var genre models.Genre
		err := rows.Scan(
			&genre.ID,
			&genre.Genre,
		)
		if err != nil {
			return nil, nil, err
		}
		genres = append(genres, &genre)
		genresArray = append(genresArray, genre.ID)
	}

	movie.Genres = genres
	movie.GenresArray = genresArray

	var allGenres []*models.Genre

	query = `SELECT id, genre FROM genres ORDER BY genre`

	rows, err = m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, nil, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			return
		}
	}(rows)

	for rows.Next() {
		var genre models.Genre
		err := rows.Scan(
			&genre.ID,
			&genre.Genre,
		)
		if err != nil {
			return nil, nil, err
		}
		allGenres = append(allGenres, &genre)
	}

	return &movie, allGenres, nil

}

// GetUserByEmail returns a user from the database by email.
func (m *PostgresDBRepo) GetUserByEmail(email string) (models.User, error) {

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `SELECT id, email, first_name, last_name, password, created_at, updated_at FROM users WHERE email = $1`

	var user models.User

	row := m.DB.QueryRowContext(ctx, query, email)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

// GetUserByID returns a user from the database by id.
func (m *PostgresDBRepo) GetUserByID(id int) (models.User, error) {

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `SELECT id, email, first_name, last_name, password, created_at, updated_at FROM users WHERE id = $1`

	var user models.User

	row := m.DB.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

// AllGenresDB returns all genres from the database.
func (m *PostgresDBRepo) AllGenresDB() ([]*models.Genre, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `SELECT id, genre, created_at, updated_at FROM genres order by genre`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	var genres []*models.Genre

	for rows.Next() {
		var genre models.Genre
		err := rows.Scan(
			&genre.ID,
			&genre.Genre,
			&genre.CreatedAt,
			&genre.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		genres = append(genres, &genre)
	}

	return genres, nil
}
