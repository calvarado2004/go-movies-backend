package main

import (
	"fmt"
	"net/http"
)

// Home is a simple handler function which writes a response.
func (app *application) Home(w http.ResponseWriter, r *http.Request) {

	_, err := fmt.Fprintf(w, "Hello, World! from %s", app.Domain)
	if err != nil {
		return
	}
}
