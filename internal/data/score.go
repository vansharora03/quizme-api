package data

import (
	"context"
	"database/sql"
	"time"

	"github.com/lib/pq"
)

type ScoreModel struct {
    DB *sql.DB
}

type Score struct {
    ID int64 `json:"id"`
    UserID int64 `json:"user_id"`
    QuizID int64 `json:"quiz_id"`
    ChosenChoicesIndices []int32 `json:"chosen_choices_indices"`
    ChosenChoicesCorrectness []bool `json:"chosen_choices_correctness"`
    CreatedAt time.Time `json:"created_at"`
    Version int32 `json:"version"`
}

func (m ScoreModel) AddScore(score *Score) error {
    stmt := `INSERT INTO score 
    (user_id, quiz_id, chosen_choices_indices, chosen_choices_correctness)
    VALUES ($1, $2, $3, $4)
    RETURNING id, created_at, version`

    ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
    defer cancel()

    err := m.DB.QueryRowContext(ctx, stmt,
        score.UserID, 
        score.QuizID, 
        pq.Array(score.ChosenChoicesIndices), 
        pq.Array(score.ChosenChoicesCorrectness)).Scan( &score.ID, &score.CreatedAt, &score.Version)

    return err
}
