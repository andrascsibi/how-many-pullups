package accounts

import (
	"net/http"

	"github.com/gorilla/mux"

	"time"

	"crypto/sha1"
	"encoding/base64"
	"io"
	"io/ioutil"
	// "strconv"
	"fmt"
	"regexp"

	"appengine"
	"appengine/datastore"
	"appengine/user"
	"encoding/json"
	"errors"
)

type Account struct {
	Email   string
	ID      string
	RegDate time.Time
}

type Profile struct {
}
type Settings struct {
}

type Challenge struct {
	AccountID    string
	ID           string
	Title        string
	Description  string
	CreationDate time.Time
}

type WorkSet struct {
	Reps int
	Date time.Time
}

type ChallengeStat struct {
	Today int
	Total int
}

func init() {
	r := mux.NewRouter()
	r.Handle("/whoami", handler(whoami)).Methods("GET")

	r.Handle("/accounts",
		handler(getAccounts)).
		Methods("GET")
	r.Handle("/accounts",
		handler(createAccount)).
		Methods("POST")

	r.Handle("/accounts/{accountId}",
		handler(getAccount)).
		Methods("GET")
	r.Handle("/accounts/{accountId}",
		handler(createAccount)). // TODO update?
		Methods("POST")

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

	r.Handle("/accounts/{accountId}/challenges/{challengeId}/stats",
		handler(getStats)).
		Methods("GET")

	r.Handle("/accounts/{accountId}/challenges/{challengeId}/sets",
		handler(getSets)).
		Methods("GET")
	r.Handle("/accounts/{accountId}/challenges/{challengeId}/sets",
		handler(createSet)).
		Methods("POST")
	r.Handle("/accounts/{accountId}/challenges/{challengeId}/import",
		handler(importSets)).
		Methods("PUT")
	r.Handle("/accounts/{accountId}/challenges/{challengeId}/export",
		handler(exportSets)).
		Methods("GET")

	http.Handle("/", r)
}

type handlerError struct {
	Error   error
	Message string
	Code    int
}

type handler func(w http.ResponseWriter, r *http.Request) (interface{}, *handlerError)

