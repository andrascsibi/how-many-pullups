package hello

import (
    "html/template"
    "net/http"
    "time"
    "strconv"

    "appengine"
    "appengine/datastore"
)

type PullupSet struct {
    Reps int
    Date time.Time
}

func init() {
    http.HandleFunc("/", root)
    http.HandleFunc("/admin", admin)
    http.HandleFunc("/add", add)
}

func pullupSetKey(c appengine.Context) *datastore.Key {
    return datastore.NewKey(c, "Pullups", "andras_pullups", 0, nil)
}

func totalPullups(c appengine.Context) (sum int, err error) {
    q := datastore.NewQuery("PullupSet").Ancestor(pullupSetKey(c)).Order("-Date")
    sets := make([]PullupSet, 0)
    _, err = q.GetAll(c, &sets)
    if err != nil {
      return
    }
    for _, s := range sets {
      sum += s.Reps
    }
    return
}

func admin(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    sum, err := totalPullups(c)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
    if err = adminTemplate.Execute(w, sum); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func root(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    sum, err := totalPullups(c)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
    if err = rootTemplate.Execute(w, sum); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

var rootTemplate = template.Must(template.New("root").ParseFiles("tmpl/root"))
var adminTemplate = template.Must(template.New("admin").ParseFiles("tmpl/admin"))

func add(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    reps, err := strconv.Atoi(r.FormValue("reps"))
    if err != nil {
      http.Error(w, "that's not a number yo", http.StatusBadRequest)
    }
    ps := PullupSet {
        Reps: reps,
        Date: time.Now(),
    }
    key := datastore.NewIncompleteKey(c, "PullupSet", pullupSetKey(c))
    _, err = datastore.Put(c, key, &ps)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    http.Redirect(w, r, "/admin", http.StatusFound)
}
