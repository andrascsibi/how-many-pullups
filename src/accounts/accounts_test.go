package accounts

import (
	"testing"

	"appengine"
	"appengine/aetest"
	"appengine/datastore"
	"appengine/user"

	"net/http"
	"net/http/httptest"
)

var c aetest.Context
var w *httptest.ResponseRecorder
var r *http.Request
var v map[string]string

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
	v = make(map[string]string)
}

func close() {
	c.Close()
}

func testCtx(r *http.Request) appengine.Context {
	return c
}

func TestAllHandlers(t *testing.T) {
	setup()
	defer close()

	key := NewKey(c, "foo")
	if _, err := datastore.Put(c, key, &Account{ID: "foo", Email: "a@b", ScreenName: "dude"}); err != nil {
		t.Fatal(err)
	}

	got, err := getAccounts(c, w, r, v)
	if got != nil || err == nil || err.Code != http.StatusForbidden {
		t.Errorf("Wanted to be forbidden")
	}

	u := user.User{Email: "a@b", Admin: false}
	c.Login(&u)
	got, err = getAccounts(c, w, r, v)
	if got != nil || err == nil || err.Code != http.StatusUnauthorized {
		t.Errorf("Wanted to be unathorized")
	}

	u.Admin = true
	c.Login(&u)
	got, err = getAccounts(c, w, r, v)
	if got == nil || err != nil {
		t.Errorf("Wanted to get some accounts back")
	}

	// XXX query result is not consistent yet
	// as := got.([]Account)
	// if len(as) != 1 || as[0].ID != "foo" {
	// 	t.Errorf("Wanted to get 'foo' back, got %v", as)
	// }

	v["accountId"] = "foo"
	got, err = getAccount(c, w, r, v)
	if err != nil {
		t.Fatal(err)
	}
	a := got.(*Account)
	if got, want := a.Email, ""; got != want {
		t.Errorf("Got email %v, want %v", got, want)
	}
	if got, want := a.ScreenName, "dude"; got != want {
		t.Errorf("Got screen name %v, want %v", got, want)
	}

}
