package challenges

import (
	"testing"

	"appengine"
	"appengine/aetest"
	"errors"
	"net/http"
	"net/http/httptest"
)

var c aetest.Context
var w *httptest.ResponseRecorder
var r *http.Request

func setup() {
	var err error
	c, err = aetest.NewContext(nil)
	if err != nil {
		panic(err.Error())
	}
	r, err = http.NewRequest("GET", "/", nil)
	if err != nil {
		panic(err.Error())
	}
}

func close() {
	c.Close()
}

func TestSuccess(t *testing.T) {
	setup()
	defer close()

	var tests = []struct {
		handler    handlerFun
		wantStatus int
		wantBody   string
	}{
		{
			func(c appengine.Context, w http.ResponseWriter, r *http.Request) (interface{}, *Error) {
				return struct{ Msg string }{"hello"}, nil
			},
			http.StatusOK,
			`{"Msg":"hello"}`,
		},
		{
			func(c appengine.Context, w http.ResponseWriter, r *http.Request) (interface{}, *Error) {
				return nil, &Error{errors.New("BOOM"), "it went boom", http.StatusTeapot}
			},
			http.StatusTeapot,
			"it went boom\n",
		},
	}
	for _, tc := range tests {
		w = httptest.NewRecorder()
		h := WithContext(tc.handler, testCtx)
		h.ServeHTTP(w, r)

		if w.Code != tc.wantStatus {
			t.Errorf("Wanted status code %d but got %d", tc.wantStatus, w.Code)
		}

		if w.Body.String() != tc.wantBody {
			t.Errorf("Wanted body '%v' but got '%v'", tc.wantBody, w.Body.String())
		}
	}
}
