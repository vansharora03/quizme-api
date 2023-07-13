package main

import (
	"net/http"
	"reflect"
	"testing"
	"vanshadhruvp/quizme-api/internal/data"
)

func TestAddUserHandler(t *testing.T) {
    ts := newTestServer(t)
    defer ts.Close()
    tests := []struct{
        name string
        payload string
        expectedCode int
        expectedBody data.User
    }{
        {
            "valid user",
            `{"username": "user123", "password": "pass12345", "email": "user@domain.com"}`,
            http.StatusCreated,
            data.User{
                Username: "user123",
                Email: "user@domain.com",
            },
        },

    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            _, code, body := testPOST[data.User](t, ts, "/v1/user", []byte(tt.payload))

            if tt.expectedCode != http.StatusCreated {
                if code != tt.expectedCode {
                    t.Fatalf("INCORRECT STATUS CODE: expected %d, got %d", tt.expectedCode, code)
                }
            }

            if !reflect.DeepEqual(body, tt.expectedBody) {
                t.Fatalf("INCORRECT BODY: expected %v, got %v", tt.expectedBody, body)
            }
        })
    }
}
