package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testServer struct {
	*httptest.Server
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
	}

	return testApp
}

func check(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

// newTestServer returns a new mock test server.
func newTestServer(t *testing.T) *testServer {
	return &testServer{httptest.NewServer(newTestApp(t).routes())}
}

// GET performs a get request on ts upon
// the supplied url path and returns the
// headers, status code, and body.
func (ts *testServer) GET(t *testing.T, urlPath string) (http.Header, int, []byte) {
	// Make request
	rs, err := ts.Client().Get(ts.URL + urlPath)
	// Fail test on error
	check(t, err)

	defer rs.Body.Close()

	// Extract body
	body, err := io.ReadAll(rs.Body)
	// Fail test on error
	check(t, err)

	return rs.Header, rs.StatusCode, body

}
