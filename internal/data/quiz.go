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
	ID        int64     // ID
	CreatedAt time.Time // When the quiz was created
	Version   int32     // Version of the quiz
	Title     string    // Title of the quiz
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

	stmt := `SELECT id, created_at, version, title FROM quiz WHERE id = $1`

	// Create a context with a timeout of 5 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Query the database for the quiz
	row := m.DB.QueryRowContext(ctx, stmt, id)

	// Create a new Quiz struct
	q := &Quiz{}

	// Copy the values from row into the Quiz struct
	err := row.Scan(&q.ID, &q.CreatedAt, &q.Version, &q.Title)

	if err != nil {
		return nil, err
	}

	// Return the quiz
	return q, nil
}

// Add will add a quiz to the database and return the id of the quiz.
func (m QuizModel) Add(title string) (int64, error) {
	stmt := `INSERT INTO quizme (title)
    VALUES($1)`

	result, err := m.DB.Exec(stmt, title)

	if err != nil {
		return 0, err
	}

	// Get the id of the quiz
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}
