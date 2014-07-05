package accounts

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

type Account struct {
	ID         string
	Email      string
	EmailMD5   string
	ScreenName string
	RegDate    time.Time

	Following []string
	Followers []string
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
	MaxReps      int
	StepReps     int
	CreationDate time.Time
}

type WorkSet struct {
	Reps int
	Date time.Time
}

type ChallengeStat struct {
	Today   int
	Total   int
	MinDate time.Time
	MaxDate time.Time
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

	r.Handle("/accounts/{op:follow|unfollow}/{follower}/{followee}",
		handler(follow)).
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

	r.Handle("/accounts/{accountId}/challenges/{challengeId}/export-csv",
		handler(exportCsv)).
		Methods("GET")
	r.Handle("/accounts/{accountId}/challenges/{challengeId}/export",
		handler(export)).
		Methods("GET")
	r.Handle("/accounts/{accountId}/challenges/{challengeId}/import-csv",
		handler(importCsv)).
		Methods("POST")

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

	authE := authorize(c, nil)
	if authE != nil {
		return nil, authE
	}

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

	authE := authorize(c, &newAccount)
	if authE != nil {
		return nil, authE
	}

	accByEmail, err := getAccountByEmail(c, newAccount.Email)
	if err != nil {
		return nil, &handlerError{e, "Error getting account by email", http.StatusInternalServerError}
	}
	if accByEmail != nil {
		return nil, &handlerError{e, "An account is already registered for this email", http.StatusConflict}
	}

	if e = validate(newAccount.ID); e != nil {
		return nil, &handlerError{e, e.Error(), http.StatusBadRequest}
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
	if indexOf([]string{"admin", "follow", "unfollow"}, username) >= 0 {
		return errors.New("Invalid username.")
	}
	return nil
}

func getAccount(w http.ResponseWriter, r *http.Request) (interface{}, *handlerError) {
	c := appengine.NewContext(r)
	accountId := mux.Vars(r)["accountId"]

	account, err := getAccountById(c, accountId)
	if err != nil {
		return nil, err
	}

	account.Email = ""
	return account, nil
}

func indexOf(list []string, item string) int {
	for i, s := range list {
		if s == item {
			return i
		}
	}
	return -1
}
func add(list []string, item string) []string {
	i := indexOf(list, item)
	if i >= 0 {
		return list
	}
	return append(list, item)
}
func remove(list []string, item string) []string {
	i := indexOf(list, item)
	if i < 0 {
		return list
	}
	return append(list[:i], list[i+1:]...)
}

func follow(w http.ResponseWriter, r *http.Request) (interface{}, *handlerError) {
	c := appengine.NewContext(r)
	follower := mux.Vars(r)["follower"]
	followee := mux.Vars(r)["followee"]
	unfollow := mux.Vars(r)["op"] == "unfollow"

	followerA, err := getAccountById(c, follower)
	if err != nil {
		return nil, err
	}

	followeeA, err := getAccountById(c, followee)
	if err != nil {
		return nil, err
	}

	authE := authorize(c, followerA)
	if authE != nil {
		return nil, authE
	}

	if unfollow {
		followerA.Following = remove(followerA.Following, followee)
		followeeA.Followers = remove(followeeA.Followers, follower)
	} else {
		followerA.Following = add(followerA.Following, followee)
		followeeA.Followers = add(followeeA.Followers, follower)
	}

	trErr := datastore.RunInTransaction(c, func(c appengine.Context) error {
		_, err := datastore.Put(c, accountKey(c, follower), followerA)
		if err != nil {
			return err
		}
		_, err = datastore.Put(c, accountKey(c, followee), followeeA)
		if err != nil {
			return err
		}
		return nil
	}, &datastore.TransactionOptions{XG: true})

	if trErr != nil {
		return nil, &handlerError{trErr, "could not set relationship", http.StatusInternalServerError}
	}

	return followerA, nil
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

func export(w http.ResponseWriter, r *http.Request) (interface{}, *handlerError) {
	c := appengine.NewContext(r)
	accountId := mux.Vars(r)["accountId"]
	challengeId := mux.Vars(r)["challengeId"]

	q := datastore.NewQuery("WorkSets").Ancestor(challengeKey(c, accountId, challengeId)).Order("-Date")
	sets := make([]WorkSet, 0)
	_, err := q.GetAll(c, &sets)

	if err != nil {
		return nil, &handlerError{err, "Error reading sets", http.StatusInternalServerError}
	}

	return sets, nil
}

func exportCsv(w http.ResponseWriter, r *http.Request) (interface{}, *handlerError) {
	s, err := export(w, r)
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

func importCsv(w http.ResponseWriter, r *http.Request) (interface{}, *handlerError) {
	c := appengine.NewContext(r)
	accountId := mux.Vars(r)["accountId"]
	challengeId := mux.Vars(r)["challengeId"]

	account, aerr := getAccountById(c, accountId)
	if aerr != nil {
		return nil, aerr
	}

	authE := authorize(c, account)
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
			return nil, &handlerError{err, "Could not read request", http.StatusBadRequest}
		}

		if len(line) != 2 {
			return nil, &handlerError{err, fmt.Sprintf("Each line should contain 2 fields. Line no: %d '%v'", i, line), http.StatusBadRequest}
		}
		if line[0] == "timestamp" && line[1] == "reps" {
			continue
		}
		date, err := time.Parse(time.RFC3339, line[0])
		if err != nil {
			return nil, &handlerError{err, fmt.Sprintf("Malformed date in line: %d '%v'", i, line[0]), http.StatusBadRequest}
		}
		reps, err := strconv.Atoi(line[1])
		if err != nil {
			return nil, &handlerError{err, fmt.Sprintf("Malformed number in line: %d '%v'", i, line[1]), http.StatusBadRequest}
		}
		if reps == 0 {
			continue
		}
		if i > 1000 {
			return nil, &handlerError{err, "Too many lines", http.StatusBadRequest}
		}
		importedSets = append(importedSets, WorkSet{Date: date, Reps: reps})
	}

	keys := make([]*datastore.Key, len(importedSets))
	for i := 0; i < len(importedSets); i++ {
		keys[i] = workSetKey(c, accountId, challengeId)
	}
	_, err := datastore.PutMulti(c, keys, importedSets)
	if err != nil {
		return nil, &handlerError{err, "Error storing in datastore", http.StatusInternalServerError}
	}

	return nil, nil
}

// TODO only works for utc
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

func accountKey(c appengine.Context, id string) *datastore.Key {
	return datastore.NewKey(c, "Accounts", id, 0, nil)
}

func challengeKey(c appengine.Context, accountId string, id string) *datastore.Key {
	return datastore.NewKey(c, "Challenges", id, 0, accountKey(c, accountId))
}

func workSetKey(c appengine.Context, accountId string, challengeId string) *datastore.Key {
	return datastore.NewIncompleteKey(c, "WorkSets", challengeKey(c, accountId, challengeId))
}

func authorize(c appengine.Context, a *Account) *handlerError {
	u := user.Current(c)
	if u == nil {
		return &handlerError{nil, "Login requried", http.StatusForbidden}
	}
	if u.Admin {
		return nil
	}
	if a == nil || u.Email != a.Email {
		return &handlerError{nil, "Unauthorized", http.StatusUnauthorized}
	}
	return nil
}

func getAccountById(c appengine.Context, id string) (*Account, *handlerError) {
	var account Account
	err := datastore.Get(c, accountKey(c, id), &account)

	if err == datastore.ErrNoSuchEntity {
		return nil, &handlerError{err, "Account not found", http.StatusNotFound}
	} else if err != nil {
		return nil, &handlerError{err, "Error getting account", http.StatusInternalServerError}
	}

	account.EmailMD5 = md5hex(account.Email)
	if account.ScreenName == "" {
		a := []rune(account.ID)
		a[0] = unicode.ToUpper(a[0])
		account.ScreenName = string(a)
	}
	return &account, nil
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
		url, err := user.LoginURL(c, "?redir=true")
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

func md5hex(src string) string {
	hasher := md5.New()
	io.WriteString(hasher, src)
	return hex.EncodeToString(hasher.Sum(nil))
}
