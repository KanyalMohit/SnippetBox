package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"snippetbox.mohit.net/internal/assert"
)

// The TestPing function tests the ping handler by sending a GET request to the "/" endpoint and
// checking if the response status code is OK with a body of "OK".
func TestPing(t *testing.T) {
	rr := httptest.NewRecorder()

	r, err := http.NewRequest(http.MethodGet, "/", nil)

	if err != nil {
		t.Fatal(err)
	}

	ping(rr, r)

	// `rs := rr.Result()` is a method call that returns the recorded response as an `http.Response`
	// struct. In this context, `rr` is an `httptest.ResponseRecorder` which is used to record the
	// response generated by the HTTP handler function being tested. The `Result()` method is used to get
	// the recorded response which can then be inspected to check the status code, headers, and body of
	// the response in the test case.
	rs := rr.Result()

	assert.Equal(t, rs.StatusCode, http.StatusOK)

	defer rs.Body.Close()

	// The line `body, err := io.ReadAll(rs.Body)` is reading the entire content of the response body
	// `rs.Body` into a byte slice named `body`. The `io.ReadAll` function reads from the provided
	// `io.Reader` until an error or EOF is reached, and returns the data read along with any error
	// encountered during the read operation.
	body, err := io.ReadAll(rs.Body)

	if err != nil {
		t.Fatal(err)
	}

	// The `bytes.TrimSpace(body)` function call in the code snippet is attempting to trim any leading and
	// trailing white space characters from the byte slice `body`. However, it is important to note that
	// the `bytes.TrimSpace()` function does not modify the original byte slice in place. Instead, it
	// returns a new byte slice with leading and trailing white space characters removed.
	bytes.TrimSpace(body)
	assert.Equal(t, string(body), "OK")
}
