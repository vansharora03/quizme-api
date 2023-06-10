package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

// readJSONSubtest contains data about a particular subtest
// for TestReadJSON
type readJSONSubtest[T any] struct {
	name           string // name of test
	js             []byte // input json to be read
	dst            T      // pointer to be filled with data
	expectedOutput T      // expected value of dst
}

type testItem struct {
	Hello int
}

func TestReadJSON(t *testing.T) {
	testItemTests := []readJSONSubtest[testItem]{
		{name: "Should work for structs", js: []byte("{\"hello\": 123}"), dst: testItem{Hello: 0}, expectedOutput: testItem{Hello: 123}},
	}

	for _, tt := range testItemTests {
		t.Run(tt.name, func(t *testing.T) {
			subtestReadJSON[testItem](t, tt)

		})
	}
}

func subtestReadJSON[T any](t *testing.T, tt readJSONSubtest[T]) {
	rr := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "", bytes.NewBuffer(tt.js))
	check(t, err)
	err = newTestApp(t).readJSON(rr, r, &tt.dst)
	check(t, err)

	if !reflect.DeepEqual(tt.dst, tt.expectedOutput) {
		t.Fatalf("INCORRECT GO VALUE FROM JSON: expected %v got %v", tt.expectedOutput, tt.dst)
	}
}
