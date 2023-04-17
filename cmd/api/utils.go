package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

// JSONResponse is a struct that is used to return a JSON response
type JSONResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// writeJSON is a helper function that writes JSON data to the response body.
func (app *application) writeJSON(w http.ResponseWriter, status int, data any, headers ...http.Header) error {

	out, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(out)
	if err != nil {
		return err
	}

	return nil

}

// readJSON is a helper function that reads JSON data from the request body.
func (app *application) readJSON(w http.ResponseWriter, r *http.Request, data any) error {

	maxBytes := 1024 * 1024 // 1MB

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)

	dec.DisallowUnknownFields() // this will return an error if the request body contains any fields that aren't defined in the data struct

	err := dec.Decode(data)
	if err != nil {
		return err
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain a single JSON object")
	}

	return nil
}
