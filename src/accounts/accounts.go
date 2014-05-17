package accounts

import (
    "net/http"

    "github.com/gorilla/mux"

    "time"

    "crypto/sha1"
    "io"
    "encoding/base64"
//    "errors"
    // "strconv"
//     "fmt"
//    "regexp"

    "errors"
    "strings"
    "appengine"
    "appengine/user"
    "encoding/json"
    "appengine/datastore"

)

type Account struct {
    Email      string
    ID         string
    ScreenName string
    Admin      bool
    RegDate    time.Time

    Challenges []string
    Settings   Settings
}

type Settings struct {

}

type Challenge struct {
  Title  string
  Description string
  Active bool
}

func init() {
    r := mux.NewRouter()
    r.HandleFunc("/whoami", whoami)
    r.HandleFunc("/accounts/", accounts)
    http.Handle("/", r)
}

func accounts(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    path := strings.Split(r.URL.Path, "/")
    if len(path) < 3 {
        return
    }
    accountId := path[2]

    var account Account
    err := datastore.Get(c, accountKey(c, accountId), &account)

    if err == datastore.ErrNoSuchEntity {
        http.Error(w, err.Error(), http.StatusNotFound)
        return;
    }

    // /account/{a_id} (element)
    if len(path) == 3 {
        // GET: public account info
        // POST: update info, like screen name (login required)
        return
    }

    // /account/{a_id}/challenges (collection)
    if len(path) == 4 && path[3] == "challenges" {
        c.Infof("challenges collection for %v", accountId)
        // GET: list challenges
        // POST: create new challenge
        return
    }

    // /accounts/{id}/challenges/{c_id} (element)
    if len(path) == 5 && path[3] == "challenges" {
        // GET: title, description
        // POST: update title/description
        return
    }

    // /accounts/{id}/challenges/{c_id}/sets (collection)
    if len(path) == 6 && path[3] == "challenges" && path[5] == "sets" {
        // GET: list of all sets (export)
        // POST: create new set (param: reps)
        // PUT: replace whole collection (import)
        return
    }

}

func getOrCreateAccount(c appengine.Context) (account Account, err error) {
    u := user.Current(c)
    if u == nil {
        err = errors.New("Login required.")
        return
    }
    email := u.Email
    id := hash(u.ID)

    key := accountKey(c, id)

    err = datastore.Get(c, key, &account)

    if err == datastore.ErrNoSuchEntity {
        c.Infof("creating new user: %v %v", email, id)
        account = Account {
            Email: email,
            ID: id,
            RegDate: time.Now(),
        }
        _, err = datastore.Put(c, key, &account)
    }
    return
}

func accountKey(c appengine.Context, id string) *datastore.Key {
    return datastore.NewKey(c, "Accounts", id, 0, nil)
}

func whoami(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    u := user.Current(c)

    var ret LoginData
    w.Header().Set("Content-type", "application/json")

    if u == nil {
        url, err := user.LoginURL(c, "")
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        ret.LoginURL = url
    } else {
        url, err := user.LogoutURL(c, "")
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        ret.LogoutURL = url

        account, err := getOrCreateAccount(c)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
        ret.Account = account
    }

    loginData, err := json.Marshal(ret)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
    w.Write(loginData)
}

type LoginData struct {
    Account Account
    LoginURL string
    LogoutURL string
}

func hash(id string) string {
    hasher := sha1.New()
    io.WriteString(hasher, id)
    io.WriteString(hasher, "salt it real good DbqOFzkk") // TODO: should come from environment
    return base64.URLEncoding.EncodeToString(hasher.Sum(nil))[:8]
}

