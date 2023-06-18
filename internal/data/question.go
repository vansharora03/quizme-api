package data

import (
	"database/sql"
	"time"
)

type QuestionModel struct {
    DB *sql.DB
}

type Question struct {
    ID int64 // ID
    QuizID int64 // The id of the quiz that this question belongs to
    Prompt string // The prompt that is to be answered
    Choices []string // The possible answer choices
    CorrectAnswer string // The correct answer
    CreatedAt time.Time // When the question was created
    Version int32 // The version of the question
}
