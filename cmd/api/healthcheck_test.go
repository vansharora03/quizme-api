package main

import (
	"fmt"
	"net/http"
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

	// Check if the server is listening and returns the expected body
	expectedBody := fmt.Sprintf("Environment: development\nVersion: %s", version)
	if expectedBody != string(body) {
		t.Errorf("Expected body %s, but got %s", expectedBody, string(body))
	}

}
