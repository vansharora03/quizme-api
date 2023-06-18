package data

import (
	"database/sql"
	"time"
)

type QuizModel struct {
    DB *sql.DB
}

type Quiz struct {
    ID int64 // ID
    CreatedAt time.Time // When the quiz was created
    Version int32 // Version of the quiz
}
