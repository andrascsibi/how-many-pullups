package handler

import (
	"testing"

	"appengine"
	"appengine/aetest"
	"errors"
	"fmt"
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
	w = httptest.NewRecorder()
	r, err = http.NewRequest("GET", "/", nil)
	if err != nil {
		panic(err.Error())
	}
}

func close() {
	c.Close()
}

func testCtx(r *http.Request) appengine.Context {
	return c
}

func TestSuccess(t *testing.T) {
	setup()
	defer close()

	h := WithContext(func(c appengine.Context, w http.ResponseWriter, r *http.Request) (interface{}, *Error) {
		return struct{ Msg string }{"hello"}, nil
	}, testCtx)

	h.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Log("expected 200 OK, got:", w.Code)
		t.Fail()
	}

	if w.Body.String() != "{\"Msg\":\"hello\"}" {
		t.Log("expected response body other than", w.Body.String())
		t.Fail()
	}
}

func TestErr(t *testing.T) {
	setup()
	defer close()

	h := WithContext(func(c appengine.Context, w http.ResponseWriter, r *http.Request) (interface{}, *Error) {
		return nil, &Error{errors.New("BOOM"), "it went boom", http.StatusTeapot}
	}, testCtx)

	h.ServeHTTP(w, r)

	fmt.Printf("%d - %s", w.Code, w.Body.String())
	if w.Code != http.StatusTeapot {
		t.Log("expected 418 OK, got:", w.Code)
		t.Fail()
	}

	if w.Body.String() != "{\"error\":\"it went boom\"}\n" {
		t.Log("expected response body other than", w.Body.String())
		t.Fail()
	}
}
