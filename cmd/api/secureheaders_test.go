package main

import (
	"testing"
)

func TestSecureHeader(t *testing.T) {
	// Create a new test server
	ts := newTestServer(t)
	defer ts.Close()

	header, _, _ := ts.GET(t, "/v1/healthcheck")

	// Check if the server is listening and returns the expected headers

	// check x-frame-options header
	xFrameOptionsHeader := header.Get("X-Frame-Options")
	if xFrameOptionsHeader != "deny" {
		t.Errorf("Expected X-Frame-Options: %s, but got %s", "deny", xFrameOptionsHeader)
	}

	// check x-xss-protection header
	xXSSProtectionHeader := header.Get("X-XSS-Protection")
	if xXSSProtectionHeader != "1; mode=block" {
		t.Errorf("Expected X-XSS-Protection: %s, but got %s", "1; mode=block", xXSSProtectionHeader)
	}

}
