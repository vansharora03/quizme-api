package main

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() *httprouter.Router {
	router := httprouter.New()

	// Router Settings
	// TODO: Create custom error handling for http errors outside defined handlers->
	// router.NotFound = ?
	// router.MethodNotAllowed = ?

	// Routes for application
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Environment: %s\nVersion: %s", app.config.env, version)
	})

	return router
}
