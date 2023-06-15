package main

import (
	"encoding/json"
	"net/http"
)

// writeJSON writes the data as json and sends it in the response body with the status code and headers. The response is guaranteed
// not to be written until errors have been handled.
func (app *application) writeJSON(w http.ResponseWriter, r *http.Request, code int, data interface{}, headers http.Header) error {
	body, err := json.Marshal(data)
	if err != nil {
		return err
	}

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(body)
	return nil
}

// readJSON reads the json data from the request body and scans the values into the destination of the pointer dst.
func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	defer r.Body.Close()

	err := dec.Decode(dst)
	if err != nil {
		return err
	}

	return nil
}
