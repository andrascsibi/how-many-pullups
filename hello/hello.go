package hello

import (
    "html/template"
    "net/http"
    "time"
    "strconv"
    "fmt"

    "appengine"
    "appengine/datastore"
)

type PullupSet struct {
    Reps int
    Date time.Time
}

func init() {
    http.HandleFunc("/admin", admin)
    http.HandleFunc("/add", add)
    http.HandleFunc("/total", total)
    http.HandleFunc("/", root)
}

func pullupSetKey(c appengine.Context) *datastore.Key {
    return datastore.NewKey(c, "Pullups", "andras_pullups", 0, nil)
}

func totalPullups(c appengine.Context) (stat Stat, err error) {
    q := datastore.NewQuery("PullupSet").Ancestor(pullupSetKey(c)).Order("-Date")
    sets := make([]PullupSet, 0)
    _, err = q.GetAll(c, &sets)
    if err != nil {
      return
    }
    today := time.Now().Truncate(24 * time.Hour)
    for _, s := range sets {
      if today.Equal(s.Date.Truncate(24 * time.Hour)) {
        stat.Today += s.Reps
      }
      stat.Total += s.Reps
    }
    return
}

func admin(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    stat, err := totalPullups(c)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
    if err = adminTemplate.Execute(w, stat.Today); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

var adminTemplate = template.Must(template.New("admin").ParseFiles("tmpl/admin"))

type Stat struct {
  Today int
  Total int
}

func root(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    stat, err := totalPullups(c)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
    if err = rootTemplate.Execute(w, stat); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

var rootTemplate = template.Must(template.New("root").ParseFiles("tmpl/root"))

func total(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    stat, err := totalPullups(c)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
    fmt.Fprint(w, stat)
}

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
