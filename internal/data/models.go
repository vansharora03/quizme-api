package data

import "database/sql"

type Models struct {
	Questions interface{}
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
