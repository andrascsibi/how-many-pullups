package accounts

import (
	"net/http"

	"github.com/gorilla/mux"

	"time"

	"crypto/sha1"
	"encoding/base64"
	"io"
	//    "errors"
	// "strconv"
	"fmt"
	//    "regexp"

	"appengine"
	"appengine/datastore"
	"appengine/user"
	"encoding/json"
	"errors"
)

type Account struct {
	Email      string
	ID         string
	ScreenName string // TODO
	Admin      bool
	RegDate    time.Time

	Challenges []string // TODO
	Settings   Settings // TODO
}

type Profile struct {
}
type Settings struct {
}

type Challenge struct {
	Title       string
	Description string
	Active      bool
}

func init() {
	r := mux.NewRouter()
	r.HandleFunc("/whoami", whoami)

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
		handler(updateAccount)).
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

	r.Handle("/accounts/{accountId}/challenges/{challengeId}/sets",
		handler(getSets)).
		Methods("GET")
	r.Handle("/accounts/{accountId}/challenges/{challengeId}/sets",
		handler(createSet)).
		Methods("POST")
	r.Handle("/accounts/{accountId}/challenges/{challengeId}/sets",
		handler(importSets)).
		Methods("PUT")

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
	q := datastore.NewQuery("Accounts").Order("-RegDate")
    as := make([]Account, 0)
    _, err := q.GetAll(c, &as)
    if err != nil {
        return nil, &handlerError{err, "Error querying datastore", http.StatusInternalServerError}
    }
	return as, nil
}

func createAccount(w http.ResponseWriter, r *http.Request) (interface{}, *handlerError) {
	c := appengine.NewContext(r)
	_ = c
	fmt.Fprintf(w, "creating challenge for account\n")
	return nil, nil
}

func getAccount(w http.ResponseWriter, r *http.Request) (interface{}, *handlerError) {
	c := appengine.NewContext(r)
	accountId := mux.Vars(r)["accountId"]

	var account Account
	err := datastore.Get(c, accountKey(c, accountId), &account)

	if err == datastore.ErrNoSuchEntity {
		return nil, &handlerError{err, "Account not found", http.StatusNotFound}
	}

	return account, nil
}

func updateAccount(w http.ResponseWriter, r *http.Request) (interface{}, *handlerError) {
	c := appengine.NewContext(r)
	accountId := mux.Vars(r)["accountId"]
	_ = c
	fmt.Fprintf(w, "updating account %v\n", accountId)
	return nil, nil
}

func getChallenges(w http.ResponseWriter, r *http.Request) (interface{}, *handlerError) {
	c := appengine.NewContext(r)
	accountId := mux.Vars(r)["accountId"]
	_ = c
	fmt.Fprintf(w, "listing challenges for account %v\n", accountId)
	return nil, nil
}

func createChallenge(w http.ResponseWriter, r *http.Request) (interface{}, *handlerError) {
	c := appengine.NewContext(r)
	accountId := mux.Vars(r)["accountId"]
	_ = c
	fmt.Fprintf(w, "creating challenge for account %v\n", accountId)
	return nil, nil
}

func getChallenge(w http.ResponseWriter, r *http.Request) (interface{}, *handlerError) {
	c := appengine.NewContext(r)
	accountId := mux.Vars(r)["accountId"]
	challengeId := mux.Vars(r)["challengeId"]
	_ = c
	fmt.Fprintf(w, "listing challenge %v/%v\n", accountId, challengeId)
	return nil, nil
}

func updateChallenge(w http.ResponseWriter, r *http.Request) (interface{}, *handlerError) {
	c := appengine.NewContext(r)
	accountId := mux.Vars(r)["accountId"]
	challengeId := mux.Vars(r)["challengeId"]
	_ = c
	fmt.Fprintf(w, "updating challenge title/descr %v/%v\n", accountId, challengeId)
	return nil, nil
}

func getSets(w http.ResponseWriter, r *http.Request) (interface{}, *handlerError) {
	c := appengine.NewContext(r)
	accountId := mux.Vars(r)["accountId"]
	challengeId := mux.Vars(r)["challengeId"]
	_ = c
	fmt.Fprintf(w, "getting sets of %v/%v\n", accountId, challengeId)
	return nil, nil
}

func createSet(w http.ResponseWriter, r *http.Request) (interface{}, *handlerError) {
	c := appengine.NewContext(r)
	accountId := mux.Vars(r)["accountId"]
	challengeId := mux.Vars(r)["challengeId"]
	_ = c
	fmt.Fprintf(w, "creating set for %v/%v\n", accountId, challengeId)
	return nil, nil
}

func importSets(w http.ResponseWriter, r *http.Request) (interface{}, *handlerError) {
	c := appengine.NewContext(r)
	accountId := mux.Vars(r)["accountId"]
	challengeId := mux.Vars(r)["challengeId"]
	_ = c
	fmt.Fprintf(w, "importing sets to %v/%v\n", accountId, challengeId)
	return nil, nil
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
		account = Account{
			Email:   email,
			ID:      id,
            Admin:   u.Admin,
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
	Account   Account
	LoginURL  string
	LogoutURL string
}

func hash(id string) string {
	hasher := sha1.New()
	io.WriteString(hasher, id)
	io.WriteString(hasher, "salt it real good DbqOFzkk") // TODO: should come from environment
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))[:8]
}
