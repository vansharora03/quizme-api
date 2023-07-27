package data

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"strings"
	"time"
	"vanshadhruvp/quizme-api/internal/validator"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
    ID int64 `json:"id"`
    Username string `json:"username"`
    PlaintextPassword string `json:"-"`
    HashedPassword []byte `json:"-"`
    Email string `json:"email"`
    CreatedAt time.Time `json:"created_at"`
    Version int32 `json:"version"`
}

type UserModel struct {
    DB *sql.DB
}


// HashUser will create a hashed version of the user's password and store it in user.HashedPassword
func HashUser(user *User) error {
    hashed, err := bcrypt.GenerateFromPassword([]byte(user.PlaintextPassword), 12)
    if err != nil {
        return err
    }

    user.HashedPassword = hashed

    return nil
}

func (user *User) MatchPassword(password string) error {
    err := bcrypt.CompareHashAndPassword(user.HashedPassword, []byte(password))
    return err
}

// ValidateEmail will run all email validation checks
func ValidateEmail(v *validator.Validator, email string) {
    if email == "" {
        v.Add("email", "must be provided")
        return
    }
    emailRX := regexp.MustCompile(`^([a-zA-Z0-9_\-\.]+)@([a-zA-Z0-9_\-]+)(\.[a-zA-Z]{2,5}){1,2}$`)
    v.Check(emailRX.MatchString(email), "email", "must be a valid email address")
}

// ValidatePassword will run all password validation checks
func ValidatePassword(v *validator.Validator, password string) {
    if password == "" {
        v.Add("password", "must be provided")
        return
    }
    v.Check(len(password) >= 8, "password", "must be 8 characters or more")
    v.Check(len(password) <= 100, "password", "must be less than 100 characters")
}

// ValidateUser will run all user related validation checks
func ValidateUser(v *validator.Validator, user *User) {
    ValidateEmail(v, user.Email)
    ValidatePassword(v, user.PlaintextPassword)
    v.Check(user.Username != "", "username", "must be provided")
    v.Check(len(user.Username) <= 50, "username", "must be less than 50 characters")
}

// AddUser will add the user to the database
func (m UserModel) AddUser(user *User) error {
    stmt := `INSERT INTO user_account (username, hashed_password, email)
    VALUES ($1, $2, $3)
    RETURNING created_at, version, id`

    ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
    defer cancel()

    err := m.DB.QueryRowContext(ctx, stmt, user.Username, user.HashedPassword, user.Email).Scan(
        &user.CreatedAt, &user.Version, &user.ID)
    if err != nil {
        switch {
        case strings.Contains(err.Error(), "user_account_email_key"):
            return ErrDuplicateEmail
        default:
            return err
        }
    }

    return nil
}

// GetUserByEmail fetches the user by the given email
func (m UserModel) GetUserByEmail(email string) (*User, error) {
    stmt := `SELECT id, email, username, hashed_password, created_at, version
    FROM user_account
    WHERE email = $1`

    ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
    defer cancel()

    var user User

    err := m.DB.QueryRowContext(ctx, stmt, email).Scan(
        &user.ID, 
        &user.Email,
        &user.Username, 
        &user.HashedPassword,
        &user.CreatedAt,
        &user.Version)
    if err != nil {
        switch {
        case errors.Is(err, sql.ErrNoRows):
            return nil, ErrNoRecords
        default:
            return nil, err
        }
    }

    return &user, nil
}











