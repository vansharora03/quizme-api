package data

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/lib/pq"
)

type QuestionModel struct {
    DB *sql.DB
}

type Question struct {
    ID int64 `json:"id"` // ID
    QuizID int64 `json:"quiz_id"` // The id of the quiz that this question belongs to
    Prompt string  `json:"prompt"` /// The prompt that is to be answered
    Choices []string  `json:"choices"` /// The possible answer choices
    CorrectIndex int32  `json:"correct_index"` /// The index at which the correct answer is found
    CreatedAt time.Time  `json:"created_at"` /// When the question was created
    Version int32  `json:"version"` /// The version of the question
}

// GetAllByQuizID returns all questions associated with the quiz with quizID
func (m QuestionModel) GetAllByQuizID(quizID string) ([]*Question, error) {
    stmt := `SELECT question.* 
    FROM quiz 
    INNER JOIN question ON question.quiz_id = quiz.id
    WHERE quiz.id = $1`

    ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
    defer cancel()

    rows, err := m.DB.QueryContext(ctx, stmt, quizID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    questions := []*Question{}

    for rows.Next() {
        question := &Question{}

        err = rows.Scan(&question.ID, &question.QuizID, &question.Prompt, 
            pq.Array(&question.Choices), &question.CorrectIndex, 
            &question.CreatedAt, &question.Version)
        if err != nil {
            return nil, err
        }
        questions = append(questions, question)
    }

    err = rows.Err()
    if err != nil {
        return nil, err
    }

    return questions, nil
}

// AddQuestion creates a question on the specified quizID and returns the newly created
// question.
func (m QuestionModel) AddQuestion(question *Question) error {

    stmt := `INSERT INTO question (prompt, choices, correct_index, quiz_id)
        VALUES ($1, $2, $3, $4)
        RETURNING id, version, created_at`
    ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
    defer cancel()

    row := m.DB.QueryRowContext(
        ctx, stmt, 
        question.Prompt, 
        pq.Array(question.Choices), 
        question.CorrectIndex, 
        question.QuizID)

    err := row.Scan(&question.ID, &question.Version, &question.CreatedAt)
    if err != nil {
        if strings.Contains(err.Error(), "foreign key") {
            return ErrNoRecords
        } else {
        return err
        }
    }

    return nil
}

// Update updates the specified question
func (m QuestionModel) Update(question *Question) error {
    stmt := `UPDATE question 
    SET prompt = $1, choices = $2, correct_index = $3, version = version + 1
    WHERE id = $4 AND version = $5
    RETURNING version`

    ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
    defer cancel()

    err := m.DB.QueryRowContext(ctx, stmt,
        question.Prompt, pq.Array(question.Choices), question.CorrectIndex, question.ID, 
        question.Version).Scan(&question.Version)
    if err != nil {
        switch {
        case err == sql.ErrNoRows:
            return ErrEditConflict
        default:
            return err
        }
    }

    return nil
}
                    







