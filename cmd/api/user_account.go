package main

import (
	"errors"
	"net/http"
	"time"
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
    
    err = data.HashUser(&user)
    if err != nil {
        app.serverErrorResponse(w, r, err)
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

// userLoginHandler reads user email and password from input, if these are valid,
// return a web token to the user and add that token's hash to the database for the user.
func (app *application) userLoginHandler(w http.ResponseWriter, r *http.Request) {
    var input struct {
        Email string `json:"email"`
        Password string `json:"password"`
    }

    err := app.readJSON(w, r, &input)
    if err != nil {
        app.errorResponse(w, r, http.StatusBadRequest, err)
        return
    }

    v := validator.New()

    data.ValidateEmail(&v, input.Email)
    data.ValidatePassword(&v, input.Password)

    if !v.Valid() {
        app.validationErrorResponse(w, r, v.Errors)
        return
    }

    invalidCredentialsMessage := "Email or password is invalid"

    user, err := app.models.Users.GetUserByEmail(input.Email)
    if err != nil {
        switch {
        case errors.Is(err, data.ErrNoRecords):
            v.Add("invalid credentials", invalidCredentialsMessage)
            app.validationErrorResponse(w, r, v.Errors)
        default:
            app.serverErrorResponse(w, r, err)
        }
        return
    }

    err = user.MatchPassword(input.Password)
    if err != nil {
        v.Add("invalid credentials", invalidCredentialsMessage)
        app.validationErrorResponse(w, r, v.Errors)
        return
    }

    token, err := app.models.Tokens.AddToken(user.ID, 24 * time.Hour)
    if err != nil {
        app.serverErrorResponse(w, r, err)
        return
    }

    err = app.writeJSON(w, r, http.StatusCreated, &token, nil)
    if err != nil {
        app.serverErrorResponse(w, r, err)
    }
}
