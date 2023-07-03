package data

import (
	"database/sql"
	"regexp"
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


