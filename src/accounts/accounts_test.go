package accounts

import (
	"testing"

	"appengine"
	"appengine/aetest"
	"appengine/datastore"

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

func testCtx(r *http.Request) appengine.Context {
	return c
}

func TestGetAccounts(t *testing.T) {
	setup()
	defer close()

	key := datastore.NewKey(c, "Accounts", "", 1, nil)
	if _, err := datastore.Put(c, key, &Account{ID: "foo", Email: "a@b"}); err != nil {
		t.Fatal(err)
	}
	// unauthorized
	as, err := getAccounts(c, w, r)
	if as != nil || err == nil || err.Code != http.StatusForbidden {
		t.Errorf("Wanted to be forbidden")
	}

}
