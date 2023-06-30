package data

import (
	"database/sql"
	"errors"
)

var ErrNoRecords = errors.New("No records found")
var ErrEditConflict = errors.New("Edit conflict")

type Models struct {
	Questions interface {
		GetAllByQuizID(quizID string) ([]*Question, error)
		AddQuestion(question *Question) error
        Update(question *Question) error
	}

	Quizzes interface {
		GetAll() ([]*Quiz, error)
		Get(id string) (*Quiz, error)
		Add(title string) (string, error)
        Update(quiz *Quiz) error
	}
}

// NewModels initializes a Models struct with
// all of quizme's data models connected on db
func NewModels(db *sql.DB) Models {
	return Models{
		Questions: QuestionModel{db},
		Quizzes:   QuizModel{db},
	}
}
