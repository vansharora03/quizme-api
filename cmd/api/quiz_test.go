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

func TestShowQuizHandler(t *testing.T) {
    ts := newTestServer(t)
    defer ts.Close()
    _, code, body := testGET[string](t, ts, "/v1/quiz/1")

    expectedCode := http.StatusOK
    expectedBody := "Showing specific quiz..."

    if code != expectedCode {
        t.Fatalf("INCORRECT STATUS CODE: expected %d, got %d", expectedCode, code)
    }

    if !reflect.DeepEqual(body, expectedBody) {
        t.Fatalf("INCORRECT BODY: expected %q, got %q", expectedBody, body)
    }
}

func TestAddQuizHandler(t *testing.T) {
    ts := newTestServer(t)
    defer ts.Close()
    _, code, body := testPOST[string](t, ts, "/v1/quiz")

    expectedCode := http.StatusOK
    expectedBody := "Adding a quiz..."

    if code != expectedCode {
        t.Fatalf("INCORRECT STATUS CODE: expected %d, got %d", expectedCode, code)
    }

    if !reflect.DeepEqual(body, expectedBody) {
        t.Fatalf("INCORRECT BODY: expected %q, got %q", expectedBody, body)
    }
}

func TestAddScoreHandler(t *testing.T) {
    ts := newTestServer(t)
    defer ts.Close()
    _, code, body := testPOST[string](t, ts, "/v1/quiz/1/score")

    expectedCode := http.StatusOK
    expectedBody := "Adding a score..."

    if code != expectedCode {
        t.Fatalf("INCORRECT STATUS CODE: expected %d, got %d", expectedCode, code)
    }

    if !reflect.DeepEqual(body, expectedBody) {
        t.Fatalf("INCORRECT BODY: expected %q, got %q", expectedBody, body)
    }
}
