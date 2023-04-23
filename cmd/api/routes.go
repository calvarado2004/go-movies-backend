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

	// add CORS middleware
	mux.Use(app.enableCORS)

	// add routes
	mux.Get("/", app.Home)
	mux.Get("/movies", app.AllMovies)
	mux.Get("/movies/{id}", app.getMovie)
	mux.Post("/authenticate", app.authenticate)
	mux.Get("/refresh", app.refreshToken)
	mux.Get("/logout", app.logout)
	mux.Get("/genres", app.allGenres)

	mux.Route("/admin", func(authMux chi.Router) {
		authMux.Use(app.authRequired)
		authMux.Get("/movies", app.movieCatalog)
		authMux.Get("/movies/{id}", app.movieForEdit)
		authMux.Put("/movies/0", app.insertMovie)
	})

	return mux

}
