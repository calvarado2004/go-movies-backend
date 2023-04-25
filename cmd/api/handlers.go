package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/calvarado2004/go-movies-backend/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v4"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
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
		err := app.errorJSON(w, err)
		if err != nil {
			return
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, movies, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

}

// authenticate is a simple handler function which writes a response.
func (app *application) authenticate(w http.ResponseWriter, r *http.Request) {
	// read json payload

	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		err := app.errorJSON(w, err, http.StatusBadRequest)
		if err != nil {
			return
		}
		return
	}

	// validate payload, user exists, password matches
	user, err := app.DB.GetUserByEmail(requestPayload.Email)
	if err != nil {
		err := app.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		if err != nil {
			return
		}
		return
	}

	// check password
	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		err := app.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		if err != nil {
			return
		}
		return
	}

	// generate token pair
	u := jwtUser{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}

	tokens, err := app.auth.generateTokenPair(&u)
	if err != nil {
		err := app.errorJSON(w, err)
		if err != nil {
			return
		}
		return

	}

	refreshCookie := app.auth.getRefreshCookie(tokens.RefreshToken)
	http.SetCookie(w, refreshCookie)

	// write json response
	err = app.writeJSON(w, http.StatusAccepted, tokens, nil)
	if err != nil {
		return
	}

}

// refreshToken is a simple handler function which writes a response.
func (app *application) refreshToken(w http.ResponseWriter, r *http.Request) {
	for _, cookie := range r.Cookies() {
		if cookie.Name == app.auth.CookieName {

			claims := &tokenClaims{}
			refreshToken := cookie.Value

			_, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (any, error) {
				return []byte(app.auth.Secret), nil
			})
			if err != nil {
				err := app.errorJSON(w, errors.New("invalid token"), http.StatusUnauthorized)
				if err != nil {
					return
				}
				return
			}

			// get user id from token claims
			userID, err := strconv.Atoi(claims.Subject)
			if err != nil {
				err := app.errorJSON(w, errors.New("unknown user"), http.StatusUnauthorized)
				if err != nil {
					return
				}
				return
			}

			// get user from database
			user, err := app.DB.GetUserByID(userID)
			if err != nil {
				err := app.errorJSON(w, errors.New("unknown user"), http.StatusUnauthorized)
				if err != nil {
					return
				}
				return
			}

			// generate token pair
			u := jwtUser{
				ID:        user.ID,
				FirstName: user.FirstName,
				LastName:  user.LastName,
			}

			tokenPairs, err := app.auth.generateTokenPair(&u)
			if err != nil {
				err := app.errorJSON(w, errors.New("error generating token pair"), http.StatusUnauthorized)
				if err != nil {
					return
				}
				return

			}

			refreshCookie := app.auth.getRefreshCookie(tokenPairs.RefreshToken)
			http.SetCookie(w, refreshCookie)

			// write json response
			err = app.writeJSON(w, http.StatusOK, tokenPairs, nil)
			if err != nil {
				return

			}

		}
	}
}

// logout is a simple handler function which writes a response.
func (app *application) logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, app.auth.getExpiredRefreshCookie())
	w.WriteHeader(http.StatusAccepted)
}

// getTokenFromHeaderAndVerify is a simple handler function which writes a response.
func (j *Auth) getTokenFromHeaderAndVerify(w http.ResponseWriter, r *http.Request) (string, *tokenClaims, error) {
	w.Header().Add("Vary", "Authorization")

	// get token from header
	authHeader := r.Header.Get("Authorization")

	if authHeader == "" {
		return "", nil, errors.New("missing authorization header")
	}

	token := strings.Split(authHeader, " ")

	if len(token) != 2 {
		return "", nil, errors.New("invalid authorization header")
	}

	if token[0] != "Bearer" {
		return "", nil, errors.New("invalid authorization header")
	}

	// verify token
	claims := &tokenClaims{}

	_, err := jwt.ParseWithClaims(token[1], claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.Secret), nil
	})
	if err != nil {
		if strings.HasPrefix(err.Error(), "token is expired") {
			return "", nil, errors.New("token is expired")
		}

		return "", nil, err
	}

	if claims.Issuer != j.Issuer {
		return "", nil, errors.New("invalid issuer")
	}

	return token[1], claims, nil

}

