package data

import "database/sql"

type Models struct {
    questions  QuestionModel
    quizzes    QuizModel
}

// NewModels initializes a Models struct with
// all of quizme's data models connected on db
func NewModels(db *sql.DB) Models {
    return Models{
        QuestionModel{db},
        QuizModel{db},
    }
}
