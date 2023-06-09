package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testServer struct {
	*httptest.Server
}

// newTestServer returns a new mock test server.
func newTestServer(t *testing.T) *testServer {
	cfg := config{
		env: "development",
	}
	testApp := &application{
		config: cfg,
		logger: nil,
	}
	return &testServer{httptest.NewServer(testApp.routes())}
}

// GET performs a get request on ts upon
// the supplied url path and returns the
// headers, status code, and body.
func (ts *testServer) GET(t *testing.T, urlPath string) (http.Header, int, []byte) {
	// Make request
	rs, err := ts.Client().Get(ts.URL + urlPath)
	// Fail test on error
	if err != nil {
		t.Fatal(err)
	}

	defer rs.Body.Close()

	// Extract body
	body, err := io.ReadAll(rs.Body)
	// Fail test on error
	if err != nil {
		t.Fatal(err)
	}

	return rs.Header, rs.StatusCode, body

}

func TestServerListening(t *testing.T) {
	// Create a new test server
	ts := newTestServer(t)
	defer ts.Close()

	// Make a GET request to the health check endpoint
	_, statusCode, _ := ts.GET(t, "/v1/healthcheck")

	// Check if the server is listening and returns the expected status code
	expectedStatus := http.StatusOK
	if statusCode != expectedStatus {
		t.Errorf("Expected status code %d, but got %d", expectedStatusCode, statusCode)
	}

}