// movieCatalog is a simple handler function which writes a response.
func (app *application) movieCatalog(w http.ResponseWriter, r *http.Request) {

	// get all movies from the database
	movies, err := app.DB.AllMovies()
	if err != nil {
		err := app.errorJSON(w, err)
		if err != nil {
			return
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, movies, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
}

// getMovie is a simple handler function which writes a response to retrieve a movie.
func (app *application) getMovie(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	movieID, err := strconv.Atoi(id)
	if err != nil {
		err := app.errorJSON(w, errors.New("invalid id parameter"), http.StatusBadRequest)
		if err != nil {
			return
		}
		return
	}

	movie, err := app.DB.OneMovie(movieID)
	if err != nil {
		err := app.errorJSON(w, err)
		if err != nil {
			return
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, movie, nil)
	if err != nil {
		return
	}

}

// movieForEdit is a simple handler function which writes a response to edit a movie.
func (app *application) movieForEdit(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	movieID, err := strconv.Atoi(id)
	if err != nil {
		err := app.errorJSON(w, errors.New("invalid id parameter"), http.StatusBadRequest)
		if err != nil {
			return
		}
		return
	}

	movie, genres, err := app.DB.OneMovieForEdit(movieID)
	if err != nil {
		err := app.errorJSON(w, err)
		if err != nil {
			return
		}
		return
	}

	var payload = struct {
		Movie  *models.Movie   `json:"movie"`
		Genres []*models.Genre `json:"genres"`
	}{
		Movie:  movie,
		Genres: genres,
	}

	err = app.writeJSON(w, http.StatusOK, payload, nil)
	if err != nil {
		return
	}
}

// allGenres is a simple handler function which writes a response to retrieve all genres.
func (app *application) allGenres(w http.ResponseWriter, r *http.Request) {

	genres, err := app.DB.AllGenresDB()
	if err != nil {
		err := app.errorJSON(w, err)
		if err != nil {
			return
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, genres, nil)
	if err != nil {
		return
	}
}

// insertMovie handler to insert a new movie
func (app *application) insertMovie(w http.ResponseWriter, r *http.Request) {

	var movie models.Movie

	err := app.readJSON(w, r, &movie)
	if err != nil {
		err := app.errorJSON(w, err)
		if err != nil {
			return
		}
		return
	}

	// try to grab an image
	movie = app.getPoster(movie)

	movie.CreatedAt = time.Now()
	movie.UpdatedAt = time.Now()

	newID, err := app.DB.InsertMovie(movie)
	if err != nil {
		err := app.errorJSON(w, err)
		if err != nil {
			return
		}
		return
	}

	// now handle genres
	err = app.DB.UpdateMovieGenres(newID, movie.GenresArray)
	if err != nil {
		err := app.errorJSON(w, err)
		if err != nil {
			return
		}
		return
	}

	response := JSONResponse{
		Error:   false,
		Message: "Movie inserted successfully",
	}

	err = app.writeJSON(w, http.StatusAccepted, response, nil)
	if err != nil {
		return
	}

}

// getPoster gets the poster from the movie db api and returns the movie
func (app *application) getPoster(movie models.Movie) models.Movie {

	type TheMovieDB struct {
		Page    int `json:"page"`
		Results []struct {
			PosterPath string `json:"poster_path"`
		} `json:"results"`
		TotalPages int `json:"total_pages"`
	}

	// get poster from the movie db
	// https://api.themoviedb.org/3/search/movie?api_key=b2225620f919fd84111a706e2dc5d872&query=Die+Hard
	client := &http.Client{}

	theURL := fmt.Sprintf("https://api.themoviedb.org/3/search/movie?api_key=%s", app.APIKey)

	req, err := http.NewRequest("GET", theURL+"&query="+url.QueryEscape(movie.Title), nil)
	if err != nil {
		log.Println("movie not found on the movie db", err)
		return movie
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Println("movie not found on the movie db", err)
		return movie
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println("error closing body", err)
		}
	}(resp.Body)

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("error reading body", err)
		return movie
	}

	var responseMovie TheMovieDB

	err = json.Unmarshal(bodyBytes, &responseMovie)
	if err != nil {
		log.Println("error unmarshalling json", err)
		return movie
	}

	if len(responseMovie.Results) > 0 {
		movie.Image = responseMovie.Results[0].PosterPath
	}

	return movie
}

// updateMovie handler to update a movie
func (app *application) updateMovie(w http.ResponseWriter, r *http.Request) {

	var payload models.Movie

	err := app.readJSON(w, r, &payload)
	if err != nil {
		err := app.errorJSON(w, err)
		if err != nil {
			return
		}
		return
	}

	movie, err := app.DB.OneMovie(payload.ID)
	if err != nil {
		err := app.errorJSON(w, err)
		if err != nil {
			return
		}
		return
	}

	movie.Title = payload.Title
	movie.ReleaseDate = payload.ReleaseDate
	movie.Description = payload.Description
	movie.Runtime = payload.Runtime
	movie.MPAARating = payload.MPAARating
	movie.UpdatedAt = time.Now()

	err = app.DB.UpdateMovie(*movie)
	if err != nil {
		err := app.errorJSON(w, err)
		if err != nil {
			return
		}
		return
	}

	err = app.DB.UpdateMovieGenres(movie.ID, payload.GenresArray)
	if err != nil {
		err := app.errorJSON(w, err)
		if err != nil {
			return
		}
		return
	}

	response := JSONResponse{
		Error:   false,
		Message: "movie updated successfully",
	}

	err = app.writeJSON(w, http.StatusAccepted, response, nil)
	if err != nil {
		return
	}

}

// deleteMovie handler to delete a movie
func (app *application) deleteMovie(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		err := app.errorJSON(w, err)
		if err != nil {
			return
		}
		return
	}

	err = app.DB.DeleteMovie(id)
	if err != nil {
		err := app.errorJSON(w, err)
		if err != nil {
			return
		}
		return
	}

	response := JSONResponse{
		Error:   false,
		Message: "movie deleted successfully",
	}

	err = app.writeJSON(w, http.StatusAccepted, response, nil)
	if err != nil {
		return
	}

}

// AllMoviesByGenre handler to get all movies by genre
func (app *application) AllMoviesByGenre(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		err := app.errorJSON(w, err)
		if err != nil {
			return
		}
		return
	}

	movies, err := app.DB.AllMovies(id)
	if err != nil {
		err := app.errorJSON(w, err)
		if err != nil {
			return
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, movies, nil)
	if err != nil {
		return
	}

}
