package main

import "net/http"

func (app *application) showAllQuizzes(w http.ResponseWriter, r *http.Request) {
    app.writeJSON(w, r, http.StatusOK, "Showing all quizzes...", nil)
}
