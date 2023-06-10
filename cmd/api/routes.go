package main

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {

	// Middleware chain
	standardMiddleware := alice.New(app.logRequest, secureHeaders)

	router := httprouter.New()

	// Router Settings
	// TODO: Create custom error handling for http errors outside defined handlers->
	// router.NotFound = ?
	// router.MethodNotAllowed = ?

	// Routes for application
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Environment: %s\nVersion: %s", app.config.env, version)
	})

	return standardMiddleware.Then(router)
}
