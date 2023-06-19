package main

import "net/http"
import _ "vanshadhruvp/quizme-api/internal/data"

// showAllQuizzesHandler sends all quizzes in the database in a JSON response to the
// client
func (app *application) showAllQuizzesHandler(w http.ResponseWriter, r *http.Request) {
    quizzes, err := app.models.Quizzes.GetAll()
    if err != nil {
        app.serverErrorResponse(w, r, err)
        return
    }

    app.writeJSON(w, r, http.StatusOK, quizzes, nil)
}

// showQuizHandler sends a specific quiz in the database in a JSON response to the
// client
func (app *application) showQuizHandler(w http.ResponseWriter, r *http.Request) {
    app.writeJSON(w, r, http.StatusOK, "Showing specific quiz...", nil)
}

// addQuizHandler adds a specific quiz to the database
func (app *application) addQuizHandler(w http.ResponseWriter, r *http.Request) {
    app.writeJSON(w, r, http.StatusOK, "Adding a quiz...", nil)
}

// addScoreHandler receives a json response containing the user's
// quiz answers and returns the user's score on the quiz.
func (app *application) addScoreHandler(w http.ResponseWriter, r *http.Request) {
    app.writeJSON(w, r, http.StatusOK, "Adding a score...", nil)
}


