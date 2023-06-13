package main

import "net/http"

// logError prints the error with app.logger
func (app *application) logError(r *http.Request, err error) {
	app.logger.Println(err)
}

// errorResponse writes an error response as JSON with the supplied message and code
func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, code int, message interface{}) {
	if err := app.writeJSON(w, r, code, message, nil); err != nil {
		app.logError(r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// notFoundResponse writes a Not Found reponse
func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "The server could not find the requested resource"

	app.errorResponse(w, r, http.StatusNotFound, message)
}

// serverErrorResponse writes an Internal Server Error response and logs err to app
func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	message := "The server encountered an error processing the request"

	app.logError(r, err)

	app.errorResponse(w, r, http.StatusInternalServerError, message)
}

// methodNotAllowedResponse writes a Method Not Allowed error response
func (app *application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := "The method " + r.Method + " is not supported for this resource"

	app.errorResponse(w, r, http.StatusMethodNotAllowed, message)
}