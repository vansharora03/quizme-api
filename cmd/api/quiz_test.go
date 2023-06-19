package main

import (
	"reflect"
	"testing"
    "net/http"
)

func TestShowAllQuizzesHandler(t *testing.T) {
    ts := newTestServer(t)
    defer ts.Close()
    _, code, body := testGET[string](t, ts, "/v1/quiz")

    expectedCode := http.StatusOK
    expectedBody := "Showing all quizzes..."

    if code != expectedCode {
        t.Fatalf("INCORRECT STATUS CODE: expected %d, got %d", expectedCode, code)
    }

    if !reflect.DeepEqual(body, expectedBody) {
        t.Fatalf("INCORRECT BODY: expected %q, got %q", expectedBody, body)
    }
}
