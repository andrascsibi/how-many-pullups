package challenges

import (
	"github.com/gorilla/mux"
	"net/http"

	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/csv"
	"encoding/hex"
	"encoding/json"
	"unicode"

	"io"
	"io/ioutil"

	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"appengine"
	"appengine/datastore"
	"appengine/user"
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

func challengeKey(c appengine.Context, accountId string, id string) *datastore.Key {
	return datastore.NewKey(c, "Challenges", id, 0, accountKey(c, accountId))
}

type ChallengeStat struct {
	Today   int
	Total   int
	MinDate time.Time
	MaxDate time.Time
}

func init() {
	r := mux.NewRouter()

	r.Handle("/accounts/{accountId}/challenges",
		handler(getChallenges)).
		Methods("GET")
	r.Handle("/accounts/{accountId}/challenges",
		handler(createChallenge)).
		Methods("POST")

	r.Handle("/accounts/{accountId}/challenges/{challengeId}",
		handler(getChallenge)).
		Methods("GET")
	r.Handle("/accounts/{accountId}/challenges/{challengeId}",
		handler(updateChallenge)).
		Methods("POST")
}

func getChallenges(c appenginge.Context, w http.ResponseWriter, r *http.Request) (interface{}, *handler.Error) {
	accountId := mux.Vars(r)["accountId"]

	q := datastore.NewQuery("Challenges").
		Ancestor(accountKey(c, accountId)).
		Order("-CreationDate")
	challenges := make([]Challenge, 0)
	_, e := q.GetAll(c, &challenges)

	if e != nil {
		return nil, &handlerError{e, "Error querying datastore", http.StatusInternalServerError}
	}

	return challenges, nil
}

func createChallenge(c appenginge.Context, w http.ResponseWriter, r *http.Request) (interface{}, *handler.Error) {

	accountId := mux.Vars(r)["accountId"]
	account, err := getAccountById(c, accountId)
	if err != nil {
		return nil, err
	}

	authE := authorize(c, account)
	if authE != nil {
		return nil, authE
	}

	data, e := ioutil.ReadAll(r.Body)
	if e != nil {
		return nil, &handlerError{e, "Could not read request", http.StatusBadRequest}
	}

	var challenge Challenge
	e = json.Unmarshal(data, &challenge)
	if e != nil {
		return nil, &handlerError{e, "Could not parse JSON", http.StatusBadRequest}
	}

	challenge.CreationDate = time.Now()
	challenge.ID = hash(challenge.CreationDate.String())
	challenge.AccountID = accountId

	key := challengeKey(c, accountId, challenge.ID)
	_, e = datastore.Put(c, key, &challenge)
	if e != nil {
		return nil, &handlerError{e, "Error storing in datastore", http.StatusInternalServerError}
	}

	return challenge, nil
}

func getChallenge(c appenginge.Context, w http.ResponseWriter, r *http.Request) (interface{}, *handler.Error) {
	accountId := mux.Vars(r)["accountId"]
	challengeId := mux.Vars(r)["challengeId"]

	var challenge Challenge
	err := datastore.Get(c, challengeKey(c, accountId, challengeId), &challenge)

	if err == datastore.ErrNoSuchEntity {
		return nil, &handlerError{err, "Challenge not found", http.StatusNotFound}
	} else if err != nil {
		return nil, &handlerError{err, "Error accessing datastore", http.StatusInternalServerError}
	}

	return challenge, nil
}

func updateChallenge(c appenginge.Context, w http.ResponseWriter, r *http.Request) (interface{}, *handler.Error) {
	accountId := mux.Vars(r)["accountId"]
	challengeId := mux.Vars(r)["challengeId"]
	_ = c

	account, err := getAccountById(c, accountId)
	if err != nil {
		return nil, err
	}

	authE := authorize(c, account)
	if authE != nil {
		return nil, authE
	}

	challenge, err := getChallenge(w, r)
	if err != nil {
		return nil, err
	}
	ch := challenge.(Challenge)

	data, e := ioutil.ReadAll(r.Body)
	if e != nil {
		return nil, &handlerError{e, "Could not read request", http.StatusBadRequest}
	}

	var updatedChallenge Challenge
	e = json.Unmarshal(data, &updatedChallenge)
	if e != nil {
		return nil, &handlerError{e, "Could not parse JSON", http.StatusBadRequest}
	}

	// protect certain fields
	updatedChallenge.CreationDate = ch.CreationDate
	updatedChallenge.ID = ch.ID
	updatedChallenge.AccountID = ch.AccountID

	key := challengeKey(c, accountId, challengeId)
	_, e = datastore.Put(c, key, &updatedChallenge)
	if e != nil {
		return nil, &handlerError{e, "Error storing in datastore", http.StatusInternalServerError}
	}

	return updatedChallenge, nil
}
