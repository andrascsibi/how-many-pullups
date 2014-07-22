package hanlder

import (
	"reflect"
	"testing"

	"httptest"

	"appengine/mail"
)

func simpleHandler(w http.ResponseWriter, r *http.Request) (interface{}, *Error) {
	return struct{ Msg string }{"hello"}
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
