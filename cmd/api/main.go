package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/calvarado2004/go-movies-backend/internal/repository"
	"github.com/calvarado2004/go-movies-backend/internal/repository/dbrepo"
	"log"
	"net/http"
	"os"
	"time"
)

const port = 8080

type application struct {
	Domain       string
	DSN          string
	DB           repository.DatabaseRepo
	auth         Auth
	JWTSecret    string
	JWTIssuer    string
	JWTAudience  string
	CookieDomain string
	APIKey       string
}

func main() {

	// set application config
	var app application

	dbServer := os.Getenv("DB_SERVER")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")

	dsnVariable := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable timezone=UTC connect_timeout=5", dbServer, dbPort, dbUser, dbPassword, dbName)

	apiMoviesKey := os.Getenv("API_MOVIES_KEY")

	// read from command line
	flag.StringVar(&app.DSN, "dsn", dsnVariable, "PostgreSQL DSN")
	flag.StringVar(&app.JWTSecret, "jwt-secret", "verysecret", "JWT Secret")
	flag.StringVar(&app.JWTIssuer, "jwt-issuer", "api-golang-movies.apps.okd.calvarado04.com", "JWT Issuer")
	flag.StringVar(&app.JWTAudience, "jwt-audience", "node-react-movies.apps.okd.calvarado04.com", "JWT Audience")
	flag.StringVar(&app.CookieDomain, "cookie-domain", "apps.okd.calvarado04.com", "Cookie Domain")
	flag.StringVar(&app.Domain, "domain", "apps.okd.calvarado04.com", "Domain")
	flag.StringVar(&app.APIKey, "api-key", apiMoviesKey, "API Key")

	flag.Parse()

	// connect to the database
	conn, err := app.connectToDB()
	if err != nil {
		log.Fatal(err)
	}

	app.DB = &dbrepo.PostgresDBRepo{DB: conn}

	defer func(connection *sql.DB) {
		err := connection.Close()
		if err != nil {
			return
		}
	}(app.DB.Connection())

	app.auth = Auth{
		Issuer:        app.JWTIssuer,
		Audience:      app.JWTAudience,
		Secret:        app.JWTSecret,
		TokenExpiry:   15 * time.Minute,
		RefreshExpiry: 24 * time.Hour,
		CookiePath:    "/",
		CookieName:    "jwt-refresh_token",
		CookieDomain:  app.CookieDomain,
	}

	log.Println(fmt.Sprintf("Starting server on port %d", port))

	// start a web server
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), app.routes())
	if err != nil {
		log.Fatal(err)
	}

}
