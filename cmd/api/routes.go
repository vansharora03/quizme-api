package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {

	// Middleware chain
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	router := httprouter.New()

	// Router Settings
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	// Routes for application
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodGet, "/v1/quiz", app.showAllQuizzesHandler)
    router.HandlerFunc(http.MethodGet, "/v1/quiz/:id", app.showQuizHandler)
    router.HandlerFunc(http.MethodPost, "/v1/quiz", app.addQuizHandler)

	return standardMiddleware.Then(router)
}
