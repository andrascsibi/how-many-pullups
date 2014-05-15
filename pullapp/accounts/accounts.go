package accounts

import (
    "html/template"
    "net/http"
    "time"
    "errors"
    // "strconv"
    // "fmt"
    "regexp"

    "appengine"
    "appengine/user"
//    "appengine/datastore"
)

type Account struct {
    Email      string
    ScreenName string
    Admin      bool
    RegDate    time.Time

    Challenges []string
    Settings   Settings
}

type Settings struct {
    
}

type Challenge struct {
  Name   string
  Title  string
  Message string
}

func init() {
    http.HandleFunc("/hello", hello)
    http.HandleFunc("/create", create)
}

func hello(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    u := user.Current(c)
    if u == nil {
        url, err := user.LoginURL(c, r.URL.String())
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        http.Redirect(w, r, url, http.StatusFound)
        return
    }
    url, err := user.LogoutURL(c, r.URL.String())
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    d := HelloData {
        Email: u.Email,
        LogoutURL: url,
    }

    if err := rootTemplate.ExecuteTemplate(w, "hello", d); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

type HelloData struct {
    Email string
    LogoutURL string
}

func validate(username string) error {
  matches, err := regexp.Match("^[a-z][a-z0-9\\-]{2,29}$", []byte(username))
  if err != nil {
    return err
  }
  if !matches {
    return errors.New("Username must be between 3 and 30 characters long, must start with a lowercase letter, and can only contain lowercase letters, numbers, and the '-' character.")
  } 
  return nil
}

func create(w http.ResponseWriter, r *http.Request) {
    //c := appengine.NewContext(r)
    //u := user.Current(c)
    username := r.FormValue("username")
    err := validate(username)
    if err != nil {
      http.Error(w, err.Error(), http.StatusBadRequest)
    }
    http.Redirect(w, r, "/" + username, http.StatusFound)
}

var rootTemplate = template.Must(template.New("root").ParseFiles("tmpl/root"))

//func (*Challenge) GetKey() *datastore.Key {
 //TODO   return datastore.NewKey(c, "Pullups", "andras_pullups", 0, nil)
//}
