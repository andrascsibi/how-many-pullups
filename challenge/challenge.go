package challenge

import (
	"github.com/gorilla/mux"
	"net/http"

	"encoding/base64"
	"encoding/json"

	"crypto/sha1"

	"io"
	"io/ioutil"

	"fmt"
	"time"

	"appengine"
	"appengine/datastore"

	"github.com/andrascsibi/how-many-pullups/account"
	"github.com/andrascsibi/how-many-pullups/handler"
)

type Challenge struct {
	AccountID    string
	ID           string
	Title        string
	Description  string
	MaxReps      int
	StepReps     int
	CreationDate time.Time
}

const kind = "Challenges"

func NewKey(c appengine.Context, accountId string, id string) *datastore.Key {
	return datastore.NewKey(c, kind, id, 0, account.NewKey(c, accountId))
}

type ChallengeStat struct {
	Today   int
	Total   int
	MinDate time.Time
	MaxDate time.Time
}

func RegisterHandlers(r *mux.Router) {

	r.Handle("/accounts/{accountId}/challenges",
		handler.New(getChallenges)).
		Methods("GET")
	r.Handle("/accounts/{accountId}/challenges",
		handler.New(createChallenge)).
		Methods("POST")

	r.Handle("/accounts/{accountId}/challenges/{challengeId}",
		handler.New(getChallenge)).
		Methods("GET")
	r.Handle("/accounts/{accountId}/challenges/{challengeId}",
		handler.New(updateChallenge)).
		Methods("POST")
}

func getChallenges(c appengine.Context, w http.ResponseWriter, r *http.Request, v map[string]string) (interface{}, *handler.Error) {
	accountId := v["accountId"]

	q := datastore.NewQuery(kind).
		Ancestor(account.NewKey(c, accountId)).
		Order("-CreationDate")
	challenges := make([]Challenge, 0)
	_, e := q.GetAll(c, &challenges)

	if e != nil {
		return nil, &handler.Error{e, "Error querying datastore", http.StatusInternalServerError}
	}

	return challenges, nil
}

func createChallenge(c appengine.Context, w http.ResponseWriter, r *http.Request, v map[string]string) (interface{}, *handler.Error) {

	accountId := v["accountId"]
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

	var challenge Challenge
	e = json.Unmarshal(data, &challenge)
	if e != nil {
		return nil, &handler.Error{e, "Could not parse JSON", http.StatusBadRequest}
	}

	challenge.CreationDate = time.Now()
	challenge.ID = hash(challenge.CreationDate.String())
	challenge.AccountID = accountId

	key := NewKey(c, accountId, challenge.ID)
	_, e = datastore.Put(c, key, &challenge)
	if e != nil {
		return nil, &handler.Error{e, "Error storing in datastore", http.StatusInternalServerError}
	}

	return challenge, nil
}

func getChallenge(c appengine.Context, w http.ResponseWriter, r *http.Request, v map[string]string) (interface{}, *handler.Error) {
	accountId := v["accountId"]
	challengeId := v["challengeId"]

	var challenge Challenge
	err := datastore.Get(c, NewKey(c, accountId, challengeId), &challenge)

	fmt.Printf("BOOOO accountId: %v, challengeId %v", accountId, challengeId)

	if err == datastore.ErrNoSuchEntity {
		return nil, &handler.Error{err, "Challenge not found", http.StatusNotFound}
	} else if err != nil {
		return nil, &handler.Error{err, "Error accessing datastore", http.StatusInternalServerError}
	}

	return challenge, nil
}

func updateChallenge(c appengine.Context, w http.ResponseWriter, r *http.Request, v map[string]string) (interface{}, *handler.Error) {
	accountId := v["accountId"]
	challengeId := v["challengeId"]
	_ = c

	account, err := account.ById(c, accountId)
	if err != nil {
		return nil, err
	}

	authE := account.Authorize(c)
	if authE != nil {
		return nil, authE
	}

	challenge, err := getChallenge(c, w, r, v)
	if err != nil {
		return nil, err
	}
	ch := challenge.(Challenge)

	data, e := ioutil.ReadAll(r.Body)
	if e != nil {
		return nil, &handler.Error{e, "Could not read request", http.StatusBadRequest}
	}

	var updatedChallenge Challenge
	e = json.Unmarshal(data, &updatedChallenge)
	if e != nil {
		return nil, &handler.Error{e, "Could not parse JSON", http.StatusBadRequest}
	}

	// protect certain fields
	updatedChallenge.CreationDate = ch.CreationDate
	updatedChallenge.ID = ch.ID
	updatedChallenge.AccountID = ch.AccountID

	key := NewKey(c, accountId, challengeId)
	_, e = datastore.Put(c, key, &updatedChallenge)
	if e != nil {
		return nil, &handler.Error{e, "Error storing in datastore", http.StatusInternalServerError}
	}

	return updatedChallenge, nil
}

func hash(id string) string {
	hasher := sha1.New()
	io.WriteString(hasher, id)
	io.WriteString(hasher, "salt it real good DbqOFzkk") // TODO: should come from environment
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))[:8]
}