// handler implements the http.Handler interface
func (fn handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	c := appengine.NewContext(r)
	response, err := fn(w, r)

	if err != nil {
		c.Errorf("%v", err.Error)
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Message), err.Code)
		return
	}
	if response == nil {
		c.Errorf("response from method is nil")
		http.Error(w, "Internal server error. Check the logs.", http.StatusInternalServerError)
		return
	}

	bytes, e := json.Marshal(response)
	if e != nil {
		http.Error(w, "Error marshalling JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(bytes)
}

func getAccounts(w http.ResponseWriter, r *http.Request) (interface{}, *handlerError) {
	c := appengine.NewContext(r)
	// TODO authorization
	q := datastore.NewQuery("Accounts").Order("-RegDate").Limit(100)
	as := make([]Account, 0)
	_, err := q.GetAll(c, &as)
	if err != nil {
		return nil, &handlerError{err, "Error querying datastore", http.StatusInternalServerError}
	}
	return as, nil
}

func createAccount(w http.ResponseWriter, r *http.Request) (interface{}, *handlerError) {
	c := appengine.NewContext(r)

	data, e := ioutil.ReadAll(r.Body)
	if e != nil {
		return nil, &handlerError{e, "Could not read request", http.StatusBadRequest}
	}

	var newAccount Account
	e = json.Unmarshal(data, &newAccount)
	if e != nil {
		return nil, &handlerError{e, "Could not parse JSON", http.StatusBadRequest}
	}

	u := user.Current(c)
	if u == nil {
		return nil, &handlerError{e, "Login requried", http.StatusForbidden}
	}
	if u.Email != newAccount.Email && u.Admin {
		return nil, &handlerError{e, "Unauthorized", http.StatusUnauthorized}
	}

	accByEmail, err := getAccountByEmail(c, newAccount.Email)
	if err != nil {
		return nil, &handlerError{e, "Error getting account by email", http.StatusInternalServerError}
	}
	if accByEmail != nil {
		return nil, &handlerError{e, "An account is already registered for this email", http.StatusConflict}
	}

	if e = validate(newAccount.ID); e != nil {
		return nil, &handlerError{e, "Invalid user name", http.StatusBadRequest}
	}

	key := accountKey(c, newAccount.ID)

	var accInDb Account
	err = datastore.Get(c, key, &accInDb)

	if err == datastore.ErrNoSuchEntity {
		c.Infof("creating new user: %v %v", newAccount.Email, newAccount.ID)
		newAccount.RegDate = time.Now()
		_, err = datastore.Put(c, key, &newAccount)
		if err != nil {
			return nil, &handlerError{e, "Error storing in datastore", http.StatusInternalServerError}
		}
		return newAccount, nil
	}
	if err != nil {
		return nil, &handlerError{e, "Error accessing datastore", http.StatusInternalServerError}
	}

	return nil, &handlerError{e, "User already exists", http.StatusConflict}
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

func getAccount(w http.ResponseWriter, r *http.Request) (interface{}, *handlerError) {
	c := appengine.NewContext(r)
	accountId := mux.Vars(r)["accountId"]

	var account Account
	err := datastore.Get(c, accountKey(c, accountId), &account)

	if err == datastore.ErrNoSuchEntity {
		return nil, &handlerError{err, "Account not found", http.StatusNotFound}
	} else if err != nil {
		return nil, &handlerError{err, "Error getting account", http.StatusInternalServerError}
	}

	return account, nil
}

func updateAccount(w http.ResponseWriter, r *http.Request) (interface{}, *handlerError) {
	c := appengine.NewContext(r)
	accountId := mux.Vars(r)["accountId"]
	_ = c
	_ = accountId
	return nil, &handlerError{errors.New("updating account not supported"), "", http.StatusMethodNotAllowed}
}

func getChallenges(w http.ResponseWriter, r *http.Request) (interface{}, *handlerError) {
	c := appengine.NewContext(r)
	accountId := mux.Vars(r)["accountId"]

	// TODO: check account exists

	// account, err := getAccount(w, r)
	// if err != nil {
	// 	return nil, err
	// }

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

func createChallenge(w http.ResponseWriter, r *http.Request) (interface{}, *handlerError) {
	c := appengine.NewContext(r)

	// TODO: authorization
	// TODO: check account exists
	accountId := mux.Vars(r)["accountId"]

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

func getChallenge(w http.ResponseWriter, r *http.Request) (interface{}, *handlerError) {
	c := appengine.NewContext(r)
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

func updateChallenge(w http.ResponseWriter, r *http.Request) (interface{}, *handlerError) {
	c := appengine.NewContext(r)
	accountId := mux.Vars(r)["accountId"]
	challengeId := mux.Vars(r)["challengeId"]
	_ = c

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

	ch.Title = updatedChallenge.Title
	ch.Description = updatedChallenge.Description

	key := challengeKey(c, accountId, challengeId)
	_, e = datastore.Put(c, key, &ch)
	if e != nil {
		return nil, &handlerError{e, "Error storing in datastore", http.StatusInternalServerError}
	}

	return challenge, nil
}

func getStats(w http.ResponseWriter, r *http.Request) (interface{}, *handlerError) {
	c := appengine.NewContext(r)
	accountId := mux.Vars(r)["accountId"]
	challengeId := mux.Vars(r)["challengeId"]

	q := datastore.NewQuery("WorkSets").Ancestor(challengeKey(c, accountId, challengeId)).Order("-Date")
	sets := make([]WorkSet, 0)
	_, err := q.GetAll(c, &sets)

	if err != nil {
		return nil, &handlerError{err, "Error reading sets", http.StatusInternalServerError}
	}

	var stat ChallengeStat
	today := time.Now().Truncate(24 * time.Hour)
	for _, s := range sets {
		if today.Equal(s.Date.Truncate(24 * time.Hour)) {
			stat.Today += s.Reps
		}
		stat.Total += s.Reps
	}

	return stat, nil
}

func getSets(w http.ResponseWriter, r *http.Request) (interface{}, *handlerError) {
	c := appengine.NewContext(r)
	accountId := mux.Vars(r)["accountId"]
	challengeId := mux.Vars(r)["challengeId"]
	_ = c
	_ = accountId
	_ = challengeId
	return nil, &handlerError{errors.New("import not supported"), "", http.StatusMethodNotAllowed}
}

func createSet(w http.ResponseWriter, r *http.Request) (interface{}, *handlerError) {
	c := appengine.NewContext(r)
	accountId := mux.Vars(r)["accountId"]
	challengeId := mux.Vars(r)["challengeId"]

	data, e := ioutil.ReadAll(r.Body)
	if e != nil {
		return nil, &handlerError{e, "Could not read request", http.StatusBadRequest}
	}

	var newSet WorkSet
	e = json.Unmarshal(data, &newSet)
	if e != nil {
		return nil, &handlerError{e, "Could not parse JSON", http.StatusBadRequest}
	}

	newSet.Date = time.Now()

	key := workSetKey(c, accountId, challengeId)
	_, e = datastore.Put(c, key, &newSet)
	if e != nil {
		return nil, &handlerError{e, "Error storing in datastore", http.StatusInternalServerError}
	}

	return newSet, nil
}

func importSets(w http.ResponseWriter, r *http.Request) (interface{}, *handlerError) {
	c := appengine.NewContext(r)
	accountId := mux.Vars(r)["accountId"]
	challengeId := mux.Vars(r)["challengeId"]
	_ = c
	_ = accountId
	_ = challengeId
	return nil, &handlerError{errors.New("import not supported"), "", http.StatusMethodNotAllowed}
}

func exportSets(w http.ResponseWriter, r *http.Request) (interface{}, *handlerError) {
	c := appengine.NewContext(r)
	accountId := mux.Vars(r)["accountId"]
	challengeId := mux.Vars(r)["challengeId"]
	_ = c
	_ = accountId
	_ = challengeId
	return nil, &handlerError{errors.New("export not supported"), "", http.StatusMethodNotAllowed}
}

func accountKey(c appengine.Context, id string) *datastore.Key {
	return datastore.NewKey(c, "Accounts", id, 0, nil)
}

func challengeKey(c appengine.Context, accountId string, id string) *datastore.Key {
	return datastore.NewKey(c, "Challenges", id, 0, accountKey(c, accountId))
}

func workSetKey(c appengine.Context, accountId string, challengeId string) *datastore.Key {
	return datastore.NewIncompleteKey(c, "WorkSets", challengeKey(c, accountId, challengeId))
}

func getAccountByEmail(c appengine.Context, email string) (*Account, error) {
	q := datastore.NewQuery("Accounts").Filter("Email =", email)

	var accounts []Account
	_, err := q.GetAll(c, &accounts)
	if err != nil {
		c.Errorf(err.Error())
		return nil, err
	}
	switch len(accounts) {
	case 0:
		return nil, nil
	case 1:
		return &accounts[0], nil
	default:
		return nil, errors.New(fmt.Sprintf("More than one accounts found with email %v", email))
	}
}

func whoami(w http.ResponseWriter, r *http.Request) (interface{}, *handlerError) {

	c := appengine.NewContext(r)
	u := user.Current(c)

	if u == nil {
		url, err := user.LoginURL(c, "") // TODO: redirect?
		if err != nil {
			return nil, &handlerError{err, "Error getting login URL", http.StatusInternalServerError}
		}
		return LoginData{LoginURL: url}, nil
	}

	url, err := user.LogoutURL(c, "")
	if err != nil {
		return nil, &handlerError{err, "Error getting logout URL", http.StatusInternalServerError}
	}

	account, err := getAccountByEmail(c, u.Email)
	if err != nil {
		return nil, &handlerError{err, "Error while getting account", http.StatusInternalServerError}
	}
	if account == nil {
		return LoginData{LogoutURL: url, Unregistered: true, Account: Account{Email: u.Email}}, nil
	}

	return LoginData{LogoutURL: url, Account: *account}, nil
}

type LoginData struct {
	Account      Account
	Unregistered bool
	LoginURL     string
	LogoutURL    string
}

func hash(id string) string {
	hasher := sha1.New()
	io.WriteString(hasher, id)
	io.WriteString(hasher, "salt it real good DbqOFzkk") // TODO: should come from environment
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))[:8]
}
