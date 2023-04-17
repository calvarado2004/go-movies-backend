package main

import (
	"fmt"
	"net/http"
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

	err := app.writeJSON(w, http.StatusOK, payload, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

}

// AllMovies is a simple handler function which writes a response.
func (app *application) AllMovies(w http.ResponseWriter, r *http.Request) {

	// get all movies from the database
	movies, err := app.DB.AllMovies()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, movies, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	
}
