package main

import "net/http"

// showAllQuizzesHandler sends all quizzes in the database in a JSON response to the
// client
func (app *application) showAllQuizzesHandler(w http.ResponseWriter, r *http.Request) {
    app.writeJSON(w, r, http.StatusOK, "Showing all quizzes...", nil)
}


