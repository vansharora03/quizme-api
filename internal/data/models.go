package data

import "database/sql"

type Models struct {
    Questions  QuestionModel
    Quizzes    QuizModel
}

// NewModels initializes a Models struct with
// all of quizme's data models connected on db
func NewModels(db *sql.DB) Models {
    return Models{
        Questions: QuestionModel{db},
        Quizzes: QuizModel{db},
    }
}
