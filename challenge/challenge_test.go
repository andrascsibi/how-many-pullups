package challenge

import (
	"testing"

	"appengine/aetest"
	"appengine/datastore"
	"appengine/user"

	"net/http"
	"net/http/httptest"

	"fmt"
	"os"
	"bytes"
	"encoding/json"

	"github.com/andrascsibi/how-many-pullups/account"
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

func TestSuccess(t *testing.T) {
	setup()
	defer close()

	wantA := account.Account{ID: "foo", Email: "a@b", ScreenName: "dude"}

	if _, err := datastore.Put(c, account.NewKey(c, "foo"), &wantA); err != nil {
		t.Fatal(err)
	}

	wantC := Challenge{AccountID: "foo", Title: "プルアップ"}
	u := user.User{Email: wantA.Email, Admin: false}
	c.Login(&u)

	bytesIn, _ := json.Marshal(wantC)
	r, e := http.NewRequest("POST", "/", bytes.NewReader(bytesIn))
	if e != nil {
		t.Fatal(e)
	}

	v["accountId"] = wantC.AccountID

	if got, err := createChallenge(c, w, r, v); err != nil {
		t.Fatal(err)
	} else {
		c := got.(Challenge)
		wantC.ID = c.ID
		v["challengeId"] = wantC.ID
	}

	if got, err := getChallenge(c, w, r, v); err != nil {
		t.Fatal(err)
	} else {
		c := got.(Challenge)
		if got, want := c.Title, wantC.Title; got != want {
			t.Errorf("Got title %v, want %v", got, want)
		}
	}
}
