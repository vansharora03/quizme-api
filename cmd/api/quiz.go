package main

import (
	"net/http"
	"strconv"
	"vanshadhruvp/quizme-api/internal/data"
	_ "vanshadhruvp/quizme-api/internal/data"

	"vanshadhruvp/quizme-api/internal/validator"

	"github.com/julienschmidt/httprouter"
)

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

	// Get the id of the quiz from the url
	params := httprouter.ParamsFromContext(r.Context())
	quizID := params.ByName("id")

	var quizInstance struct {
        Quiz      *data.Quiz `json:"quiz"`
        Questions []*data.Question `json:"questions"`
	}

	// Get the quiz from the database
	quiz, err := app.models.Quizzes.Get(quizID)
	if err == data.ErrNoRecords {
		app.notFoundResponse(w, r)
		return
	} else if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	quizInstance.Quiz = quiz

	// Get the questions from the database
	questions, err := app.models.Questions.GetAllByQuizID(quizID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	quizInstance.Questions = questions

	app.writeJSON(w, r, http.StatusOK, quizInstance, nil)
}

// addQuizHandler adds a specific quiz to the database
func (app *application) addQuizHandler(w http.ResponseWriter, r *http.Request) {
    userVal := (r.Context().Value("user"))
    user := userVal.(*data.User)
    if user == data.AnonymousUser {
        app.forbiddenResponse(w, r)
        return
    }

	// Create a struct to hold the quiz data
    var input struct {
        Title string `json:"title"`
    }

	// Read the json request body into the quiz struct
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
    
	var quiz data.Quiz 
    quiz.Title = input.Title

    v := validator.New()

    if data.ValidateQuiz(&v, &quiz); !v.Valid() {
        app.validationErrorResponse(w, r, v.Errors)
        return
    }

	// Add the quiz to the database
    q, err := app.models.Quizzes.Add(quiz.Title, user.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Send the id of the quiz in the response
	app.writeJSON(w, r, http.StatusCreated, q, nil)

}

// addQuestionHandler receives a json response of a question and
// adds the question to the database, as well as responding to the
// client with the added question.
func (app *application) addQuestionHandler(w http.ResponseWriter, r *http.Request) {
    userVal := r.Context().Value("user")
    user := userVal.(*data.User)
    if user == data.AnonymousUser {
        app.forbiddenResponse(w, r)
        return
    }
	// Get the id of the quiz from the url
	params := httprouter.ParamsFromContext(r.Context())
	stringID := params.ByName("id")
	quizID, err := strconv.ParseInt(stringID, 10, 64)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	var input struct {
		Prompt       string   `json:"prompt"`
		Choices      []string `json:"choices"`
		CorrectIndex int32    `json:"correct_index"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, err)
		return
	}

	question := data.Question{
		Prompt:       input.Prompt,
		Choices:      input.Choices,
		QuizID:       quizID,
		CorrectIndex: input.CorrectIndex,
	}

    v := validator.New()

    if data.ValidateQuestion(&v, &question); !v.Valid() {
        app.validationErrorResponse(w, r, v.Errors)
        return
    }

    if !app.userMatchesQuiz(w, r, user.ID, quizID) {
        return
    }

	err = app.models.Questions.AddQuestion(&question)
	if err == data.ErrNoRecords {
		app.notFoundResponse(w, r)
		return
	} else if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, r, http.StatusCreated, question, nil)
}

// updateQuizHandler receives a quiz from the json request and 
// updates the corresponding quiz in the database
func (app *application) updateQuizHandler(w http.ResponseWriter, r *http.Request) {
    userVal := r.Context().Value("user");
    user := userVal.(*data.User)
    if user == data.AnonymousUser {
        app.forbiddenResponse(w, r)
        return
    }

    params := httprouter.ParamsFromContext(r.Context())
    id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
    if err != nil {
        app.notFoundResponse(w, r)
        return
    }

    var input struct {
        Title string `json:"title"`
        Version int32 `json:"version"`
    }

    err = app.readJSON(w, r, &input)
    if err != nil {
        app.errorResponse(w, r, http.StatusBadRequest, err)
        return
    }

    var quiz data.Quiz
    quiz.Title = input.Title
    quiz.Version = input.Version
    quiz.ID = id

    v := validator.New()

    if data.ValidateQuiz(&v, &quiz); !v.Valid() {
        app.validationErrorResponse(w, r, v.Errors)
        return
    }

    if !app.userMatchesQuiz(w, r, user.ID, quiz.ID) {
        return
    }

    err = app.models.Quizzes.Update(&quiz)
    if err != nil {
        switch {
        case err == data.ErrEditConflict:
            app.errorResponse(w, r, http.StatusConflict, "Please try again")
            return
        default:
            app.serverErrorResponse(w, r, err)
            return
        }
    }

    app.writeJSON(w, r, http.StatusOK, quiz, nil)
}
        

// updateQuestionHandler receives an updated question from the user and updates the
// question in the database
func (app *application) updateQuestionHandler(w http.ResponseWriter, r *http.Request) {
    userVal := r.Context().Value("user");
    user := userVal.(*data.User)
    if user == data.AnonymousUser {
        app.forbiddenResponse(w, r)
        return
    }
    

    params := httprouter.ParamsFromContext(r.Context())
    quizID, err := strconv.ParseInt(params.ByName("id"), 10, 64)
    questionID, err := strconv.ParseInt(params.ByName("questionID"), 10, 64)
    if err != nil {
        app.notFoundResponse(w, r)
        return
    }

    var input struct {
        Prompt string `json:"prompt"`
        Choices []string `json:"choices"`
        CorrectIndex int32 `json:"correct_index"`
        Version int32 `json:"version"`
    }

    err = app.readJSON(w, r, &input)
    if err != nil {
        app.errorResponse(w, r, http.StatusBadRequest, err.Error())
        return
    }

    var question data.Question
    question.ID = questionID
    question.QuizID = quizID
    question.Prompt = input.Prompt
    question.Choices = input.Choices
    question.CorrectIndex = input.CorrectIndex
    question.Version = input.Version

    v := validator.New()

    if data.ValidateQuestion(&v, &question); !v.Valid() {
        app.validationErrorResponse(w, r, v.Errors)
        return
    }

    if !app.userMatchesQuiz(w, r, user.ID, quizID) {
        return
    }

    err = app.models.Questions.Update(&question)
    if err != nil {
        switch {
        case err == data.ErrEditConflict:
            app.errorResponse(w, r, http.StatusConflict, "Please try again")
            return
        default:
            app.serverErrorResponse(w, r, err)
            return
        }
    }

    app.writeJSON(w, r, http.StatusOK, question, nil)
}

// addScoreHandler takes a user's answers and scores it, sending the score to the database.
func (app *application) addScoreHandler(w http.ResponseWriter, r *http.Request) {
    userVal := r.Context().Value("user")
    user := userVal.(*data.User)
    if user == data.AnonymousUser {
        app.forbiddenResponse(w, r)
        return
    }

    params := httprouter.ParamsFromContext(r.Context())
    quizIDRaw := params.ByName("id")
    quizID, err := strconv.ParseInt(quizIDRaw, 10, 64)
    if err != nil || quizID < 1 {
        app.notFoundResponse(w, r)
        return
    }

    var input struct {
        answers []int32
    }

    err = app.readJSON(w, r, &input)
    if err != nil {
        app.errorResponse(w, r, http.StatusBadRequest, err.Error())
        return
    }

    score := &data.Score{
        QuizID: quizID,
        UserID: user.ID,
    }

    questions, err := app.models.Questions.GetAllByQuizID(quizIDRaw)
    if err != nil {
        switch {
        case err == data.ErrNoRecords:
            app.notFoundResponse(w, r)
            return
        default:
            app.serverErrorResponse(w, r, err)
            return
        }
    }

    err = score.CalcScore(questions, input.answers)
    if err != nil {
        app.errorResponse(w, r, http.StatusBadRequest, err.Error())
        return
    }

    err = app.models.Scores.AddScore(score)
    if err != nil {
        app.serverErrorResponse(w, r, err)
        return
    }

    err = app.writeJSON(w, r, http.StatusCreated, score, nil)
    if err != nil {
        app.serverErrorResponse(w, r, err)
    }

}

// showScoresHandler shows all of the user's scores on the chosen quiz
func (app *application) showScoresHandler(w http.ResponseWriter, r *http.Request) {
    userVal := r.Context().Value("user")
    user := userVal.(*data.User)
    if user == data.AnonymousUser {
        app.forbiddenResponse(w, r)
        return
    }

    params := httprouter.ParamsFromContext(r.Context())
    quizID, err := strconv.ParseInt(params.ByName("id"), 10, 64)
    if err != nil || quizID < 1 {
        app.notFoundResponse(w, r)
        return
    }

    scores, err := app.models.Scores.GetScoresByUserAndQuiz(user.ID, quizID)
    if err != nil {
        switch {
        case err == data.ErrNoRecords:
            app.notFoundResponse(w, r)
            return
        default:
            app.serverErrorResponse(w, r, err)
            return
        }
    }

    err = app.writeJSON(w, r, http.StatusOK, scores, nil)
    if err != nil {
        app.serverErrorResponse(w, r, err)
    }
}

















