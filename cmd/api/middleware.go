package main

import "net/http"

// enableCORS is a middleware function that adds the appropriate CORS headers to the response.
func (app *application) enableCORS(h http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "https://node-react-movies.apps.okd.calvarado04.com")
		w.Header().Set("Access-Control-Allow-Origin", "https://node-react-svc:3000")
		w.Header().Set("Access-Control-Allow-Origin", "https://localhost:3000")

		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

			return
		}
		h.ServeHTTP(w, r)
	})
}

// authRequired is a middleware function that checks that the request contains a valid JWT token.
func (app *application) authRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _, err := app.auth.getTokenFromHeaderAndVerify(w, r)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
