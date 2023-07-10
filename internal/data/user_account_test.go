package data

import (
	"errors"
	"reflect"
	"testing"
	"vanshadhruvp/quizme-api/internal/validator"
)

func TestValidateUser(t *testing.T) {
    tests := []struct{
        name string
        inputUser User
        expectedErrs map[string]string
    }{
        {"Valid user", User{
            Username: "username", 
            PlaintextPassword: "abcdefg123",
            Email: "email@domain.com",
        }, make(map[string]string)},
        {"Empty username and short password", User{
            Username: "", 
            PlaintextPassword: "ag123",
            Email: "email@domain.com",
        }, map[string]string{"username":"must be provided", "password":"must be 8 characters or more"}},
        {"Invalid email", User{
            Username: "username", 
            PlaintextPassword: "abcdefg123",
            Email: "emailinvalid",
        }, map[string]string{"email":"must be a valid email address"}},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            v := validator.New()

            ValidateUser(&v, &tt.inputUser)

            if !reflect.DeepEqual(v.Errors, tt.expectedErrs) {
                t.Fatalf("INCORRECT VALIDATION: expected %v, got %v", tt.expectedErrs, v.Errors)
            }
        })
    }
}

func TestAddUser(t *testing.T) {
    users := UserModel{openTestDB(t)}
    tests := []struct{
        name string
        user User
        expectedErr error
        expectedID int64
    }{
        {
            "Valid user_account", 
            User{ Username: "user123", HashedPassword: []byte("abcdef"), Email: "user@domain.com",}, 
            nil, int64(1),
        }, 
        {
            "Repeated email", 
            User{ Username: "user123", HashedPassword: []byte("abcdef"), Email: "user@domain.com",}, 
            ErrDuplicateEmail, int64(0),
        }, 
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := users.AddUser(&tt.user)
            
            if tt.expectedErr != nil {
                switch {
                case errors.Is(err, tt.expectedErr):
                    return
                default:
                    t.Fatalf("INCORRECT ERROR: expected %v, got %v", tt.expectedErr, err)
                }
            }

            if tt.user.ID != tt.expectedID {
                t.Fatalf("INCORRECT ENTRY ID: expected %d, got %d", tt.expectedID, tt.user.ID)
            }

        })

    }
}
