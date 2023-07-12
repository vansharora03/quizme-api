package main

import (
	"errors"
	"net/http"
	"vanshadhruvp/quizme-api/internal/data"
	"vanshadhruvp/quizme-api/internal/validator"
)

// addUserHandler reads json user data from the request body and adds a user to the
// database
func (app *application) addUserHandler(w http.ResponseWriter, r *http.Request) {
    var input struct {
        Email string `json:"email"`
        Username string `json:"username"`
        Password string `json:"password"`
    }

    err := app.readJSON(w, r, &input)
    if err != nil {
        app.errorResponse(w, r, http.StatusBadRequest, err)
        return
    }

    user := data.User{
        Email: input.Email,
        Username: input.Username,
        PlaintextPassword: input.Password,
    }
    
    v := validator.New()

    if data.ValidateUser(&v, &user); !v.Valid() {
        app.validationErrorResponse(w, r, v.Errors)
        return
    }

    err = app.models.Users.AddUser(&user)
    if err != nil {
        switch {
        case errors.Is(err, data.ErrDuplicateEmail):
            v.Add("email", "email already exists")
            app.validationErrorResponse(w, r, v.Errors)
        default:
            app.serverErrorResponse(w, r, err)
        }

        return
    }

    err = app.writeJSON(w, r, http.StatusCreated, &user, nil)
    if err != nil {
        app.serverErrorResponse(w, r, err)
    }
}
