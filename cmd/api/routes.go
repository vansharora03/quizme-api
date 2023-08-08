package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {

	// Middleware chain
	standardMiddleware := alice.New(app.recoverPanic, app.rateLimit, app.logRequest, secureHeaders)

	dynamicMiddleware := alice.New(app.authenticate)

	router := httprouter.New()

	// Router Settings
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	// Routes for application
	router.Handler(http.MethodGet, "/v1/healthcheck", dynamicMiddleware.ThenFunc(app.healthcheckHandler))
	router.Handler(http.MethodGet, "/v1/quiz", dynamicMiddleware.ThenFunc(app.showAllQuizzesHandler))
    router.Handler(http.MethodGet, "/v1/quiz/:id", dynamicMiddleware.ThenFunc(app.showQuizHandler))
    router.Handler(http.MethodGet, "/v1/quiz/:id/score", dynamicMiddleware.ThenFunc(app.showScoresHandler))
    router.Handler(http.MethodPost, "/v1/quiz", dynamicMiddleware.ThenFunc(app.addQuizHandler))
    router.Handler(http.MethodPost, "/v1/quiz/:id/score", dynamicMiddleware.ThenFunc(app.addScoreHandler))
    router.Handler(http.MethodPost, "/v1/quiz/:id/question", dynamicMiddleware.ThenFunc(app.addQuestionHandler))
    router.Handler(http.MethodPut, "/v1/quiz/:id", dynamicMiddleware.ThenFunc(app.updateQuizHandler))
    router.Handler(http.MethodPut, "/v1/quiz/:id/question/:questionID", dynamicMiddleware.ThenFunc(app.updateQuestionHandler))

    router.Handler(http.MethodPost, "/v1/user", dynamicMiddleware.ThenFunc(app.addUserHandler))
    router.Handler(http.MethodPost, "/v1/user/login", dynamicMiddleware.ThenFunc(app.userLoginHandler))


	return standardMiddleware.Then(router)
}
