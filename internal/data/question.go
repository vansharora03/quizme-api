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
    ID int64 // ID
    QuizID int64 // The id of the quiz that this question belongs to
    Prompt string // The prompt that is to be answered
    Choices []string // The possible answer choices
    CorrectIndex int32 // The index at which the correct answer is found
    CreatedAt time.Time // When the question was created
    Version int32 // The version of the question
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









