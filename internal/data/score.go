package data

import (
	"context"
	"database/sql"
	"errors"
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

var ErrInvalidAnswers = errors.New("Answers must be of the same length as questions")

func (score *Score) CalcScore(questions []*Question, answers []int32) error {
    if len(answers) != len(questions) {
        return ErrInvalidAnswers
    }

    score.ChosenChoicesIndices = answers
    score.ChosenChoicesCorrectness = []bool{}

    for i, answer := range answers {
        correct := false
        if answer == questions[i].CorrectIndex {
            correct = true
        }

        score.ChosenChoicesCorrectness = append(score.ChosenChoicesCorrectness, correct)
    }

    return nil
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

func (m ScoreModel) GetScoresByUserAndQuiz(userID, quizID int64) ([]*Score, error) {
    stmt := `SELECT id, user_id, quiz_id, chosen_choices_indices,
        chosen_choices_correctness, created_at, version
        FROM score
        WHERE user_id = $1 AND quiz_id = $2`

    ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
    defer cancel()

    rows, err := m.DB.QueryContext(ctx, stmt, userID, quizID);
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    scores := []*Score{}

    for rows.Next() {
        score := &Score{}
        err := rows.Scan(&score.ID, &score.UserID, &score.QuizID, &score.ChosenChoicesIndices,
            &score.ChosenChoicesCorrectness, &score.CreatedAt, &score.Version)
        if err != nil {
            return nil, err
        }
        scores = append(scores, score)
    }

    if rows.Err() != nil {
        return nil, rows.Err()
    }

    if len(scores) == 0 {
        return nil, ErrNoRecords
    }


    return scores, nil
}
