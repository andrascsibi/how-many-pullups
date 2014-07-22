package handler

import (
	"testing"

	"net/http"
	"net/http/httptest"
)

func simpleHandler(w http.ResponseWriter, r *http.Request) (interface{}, *Error) {
	return struct{ Msg string }{"hello"}, nil
}

func TestSimpleHandler(t *testing.T) {

	resp := httptest.NewRecorder()
	h := New(simpleHandler)

	req, err := http.NewRequest("GET", "foo", nil)
	if err != nil {
		t.Fatal(err)
	}

	h.ServeHTTP(resp, req)

}
