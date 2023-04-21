package main

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"strconv"
	"strings"
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
