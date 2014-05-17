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
	ScreenName string
	Admin      bool
	RegDate    time.Time

	Challenges []string
	Settings   Settings
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

	r.HandleFunc("/accounts/",
		createAccount).
		Methods("POST")

	r.HandleFunc("/accounts/{accountId}",
		getAccount).
		Methods("GET")
	r.HandleFunc("/accounts/{accountId}",
		updateAccount).
		Methods("POST")

	r.HandleFunc("/accounts/{accountId}/challenges",
		getChallenges).
		Methods("GET")
	r.HandleFunc("/accounts/{accountId}/challenges",
		createChallenge).
		Methods("POST")

	r.HandleFunc("/accounts/{accountId}/challenges/{challengeId}",
		getChallenge).
		Methods("GET")
	r.HandleFunc("/accounts/{accountId}/challenges/{challengeId}",
		updateChallenge).
		Methods("POST")

	r.HandleFunc("/accounts/{accountId}/challenges/{challengeId}/sets",
		getSets).
		Methods("GET")
	r.HandleFunc("/accounts/{accountId}/challenges/{challengeId}/sets",
		createSet).
		Methods("POST")
	r.HandleFunc("/accounts/{accountId}/challenges/{challengeId}/sets",
		importSets).
		Methods("PUT")

	http.Handle("/", r)
}

func createAccount(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	c.Infof("creating account")
}

func getAccount(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	accountId := mux.Vars(r)["accountId"]
	c.Infof("get account %v", accountId)

	var account Account
	err := datastore.Get(c, accountKey(c, accountId), &account)

	if err == datastore.ErrNoSuchEntity {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-type", "application/json")

	ret, err := json.Marshal(Account{
		Email: account.Email,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write(ret)
}

func updateAccount(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	accountId := mux.Vars(r)["accountId"]
	_ = c
	fmt.Fprintf(w, "updating account %v\n", accountId)
}

func getChallenges(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	accountId := mux.Vars(r)["accountId"]
	_ = c
	fmt.Fprintf(w, "listing challenges for account %v\n", accountId)
}

func createChallenge(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	accountId := mux.Vars(r)["accountId"]
	_ = c
	fmt.Fprintf(w, "creating challenge for account %v\n", accountId)
}

func getChallenge(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	accountId := mux.Vars(r)["accountId"]
	challengeId := mux.Vars(r)["challengeId"]
	_ = c
	fmt.Fprintf(w, "listing challenge %v/%v\n", accountId, challengeId)
}

func updateChallenge(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	accountId := mux.Vars(r)["accountId"]
	challengeId := mux.Vars(r)["challengeId"]
	_ = c
	fmt.Fprintf(w, "updating challenge title/descr %v/%v\n", accountId, challengeId)
}

func getSets(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	accountId := mux.Vars(r)["accountId"]
	challengeId := mux.Vars(r)["challengeId"]
	_ = c
	fmt.Fprintf(w, "getting sets of %v/%v\n", accountId, challengeId)
}

func createSet(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	accountId := mux.Vars(r)["accountId"]
	challengeId := mux.Vars(r)["challengeId"]
	_ = c
	fmt.Fprintf(w, "creating set for %v/%v\n", accountId, challengeId)
}

func importSets(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	accountId := mux.Vars(r)["accountId"]
	challengeId := mux.Vars(r)["challengeId"]
	_ = c
	fmt.Fprintf(w, "importing sets to %v/%v\n", accountId, challengeId)
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
