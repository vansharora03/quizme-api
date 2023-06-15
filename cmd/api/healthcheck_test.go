package main

import (
	"net/http"
	"reflect"
	"testing"
)

func TestHealthCheck(t *testing.T) {
	// Create a new test server
	ts := newTestServer(t)
	defer ts.Close()

	// Make a GET request to the health check endpoint
	headers, statusCode, body := testGET[map[string]string](t, ts, "/v1/healthcheck")

	// Check if the server is listening and returns the expected status code
	expectedStatus := http.StatusOK
	if statusCode != expectedStatus {
		t.Errorf("Expected status code %d, but got %d", expectedStatus, statusCode)
	}

	if headers.Get("Content-Type") != "application/json" {
		t.Fatalf("INCORRECT HEADER Content-type: expected %q got %q", "application/json", headers.Get("Content-Type"))
	}

	expectedBody := map[string]string{
		"environment": "development",
		"version":     "1.0.0",
	}

	if !reflect.DeepEqual(body, expectedBody) {
		t.Fatalf("INCORRECT JSON RESPONSE: expected %v, got %v", expectedBody, body)
	}

}
