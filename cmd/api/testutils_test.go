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
	app *application
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

	// Extract body
	body, err := io.ReadAll(rs.Body)
	// Fail test on error
	check(t, err)

	var dst T
	rr := httptest.NewRecorder()
	r, err := http.NewRequest("GET", urlPath, bytes.NewBuffer(body))
	check(t, err)
	err = ts.app.readJSON(rr, r, &dst)
	check(t, err)

	return rs.Header, rs.StatusCode, dst

}
