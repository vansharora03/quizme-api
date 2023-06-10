package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// writeJSONSubtest contains data about a particular subtest for
// TestWriteJSON
type writeJSONSubtest struct {
	name         string      // name of subtest
	code         int         // status code subtest will send
	data         interface{} // data subtest will send
	expectedJSON []byte      // json that is expected to be produced
}

func TestWriteJSON(t *testing.T) {
	tests := []writeJSONSubtest{
		{name: "Should work for structs",
			data: struct{ Hello string }{Hello: "world"}, code: 200, expectedJSON: []byte("{\"Helo\":\"world\"}")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			subtestWriteJson(t, tt)
		})
	}

}

func subtestWriteJson(t *testing.T, tt writeJSONSubtest) {
	rr := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "", nil)
	check(t, err)
	err = newTestApp(t).writeJSON(rr, r, tt.code, tt.data, nil)
	check(t, err)
	body, err := io.ReadAll(rr.Body)
	check(t, err)

	if rr.Code != tt.code {
		t.Fatalf("WRONG STATUS CODE: expected '%d' got '%d'", tt.code, rr.Code)
	}

	if !bytes.Contains(body, tt.expectedJSON) {
		t.Fatalf("INCORRECT JSON RESULT: expected '%q' got '%q'", tt.expectedJSON, body)
	}
}
