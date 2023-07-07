package main

import (
	"net/http"
	"reflect"
	"testing"
	"vanshadhruvp/quizme-api/internal/data"
)

func TestShowAllQuizzesHandler(t *testing.T) {
	ts := newTestServer(t)
	defer ts.Close()
	_, code, body := testGET[[]*data.Quiz](t, ts, "/v1/quiz")

	expectedCode := http.StatusOK
	expectedBody := []*data.Quiz{quiz1, quiz2}

	if code != expectedCode {
		t.Fatalf("INCORRECT STATUS CODE: expected %d, got %d", expectedCode, code)
	}

	if !reflect.DeepEqual(body, expectedBody) {
		t.Fatalf("INCORRECT BODY: expected %+v, got %+v", expectedBody, body)
	}
}

func TestShowQuizHandler(t *testing.T) {
	ts := newTestServer(t)
	defer ts.Close()
	type quizInstance struct {
		Quiz      *data.Quiz
		Questions []*data.Question
	}

	_, code, body := testGET[quizInstance](t, ts, "/v1/quiz/1")

	expectedCode := http.StatusOK
    expectedBody := quizInstance{Quiz: quiz1, Questions: []*data.Question{&question1, &question2,}}

	if code != expectedCode {
		t.Fatalf("INCORRECT STATUS CODE: expected %d, got %d", expectedCode, code)
	}

	if !reflect.DeepEqual(body.Quiz, expectedBody.Quiz) {
		t.Fatalf("INCORRECT BODY: expected %+v, got %+v", expectedBody, body)
	}

    if len(body.Questions) != len(expectedBody.Questions) {
        t.Fatalf("INCORRECT ENTRY QUESTIONS LENGTH: expected %d, got %d", 
            len(expectedBody.Questions), len(body.Questions))
    }

    for i, question := range body.Questions {
        expectedQuestion := expectedBody.Questions[i]
        if !reflect.DeepEqual(*question, *expectedQuestion) {
            t.Fatalf("INCORRECT ENTRY QUESTION: expected %v, got %v",
                *expectedQuestion, *question)
        }
    }
    
}

func TestAddQuizHandler(t *testing.T) {
	// Test valid quiz request
	t.Run("ValidRequest", func(t *testing.T) {
		ts := newTestServer(t)
		defer ts.Close()
		_, code, body := testPOST[string](t, ts, "/v1/quiz", []byte(`{"title": "quiz"}`))

		expectedCode := http.StatusCreated
		expectedBody := "quiz"

		if code != expectedCode {
			t.Fatalf("INCORRECT STATUS CODE: expected %d, got %d", expectedCode, code)
		}

		if !reflect.DeepEqual(body, expectedBody) {
			t.Fatalf("INCORRECT BODY: expected %q, got %q", expectedBody, body)
		}

	})

}

func TestAddQuestionHandler(t *testing.T) {
    ts := newTestServer(t)
    defer ts.Close()
    tests := []struct{
        name string
        quizID string
        payload []byte
        expectedCode int
        expectedErr error
        expected data.Question
    }{
        {"Valid post question", 
            "1", 
            []byte(`{"prompt":"q1", "choices": ["a", "b", "c"], "correct_index":2}`),
            http.StatusCreated,
            nil, 
            data.Question{
                Prompt: "q1",
                Choices: []string{"a", "b", "c"},
                CorrectIndex: 2,
                Version: 1,
                CreatedAt: testDate,
                ID: 1,
                QuizID: 1,
            }},
        {"Quiz does not exist", 
            "3", 
            []byte(`{"prompt":"q1", "choices": ["a", "b", "c"], "correct_index":2}`),
            http.StatusNotFound,
            data.ErrNoRecords, 
            data.Question{}},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            _, code, body := testPOST[data.Question](t, ts, "/v1/quiz/" + tt.quizID + "/question", tt.payload) 
            
            if code != tt.expectedCode {
                t.Fatalf("INCORRECT STATUS CODE: expected %d, got %d", tt.expectedCode, code)
            }

            if code != http.StatusCreated {
                return
            }

            if !reflect.DeepEqual(body, tt.expected) {
                t.Fatalf("INCORRECT BODY: expected %v, got %v", tt.expected, body)
            }

            
        })
    }
}


func TestAddScoreHandler(t *testing.T) {
	ts := newTestServer(t)
	defer ts.Close()
    _, code, body := testPOST[int32](t, ts, "/v1/quiz/1/score", []byte(`{"answers": [1, 2]}`))

	expectedCode := http.StatusOK
	expectedBody := int32(50) 

	if code != expectedCode {
		t.Fatalf("INCORRECT STATUS CODE: expected %d, got %d", expectedCode, code)
	}

	if !reflect.DeepEqual(body, expectedBody) {
		t.Fatalf("INCORRECT BODY: expected %d, got %d", expectedBody, body)
	}
} 

