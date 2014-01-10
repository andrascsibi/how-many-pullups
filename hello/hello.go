package hello

import (
    "html/template"
    "net/http"
    "time"
    "strconv"
    "encoding/json"

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

//func pullupSetKey(c appengine.Context, user string) *datastore.Key {
//    return datastore.NewKey(c, "Pullups", "andras_pullups", 0, nil)
//}

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
    if err = adminTemplate.ExecuteTemplate(w, "admin", stat); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

var adminTemplate = template.Must(template.New("root").ParseFiles("tmpl/root"))

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
    if err = rootTemplate.ExecuteTemplate(w, "root", stat); err != nil {
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
    statJson, err := json.Marshal(stat)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
    w.Header().Set("Content-type", "application/json")
    w.Write(statJson)
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
