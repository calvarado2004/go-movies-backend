package main

import (
	"fmt"
	"net/http"
)

// Hello is a simple handler function which writes a response.
func Hello(w http.ResponseWriter, r *http.Request) {

	_, err := fmt.Fprint(w, "Hello, World!")
	if err != nil {
		return
	}
}
