package workset

import (
  "testing"

  "appengine/aetest"
  "appengine/datastore"
  "appengine/user"

  "net/http"
  "net/http/httptest"

  "bytes"
  "encoding/json"
  "fmt"
  "os"

  "github.com/andrascsibi/how-many-pullups/account"
  "github.com/andrascsibi/how-many-pullups/challenge"
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
  wantC := challenge.Challenge{AccountID: "foo", ID: "bar", Title: "プルアップ"}
  if _, err := datastore.Put(c, challenge.NewKey(c, "foo", "bar"), &wantC); err != nil {
    t.Fatal(err)
  }

  u := user.User{Email: wantA.Email, Admin: false}
  c.Login(&u)

  v["accountId"] = wantC.AccountID
  v["challengeId"] = wantC.ID

  wantS := WorkSet{Reps: 12}
  bytesIn, _ := json.Marshal(wantS)
  r, e := http.NewRequest("POST", "/", bytes.NewReader(bytesIn))
  if e != nil {
    t.Fatal(e)
  }

  if got, err := createSet(c, w, r, v); err != nil {
    t.Fatal(err)
  } else {
    s := got.(WorkSet)
    if got, want := s.Reps, wantS.Reps; got != want {
      t.Errorf("Wanted %v got %v", want, got)
    }
  }

  if got, err := export(c, w, r, v); err != nil {
    t.Fatal(err)
  } else {
    ss := got.([]WorkSet)
    if len(ss) != 1 || ss[0].Reps != wantS.Reps {
      t.Errorf("Wanted to get set back")
    }
  }
}
