package main

import (
    "net/http"
    "fmt"
    "vanshadhruvp/quizme-api/internal/data"
)

// quizBelongsToUser returns whether the user is the owner of the quiz, if this returns false, the function
// will handle the error response handling.
func (app *application) userMatchesQuiz(w http.ResponseWriter, r *http.Request, userID, quizID int64) bool {
    quiz, err := app.models.Quizzes.Get(fmt.Sprintf("%d", quizID))
    if err == data.ErrNoRecords {
        app.notFoundResponse(w, r)
        return false;
    } else if err != nil {
        app.serverErrorResponse(w, r, err)
        return false;
    } else if quiz.UserID != userID {
        app.forbiddenResponse(w, r)
        return false;
    }
    return true;
}
