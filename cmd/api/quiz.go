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
		Quiz      *data.Quiz
		Questions []*data.Question
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
	// Create a struct to hold the quiz data
	var quiz struct {
		Title string `json:"title"`
	}

	// Read the json request body into the quiz struct
	err := app.readJSON(w, r, &quiz)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	if err := validator.ValidateQuiz(data.Quiz{Title: quiz.Title}); err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Add the quiz to the database
	title, err := app.models.Quizzes.Add(quiz.Title)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Send the id of the quiz in the response
	app.writeJSON(w, r, http.StatusCreated, title, nil)

}

// addQuestionHandler receives a json response of a question and
// adds the question to the database, as well as responding to the
// client with the added question.
func (app *application) addQuestionHandler(w http.ResponseWriter, r *http.Request) {
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

	if err := validator.ValidateQuestion(
        data.Question{Prompt: input.Prompt, Choices: input.Choices, CorrectIndex: input.CorrectIndex}); 
        err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	question := data.Question{
		Prompt:       input.Prompt,
		Choices:      input.Choices,
		QuizID:       quizID,
		CorrectIndex: input.CorrectIndex,
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

// addScoreHandler receives a json response containing the user's
// quiz answers and returns the user's score on the quiz.
func (app *application) addScoreHandler(w http.ResponseWriter, r *http.Request) {
    params := httprouter.ParamsFromContext(r.Context())
    id := params.ByName("id")
    var input struct {
        Answers []int32
    }

    err := app.readJSON(w, r, &input)
    if err != nil {
        app.errorResponse(w, r, http.StatusBadRequest, err)
        return
    }

    questions, err := app.models.Questions.GetAllByQuizID(id)
    if err != nil {
        app.serverErrorResponse(w, r, err)
        return
    }

    correct := 0
    for i, question := range questions {
        if input.Answers[i] == question.CorrectIndex {
            correct += 1
        }
    }

    score := float32(correct) / float32(len(questions)) * 100.0

    err = app.writeJSON(w, r, http.StatusOK, score, nil)
    if err != nil {
        app.serverErrorResponse(w, r, err)
    }
}
