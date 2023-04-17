package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/calvarado2004/go-movies-backend/internal/repository"
	"github.com/calvarado2004/go-movies-backend/internal/repository/dbrepo"
	"log"
	"net/http"
)

const port = 8080

type application struct {
	Domain string
	DSN    string
	DB     repository.DatabaseRepo
}

func main() {

	// set application config
	var app application

	// read from command line
	flag.StringVar(&app.DSN, "dsn", "host=localhost port=5432 user=postgres password=postgres dbname=movies sslmode=disable timezone=UTC connect_timeout=5", "PostgreSQL DSN")

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

	app.Domain = "example.com"

	log.Println(fmt.Sprintf("Starting server on port %d", port))

	// start a web server
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), app.routes())
	if err != nil {
		log.Fatal(err)
	}

}
