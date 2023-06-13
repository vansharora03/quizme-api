package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNotFoundResponse(t *testing.T) {
	rr := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "", nil)
	check(t, err)

	newTestApp(t).notFoundResponse(rr, r)

	testErrorResponse(t, rr, http.StatusNotFound, []byte("The server could not find the requested resource"))

}

func TestServerErrorResponse(t *testing.T) {
	rr := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "", nil)
	check(t, err)

	newTestApp(t).serverErrorResponse(rr, r, nil)

	testErrorResponse(t, rr, http.StatusInternalServerError, []byte("The server encountered an error processing the request"))

}

func TestMethodNotAllowedResponse(t *testing.T) {
	rr := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "", nil)
	check(t, err)

	newTestApp(t).methodNotAllowedResponse(rr, r)

	testErrorResponse(t, rr, http.StatusMethodNotAllowed, []byte("The method GET is not supported for this resource"))
}

func testErrorResponse(t *testing.T, rr *httptest.ResponseRecorder, expectedCode int, expectedBody []byte) {
	body, err := io.ReadAll(rr.Body)
	check(t, err)

	if rr.Code != expectedCode {
		t.Fatalf("INCORRECT STATUS CODE: expected %d, got %d", expectedCode, rr.Code)
	}

	if !bytes.Contains(body, expectedBody) {
		t.Fatalf("INCORRECT RESPONSE BODY: expected %q, got %q", body, expectedBody)
	}
}
