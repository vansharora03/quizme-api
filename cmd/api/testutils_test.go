package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"runtime/debug"
	"testing"
	"time"
	"vanshadhruvp/quizme-api/internal/data"
)

type testServer struct {
	app *application
	*httptest.Server
}

// Mocks QuizModel
type TestQuizModel struct{}

// test quiz data
var quiz1 *data.Quiz = &data.Quiz{
	Title:     "quiz1",
	Version:   4,
	ID:        1,
	CreatedAt: time.Date(2005, time.April, 5, 1, 2, 3, 0, time.UTC),
}
var quiz2 *data.Quiz = &data.Quiz{
	Title:     "quiz2",
	Version:   1,
	ID:        2,
	CreatedAt: time.Date(2020, time.March, 20, 0, 0, 5, 0, time.UTC),
}

// Mocks data.QuizModel.GetAll to test showAllQuizzesHandler
func (m TestQuizModel) GetAll() ([]*data.Quiz, error) {
	return []*data.Quiz{quiz1, quiz2}, nil
}

// Mocks data.QuizModel.Get to test showQuizHandler
func (m TestQuizModel) Get(id string) (*data.Quiz, error) {
	// If id is 1, return quiz1
	if id == "1" {
		return quiz1, nil
		// If id is 2, return quiz2
	} else if id == "2" {
		return quiz2, nil

	}
	return nil, data.ErrNoRecords 
}

// Mocks data.QuizModel.Add to test addQuizHandler
func (m TestQuizModel) Add(title string) (string, error) {
	return "quiz", nil
}

// Mocks data.QuizModel.Update
func (m TestQuizModel) Update(quiz *data.Quiz) error {
    return nil
}

// Mocks QuestionModel
type TestQuestionModel struct{}

var question1 data.Question = data.Question{
    ID: 1,
    Version: 1,
    Prompt: "testq1",
    Choices: []string{"a", "b", "c"},
    CorrectIndex: 1,
    QuizID: 1,
}

var question2 data.Question = data.Question{
    ID: 1,
    Version: 1,
    Prompt: "testq1",
    Choices: []string{"a", "b", "c"},
    CorrectIndex: 0,
    QuizID: 1,
}

var testDate time.Time = time.Date(2000, time.November, 28, 12, 23, 12, 0, time.UTC)

func (m TestQuestionModel) GetAllByQuizID(quizID string) ([]*data.Question, error) {
    if quizID == "1" {
        return []*data.Question{&question1, &question2}, nil
    }
	return []*data.Question{}, nil
}

func (m TestQuestionModel) AddQuestion(question *data.Question) error {
    if question.QuizID != 1 && question.QuizID != 2 {
        return data.ErrNoRecords
    }
    question.CreatedAt = testDate 
    question.Version = 1
    question.ID = 1

	return nil
}

func (m TestQuestionModel) Update(question *data.Question) error {
    return nil
}

type TestUserModel struct {}

func (m TestUserModel) AddUser(user *data.User) error {
    return nil
}

// Create mock Models
func newTestModel(t *testing.T) data.Models {
	return data.Models{
		Quizzes:   TestQuizModel{},
		Questions: TestQuestionModel{},
        Users: TestUserModel{},
	}
}

func newTestApp(t *testing.T) *application {
	cfg := config{
		env: "development",
	}

	// Create a new buffer for the logger
	buffer := []byte{}
	testApp := &application{
		config: cfg,
		// Create a new logger with the buffer
		logger: log.New(bytes.NewBuffer(buffer), "", log.Default().Flags()),
		// Pass in mock models
		models: newTestModel(t),
	}

	return testApp
}

func check(t *testing.T, err error) {
	if err != nil {
		fmt.Printf("%q", debug.Stack())
		t.Fatal(err)
	}
}

// newTestServer returns a new mock test server.
func newTestServer(t *testing.T) *testServer {
	app := newTestApp(t)
	return &testServer{app, httptest.NewServer(app.routes())}
}

// testGET performs a get request on ts upon
// the supplied url path and returns the
// headers, status code, and body. Initialize
// the generic to be the type of the json response
// once converted to a Go value.
func testGET[T any](t *testing.T, ts *testServer, urlPath string) (http.Header, int, T) {
	// Make request
	rs, err := ts.Client().Get(ts.URL + urlPath)
	// Fail test on error
	check(t, err)

	defer rs.Body.Close()

	return openResponse[T](t, ts, rs, urlPath)

}

// testPOST performs a post request on ts with the specified json payload. Returns headers,
// status code, and body. Initialize the generic to be the type of the json response
// once converted to a Go value.
func testPOST[T any](
	t *testing.T, ts *testServer, urlPath string, payload []byte) (http.Header, int, T) {
	// Make request
	rs, err := ts.Client().Post(ts.URL+urlPath, "application/json", bytes.NewBuffer(payload))
    check(t, err)

	defer rs.Body.Close()

	return openResponse[T](t, ts, rs, urlPath)
}

// openResponse takes the rs response that contains json data, and attempts to extract
// the headers, status code, and body as a Go value of type T.
func openResponse[T any](
	t *testing.T, ts *testServer, rs *http.Response, urlPath string) (http.Header, int, T) {

	// Extract body
	body, err := io.ReadAll(rs.Body)
	// Fail test on error
	check(t, err)
	var dst T
	rr := httptest.NewRecorder()
	r, err := http.NewRequest("GET", urlPath, bytes.NewBuffer(body))
	check(t, err)
	_ = ts.app.readJSON(rr, r, &dst)
	return rs.Header, rs.StatusCode, dst

}
