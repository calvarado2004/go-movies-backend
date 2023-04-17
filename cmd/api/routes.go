package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

// routes returns a http.Handler containing all the routes for the application.
func (app *application) routes() http.Handler {
	// create a router mux

	mux := chi.NewRouter()

	// add middleware
	mux.Use(middleware.Recoverer)

	// add routes
	mux.Get("/", app.Home)

	return mux

}
