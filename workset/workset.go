package workset

import (
  "github.com/gorilla/mux"
  "net/http"

  "encoding/csv"
  "encoding/json"

  "io"
  "io/ioutil"

  "fmt"
  "strconv"
  "time"

  "appengine"
  "appengine/datastore"

  "github.com/andrascsibi/how-many-pullups/account"
  "github.com/andrascsibi/how-many-pullups/challenge"
  "github.com/andrascsibi/how-many-pullups/handler"
)

type WorkSet struct {
  Reps int
  Date time.Time
}

const kind = "WorkSets"

func NewKey(c appengine.Context, accountId string, challengeId string) *datastore.Key {
  return datastore.NewIncompleteKey(c, kind, challenge.NewKey(c, accountId, challengeId))
}

type ChallengeStat struct {
  Today   int
  Total   int
  MinDate time.Time
  MaxDate time.Time
}

func RegisterHandlers(r *mux.Router) {

  r.Handle("/accounts/{accountId}/challenges/{challengeId}/stats",
    handler.New(getStats)).
    Methods("GET")

  r.Handle("/accounts/{accountId}/challenges/{challengeId}/sets",
    handler.New(export)).
    Methods("GET")
  r.Handle("/accounts/{accountId}/challenges/{challengeId}/sets",
    handler.New(createSet)).
    Methods("POST")

  r.Handle("/accounts/{accountId}/challenges/{challengeId}/export-csv",
    handler.New(exportCsv)).
    Methods("GET")
  r.Handle("/accounts/{accountId}/challenges/{challengeId}/export",
    handler.New(export)).
    Methods("GET")
  r.Handle("/accounts/{accountId}/challenges/{challengeId}/import-csv",
    handler.New(importCsv)).
    Methods("POST")
}

func export(c appengine.Context, w http.ResponseWriter, r *http.Request, v map[string]string) (interface{}, *handler.Error) {
  accountId := v["accountId"]
  challengeId := v["challengeId"]

  q := datastore.NewQuery(kind).Ancestor(challenge.NewKey(c, accountId, challengeId)).Order("-Date")
  sets := make([]WorkSet, 0)
  _, err := q.GetAll(c, &sets)

  if err != nil {
    return nil, &handler.Error{err, "Error reading sets", http.StatusInternalServerError}
  }

  return sets, nil
}

func exportCsv(c appengine.Context, w http.ResponseWriter, r *http.Request, v map[string]string) (interface{}, *handler.Error) {
  s, err := export(c, w, r, v)
  sets := s.([]WorkSet)

  if err != nil {
    return nil, err
  }

  w.Header().Set("Content-type", "application/csv")
  cw := csv.NewWriter(w)
  cw.Write([]string{"timestamp", "reps"})
  for _, s := range sets {
    cw.Write([]string{s.Date.Format(time.RFC3339), strconv.Itoa(s.Reps)})
  }
  cw.Flush()

  return nil, nil
}

func importCsv(c appengine.Context, w http.ResponseWriter, r *http.Request, v map[string]string) (interface{}, *handler.Error) {
  accountId := v["accountId"]
  challengeId := v["challengeId"]

  account, aerr := account.ById(c, accountId)
  if aerr != nil {
    return nil, aerr
  }

  authE := account.Authorize(c)
  if authE != nil {
    return nil, authE
  }

  csvIn := csv.NewReader(r.Body)
  importedSets := make([]WorkSet, 0, 1000)

  for i := 0; ; i++ {
    line, err := csvIn.Read()
    if err == io.EOF {
      break
    } else if err != nil {
      return nil, &handler.Error{err, "Could not read request", http.StatusBadRequest}
    }

    if len(line) != 2 {
      return nil, &handler.Error{err, fmt.Sprintf("Each line should contain 2 fields. Line no: %d '%v'", i, line), http.StatusBadRequest}
    }
    if line[0] == "timestamp" && line[1] == "reps" {
      continue
    }
    date, err := time.Parse(time.RFC3339, line[0])
    if err != nil {
      return nil, &handler.Error{err, fmt.Sprintf("Malformed date in line: %d '%v'", i, line[0]), http.StatusBadRequest}
    }
    reps, err := strconv.Atoi(line[1])
    if err != nil {
      return nil, &handler.Error{err, fmt.Sprintf("Malformed number in line: %d '%v'", i, line[1]), http.StatusBadRequest}
    }
    if reps == 0 {
      continue
    }
    if i > 1000 {
      return nil, &handler.Error{err, "Too many lines", http.StatusBadRequest}
    }
    importedSets = append(importedSets, WorkSet{Date: date, Reps: reps})
  }

  keys := make([]*datastore.Key, len(importedSets))
  for i := 0; i < len(importedSets); i++ {
    keys[i] = NewKey(c, accountId, challengeId)
  }
  _, err := datastore.PutMulti(c, keys, importedSets)
  if err != nil {
    return nil, &handler.Error{err, "Error storing in datastore", http.StatusInternalServerError}
  }

  return nil, nil
}

// TODO only works for utc
func getStats(c appengine.Context, w http.ResponseWriter, r *http.Request, v map[string]string) (interface{}, *handler.Error) {
  accountId := v["accountId"]
  challengeId := v["challengeId"]

  q := datastore.NewQuery(kind).Ancestor(challenge.NewKey(c, accountId, challengeId)).Order("-Date")
  sets := make([]WorkSet, 0)
  _, err := q.GetAll(c, &sets)

  if err != nil {
    return nil, &handler.Error{err, "Error reading sets", http.StatusInternalServerError}
  }

  var stat ChallengeStat
  if len(sets) == 0 {
    return stat, nil
  }
  stat.MaxDate = sets[0].Date
  stat.MinDate = sets[len(sets)-1].Date
  today := time.Now().Truncate(24 * time.Hour)
  for _, s := range sets {
    if today.Equal(s.Date.Truncate(24 * time.Hour)) { // TODO only works for utc
      stat.Today += s.Reps
    }
    stat.Total += s.Reps
  }

  return stat, nil
}

func createSet(c appengine.Context, w http.ResponseWriter, r *http.Request, v map[string]string) (interface{}, *handler.Error) {
  accountId := v["accountId"]
  challengeId := v["challengeId"]

  account, err := account.ById(c, accountId)
  if err != nil {
    return nil, err
  }

  authE := account.Authorize(c)
  if authE != nil {
    return nil, authE
  }

  data, e := ioutil.ReadAll(r.Body)
  if e != nil {
    return nil, &handler.Error{e, "Could not read request", http.StatusBadRequest}
  }

  var newSet WorkSet
  e = json.Unmarshal(data, &newSet)
  if e != nil {
    return nil, &handler.Error{e, "Could not parse JSON", http.StatusBadRequest}
  }

  newSet.Date = time.Now()

  key := NewKey(c, accountId, challengeId)
  _, e = datastore.Put(c, key, &newSet)
  if e != nil {
    return nil, &handler.Error{e, "Error storing in datastore", http.StatusInternalServerError}
  }

  return newSet, nil
}

