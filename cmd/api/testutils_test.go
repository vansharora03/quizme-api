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

func newTestApp(t *testing.T) *application {
	cfg := config{
		env: "development",
	}
	testApp := &application{
		config: cfg,
		logger: nil,
	}
	return testApp
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
