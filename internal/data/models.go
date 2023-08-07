package data

import (
	"database/sql"
	"errors"
    "time"
)

var ErrNoRecords = errors.New("No records found")
var ErrEditConflict = errors.New("Edit conflict")
var ErrDuplicateEmail = errors.New("Email exists already")

type Models struct {
	Questions interface {
		GetAllByQuizID(quizID string) ([]*Question, error)
		AddQuestion(question *Question) error
        Update(question *Question) error
	}

	Quizzes interface {
		GetAll() ([]*Quiz, error)
		Get(id string) (*Quiz, error)
		Add(title string, userID int64) (*Quiz, error)
        Update(quiz *Quiz) error
	}

    Users interface {
        AddUser(user *User) error
        GetUserByEmail(email string) (*User, error)
        GetUserByToken(token string) (*User, error)
    }

    Tokens interface {
        AddToken(userID int64, lifetime time.Duration) (*Token, error) 
    }

    Scores interface {
        AddScore(score *Score) error
        GetScoresByUserAndQuiz(userID, quizID int64) ([]*Score, error)
    }
}

// NewModels initializes a Models struct with
// all of quizme's data models connected on db
func NewModels(db *sql.DB) Models {
	return Models{
		Questions: QuestionModel{db},
		Quizzes:   QuizModel{db},
        Users: UserModel{db},
        Tokens: TokenModel{db},
        Scores: ScoreModel{db},
	}
}
