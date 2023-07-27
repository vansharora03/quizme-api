package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {

	// Middleware chain
	standardMiddleware := alice.New(app.recoverPanic, app.rateLimit, app.logRequest, secureHeaders)

	// dynamic middleware chain for future
	//dynamicMiddleware := alice.New(nosurf)

	router := httprouter.New()

	// Router Settings
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	// Routes for application
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodGet, "/v1/quiz", app.showAllQuizzesHandler)
    router.HandlerFunc(http.MethodGet, "/v1/quiz/:id", app.showQuizHandler)
    router.HandlerFunc(http.MethodPost, "/v1/quiz", app.addQuizHandler)
    router.HandlerFunc(http.MethodPost, "/v1/quiz/:id/score", app.addScoreHandler)
    router.HandlerFunc(http.MethodPost, "/v1/quiz/:id/question", app.addQuestionHandler)
    router.HandlerFunc(http.MethodPut, "/v1/quiz/:id", app.updateQuizHandler)
    router.HandlerFunc(http.MethodPut, "/v1/quiz/:id/question/:questionID", app.updateQuestionHandler)

    router.HandlerFunc(http.MethodPost, "/v1/user", app.addUserHandler)
    router.HandlerFunc(http.MethodPost, "/v1/user/login", app.userLoginHandler)


	return standardMiddleware.Then(router)
}
