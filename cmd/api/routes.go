package main

import (
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
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	return standardMiddleware.Then(router)
}
