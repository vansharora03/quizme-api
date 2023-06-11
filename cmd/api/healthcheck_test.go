package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestHealthCheck(t *testing.T) {
	// Create a new test server
	ts := newTestServer(t)
	defer ts.Close()

	// Make a GET request to the health check endpoint
	_, statusCode, body := ts.GET(t, "/v1/healthcheck")

	// Check if the server is listening and returns the expected status code
	expectedStatus := http.StatusOK
	if statusCode != expectedStatus {
		t.Errorf("Expected status code %d, but got %d", expectedStatus, statusCode)
	}

	expectedBody := map[string]string{
		"environment": "development",
		"version":     "1.0.0",
	}

	dst := map[string]string{}
	r, err := http.NewRequest("GET", "/v1/healthcheck", bytes.NewBuffer(body))
	check(t, err)
	err = ts.app.readJSON(httptest.NewRecorder(), r, &dst)
	check(t, err)
	if !reflect.DeepEqual(dst, expectedBody) {
		t.Fatalf("INCORRECT JSON RESPONSE: expected %v, got %v", expectedBody, dst)
	}

}
