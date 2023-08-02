package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
	"vanshadhruvp/quizme-api/internal/validator"
)

type QuizModel struct {
	DB *sql.DB
}

type Quiz struct {
    ID        int64 `json:"id"` // ID
	CreatedAt time.Time `json:"created_at"` // When the quiz was created
	Version   int32 `json:"version"` // Version of the quiz
	Title     string `json:"title"` // Title of the quiz
    UserID int64 `json:"user_id"` // ID of the author
}

// ValidateQuiz runs all Quiz Checks
func ValidateQuiz(v *validator.Validator, quiz *Quiz) {
    // Run checks
    v.Check(quiz.Title != "", "title", "must be provided")
    v.Check(len(quiz.Title) <= 100, "title", "must be 100 characters or less")
}

// GetAll will get all quizzes in the database and return them.
func (m QuizModel) GetAll() ([]*Quiz, error) {
	stmt := `SELECT id, created_at, version, title FROM quiz`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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

// Get will get a specific quiz from the database and return it based on the id.
func (m QuizModel) Get(id string) (*Quiz, error) {

	stmt := `SELECT id, created_at, version, title, user_id FROM quiz WHERE id = $1`

	// Create a context with a timeout of 5 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Query the database for the quiz
	row := m.DB.QueryRowContext(ctx, stmt, id)

	// Create a new Quiz struct
	q := &Quiz{}

	// Copy the values from row into the Quiz struct
	err := row.Scan(&q.ID, &q.CreatedAt, &q.Version, &q.Title, &q.UserID)

	if err != nil {
        if err == sql.ErrNoRows {
            return nil, ErrNoRecords
        }
		return nil, err
	}

	// Return the quiz
	return q, nil
}

// Add will add a quiz to the database and return the id of the quiz.
func (m QuizModel) Add(title string, userID int64) (*Quiz, error) {
	stmt := `INSERT INTO quiz (title, user_id)
    VALUES($1, $2) RETURNING title, created_at, version, id, user_id`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	row := m.DB.QueryRowContext(ctx, stmt, title, userID)

    quiz := &Quiz{UserID: userID}

	err := row.Scan(&quiz.Title, &quiz.CreatedAt, &quiz.Version, &quiz.ID, &quiz.UserID)
	if err != nil {
		return nil, err
	}

	return quiz, nil
}




// Update will update a quiz in the database
func (m QuizModel) Update(quiz *Quiz) error {
    stmt := `UPDATE quiz
    SET title = $1, version = version + 1
    WHERE id = $2 AND version = $3
    RETURNING version, user_id`

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    err := m.DB.QueryRowContext(ctx, stmt, 
        quiz.Title, 
        quiz.ID, 
        quiz.Version).Scan(&quiz.Version, &quiz.UserID)
    
    if err != nil {
        switch {
        case errors.Is(err, sql.ErrNoRows):
            return ErrEditConflict
        default:
            return err
        }
    }

    return nil
}













