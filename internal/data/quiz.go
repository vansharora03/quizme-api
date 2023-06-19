package data

import (
	"context"
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
    Title string // Title of the quiz
}


// GetAll will get all quizzes in the database and return them.
func (m QuizModel) GetAll() ([]*Quiz, error) {
    stmt := `SELECT id, created_at, version, title FROM quiz`
    ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
    defer cancel()

    rows, err := m.DB.QueryContext(ctx, stmt)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    quizzes := []*Quiz{}

    for rows.Next() {
        q := &Quiz{}
        err := rows.Scan(&q.ID, &q.CreatedAt, &q.Version, &q.Title)
        if err != nil {
            return nil, err
        }

        quizzes = append(quizzes, q)
    }

    err = rows.Err()

    if err != nil {
        return nil, err
    }

    return quizzes, nil
}
