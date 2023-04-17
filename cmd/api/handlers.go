package main

import (
	"encoding/json"
	"fmt"
	"github.com/calvarado2004/go-movies-backend/internal/models"
	"net/http"
	"time"
)

// Home is a simple handler function which writes a response.
func (app *application) Home(w http.ResponseWriter, r *http.Request) {

	var payload = struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Version string `json:"version"`
	}{
		Status:  "active",
		Message: "Welcome to the Go Movies API",
		Version: "1.0.0",
	}

	out, err := json.Marshal(payload)
	if err != nil {
		fmt.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(out)
	if err != nil {
		fmt.Println(err)
		return
	}
}

// AllMovies is a simple handler function which writes a response.
func (app *application) AllMovies(w http.ResponseWriter, r *http.Request) {
	// get all movies
	var movies []models.Movie

	releaseDate, err := time.Parse("2006-01-02", "1986-03-07")
	if err != nil {
		fmt.Println(err)
	}

	// create one movie
	highlander := models.Movie{
		ID:          1,
		Title:       "Highlander",
		ReleaseDate: releaseDate,
		Runtime:     116,
		MPAARating:  "R",
		Description: "An immortal Scottish swordsman must confront the last of his immortal opponent, a murderously brutal barbarian who lusts for the fabled Prize.",
		Genre:       "Action, Adventure, Fantasy, Romance",
		Image:       "https://m.media-amazon.com/images/M/MV5BMWY4MDg1NmUtMDg1ZC00ODUyLWJkOWQtNWIxMWZhMjRlM2Q1XkEyXkFqcGdeQXVyODc0OTEyNDU@._V1_FMjpg_UX1000_.jpg",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// add movie to movies
	movies = append(movies, highlander)

	releaseDate, err = time.Parse("2006-01-02", "1981-06-12")
	if err != nil {
		fmt.Println(err)
	}

	// create one movie
	rotla := models.Movie{
		ID:          1,
		Title:       "Raiders of the Lost Ark",
		ReleaseDate: releaseDate,
		Runtime:     115,
		MPAARating:  "PG-13",
		Description: "Raiders of the Lost Ark is a 1981 American action adventure film directed by Steven Spielberg, produced by George Lucas, and starring Harrison Ford. It was the first installment in the Indiana Jones franchise.",
		Genre:       "Action, Adventure",
		Image:       "https://image.tmdb.org/t/p/w500/4q2hz2m8hubgvijz8Ez0T2Os2Yv.jpg",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// add movie to movies
	movies = append(movies, rotla)

	out, err := json.Marshal(movies)
	if err != nil {
		fmt.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(out)
	if err != nil {
		fmt.Println(err)
		return
	}
}
