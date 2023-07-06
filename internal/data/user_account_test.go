package data

import (
    "vanshadhruvp/quizme-api/internal/validator"
    "testing"
    "reflect"
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
