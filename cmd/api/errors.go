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

// clientError sends a specific status code to the client and writes the status text
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
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

// validationErrorResponse writes a Bad Request response along with a map of errors,
// to be used on failed validation
func (app *application) validationErrorResponse(w http.ResponseWriter, r *http.Request, 
    errors map[string]string) {
        app.errorResponse(w, r, http.StatusBadRequest, errors)
}

// forbiddenResponse writes a Forbidden response to the user, to be used when
// a user must be logged in to use a resource
func (app *application) forbiddenResponse(w http.ResponseWriter, r *http.Request) {
        app.errorResponse(w, r, http.StatusForbidden, "Please log in to use this route")
}
