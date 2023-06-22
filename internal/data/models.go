package data

import (
	"database/sql"
	"errors"
)

var ErrNoRecords = errors.New("No records found")

type Models struct {
	Questions interface{
        GetAllByQuizID(quizID string) ([]*Question, error)
        AddQuestion(question *Question) error 
    }

	Quizzes   interface {
		GetAll() ([]*Quiz, error)
		Get(id string) (*Quiz, error)
		Add(title string) (int64, error)
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
