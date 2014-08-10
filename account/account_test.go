package account

import (
	"testing"

	"appengine/aetest"
	"appengine/datastore"
	"appengine/user"

  "fmt"
  "os"

	"net/http"
	"net/http/httptest"
)

var c aetest.Context
var w *httptest.ResponseRecorder
var r *http.Request
var v map[string]string

func setup() {
	var err error
	opts := aetest.Options{AppID: fmt.Sprintf("app-test-%v", os.Getpid())}
	c, err = aetest.NewContext(&opts)
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

func TestAllHandlers(t *testing.T) {
	setup()
	defer close()

	wantA := Account{ID: "foo", Email: "a@b", ScreenName: "dude"}
	key := NewKey(c, "foo")
	if _, err := datastore.Put(c, key, &wantA); err != nil {
		t.Fatal(err)
	}

	// Authorize tests
	if gotE, wantC := wantA.Authorize(c), http.StatusForbidden; gotE == nil || gotE.Code != wantC {
		t.Errorf("Wanted %v, got %v", wantC, gotE)
	}

	u := user.User{Email: wantA.Email, Admin: false}
	c.Login(&u)
	if gotE := wantA.Authorize(c); gotE != nil {
		t.Errorf("Wanted to be let through because email matches")
	}

	u.Email = "b@c"
	c.Login(&u)
	if gotE, wantC := wantA.Authorize(c), http.StatusUnauthorized; gotE == nil || gotE.Code != wantC {
		t.Errorf("Wanted %v, got %v", wantC, gotE)
	}

	u.Email = "b@c"
	u.Admin = true
	c.Login(&u)
	if gotE := wantA.Authorize(c); gotE != nil {
		t.Errorf("Wanted to be let through because admin")
	}

	// Get account test

	v["accountId"] = "foo"
	if got, err := getAccount(c, w, r, v); err != nil {
		t.Fatal(err)
	} else {
		a := got.(Account)
		if got, want := a.Email, ""; got != want {
			t.Errorf("Got email %v, want %v", got, want)
		}
		if got, want := a.ScreenName, wantA.ScreenName; got != want {
			t.Errorf("Got screen name %v, want %v", got, want)
		}
	}

}
