package account

import (
	"github.com/gorilla/mux"
	"net/http"

	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"unicode"

	"io"
	"io/ioutil"

	"errors"
	"fmt"
	"regexp"
	"time"

	"appengine"
	"appengine/datastore"
	"appengine/user"

	"github.com/andrascsibi/how-many-pullups/handler"
	"github.com/andrascsibi/how-many-pullups/stringset"
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

const kind = "Accounts"

func NewKey(c appengine.Context, id string) *datastore.Key {
	return datastore.NewKey(c, kind, id, 0, nil)
}

func RegisterHandlers(r *mux.Router) {

	r.Handle("/whoami", handler.New(whoami)).Methods("GET")

	r.Handle("/accounts",
		handler.New(getAccounts)).
		Methods("GET")
	r.Handle("/accounts",
		handler.New(createAccount)).
		Methods("POST")

	r.Handle("/accounts/{op:follow|unfollow}/{follower}/{followee}",
		handler.New(follow)).
		Methods("POST")

	r.Handle("/accounts/{accountId}",
		handler.New(getAccount)).
		Methods("GET")
	r.Handle("/accounts/{accountId}",
		handler.New(createAccount)). // TODO update?
		Methods("POST")

}

func ById(c appengine.Context, id string) (*Account, *handler.Error) {
	var account Account
	err := datastore.Get(c, NewKey(c, id), &account)

	if err == datastore.ErrNoSuchEntity {
		return nil, &handler.Error{err, "Account not found", http.StatusNotFound}
	} else if err != nil {
		return nil, &handler.Error{err, "Error getting account", http.StatusInternalServerError}
	}

	account.EmailMD5 = md5hex(account.Email)
	if account.ScreenName == "" {
		a := []rune(account.ID)
		a[0] = unicode.ToUpper(a[0])
		account.ScreenName = string(a)
	}
	return &account, nil
}

func (a *Account) Authorize(c appengine.Context) *handler.Error {
	u := user.Current(c)
	if u == nil {
		return &handler.Error{nil, "Login requried", http.StatusForbidden}
	}
	if u.Admin {
		return nil
	}
	// XXX
	if a == nil || u.Email != a.Email {
		return &handler.Error{nil, "Unauthorized", http.StatusUnauthorized}
	}
	return nil
}

func isAdmin(c appengine.Context) bool {
	u := user.Current(c)
	if u == nil {
		return false
	}
	return u.Admin
}

func getAccounts(c appengine.Context, w http.ResponseWriter, r *http.Request, v map[string]string) (interface{}, *handler.Error) {

	if !isAdmin(c) {
		return nil, &handler.Error{nil, "only admins have access to this", http.StatusForbidden}
	}

	q := datastore.NewQuery(kind).Order("-RegDate").Limit(100)
	as := make([]Account, 0)
	_, err := q.GetAll(c, &as)
	if err != nil {
		return nil, &handler.Error{err, "Error querying datastore", http.StatusInternalServerError}
	}
	return as, nil
}

func createAccount(c appengine.Context, w http.ResponseWriter, r *http.Request, v map[string]string) (interface{}, *handler.Error) {

	data, e := ioutil.ReadAll(r.Body)
	if e != nil {
		return nil, &handler.Error{e, "Could not read request", http.StatusBadRequest}
	}

	var newAccount Account
	e = json.Unmarshal(data, &newAccount)
	if e != nil {
		return nil, &handler.Error{e, "Could not parse JSON", http.StatusBadRequest}
	}

	authE := newAccount.Authorize(c)
	if authE != nil {
		return nil, authE
	}

	accByEmail, err := getAccountByEmail(c, newAccount.Email)
	if err != nil {
		return nil, &handler.Error{e, "Error getting account by email", http.StatusInternalServerError}
	}
	if accByEmail != nil {
		return nil, &handler.Error{e, "An account is already registered for this email", http.StatusConflict}
	}

	if e = validate(newAccount.ID); e != nil {
		return nil, &handler.Error{e, e.Error(), http.StatusBadRequest}
	}

	key := NewKey(c, newAccount.ID)

	var accInDb Account
	err = datastore.Get(c, key, &accInDb)

	if err == datastore.ErrNoSuchEntity {
		c.Infof("creating new user: %v %v", newAccount.Email, newAccount.ID)
		newAccount.RegDate = time.Now()
		_, err = datastore.Put(c, key, &newAccount)
		if err != nil {
			return nil, &handler.Error{e, "Error storing in datastore", http.StatusInternalServerError}
		}
		return newAccount, nil
	}
	if err != nil {
		return nil, &handler.Error{e, "Error accessing datastore", http.StatusInternalServerError}
	}

	return nil, &handler.Error{e, "User already exists", http.StatusConflict}
}

func validate(username string) error {
	matches, err := regexp.Match("^[a-z][a-z0-9\\-]{2,29}$", []byte(username))
	if err != nil {
		return err
	}
	if !matches {
		return errors.New("Username must be between 3 and 30 characters long, must start with a lowercase letter, and can only contain lowercase letters, numbers, and the '-' character.")
	}
	if stringset.IndexOf([]string{"admin", "follow", "unfollow"}, username) >= 0 {
		return errors.New("Invalid username.")
	}
	return nil
}

func getAccount(c appengine.Context, w http.ResponseWriter, r *http.Request, v map[string]string) (interface{}, *handler.Error) {
	accountId := v["accountId"]

	account, err := ById(c, accountId)
	if err != nil {
		return nil, err
	}

	account.Email = ""
	return *account, nil
}

func follow(c appengine.Context, w http.ResponseWriter, r *http.Request, v map[string]string) (interface{}, *handler.Error) {
	follower := v["follower"]
	followee := v["followee"]
	unfollow := v["op"] == "unfollow"

	followerA, err := ById(c, follower)
	if err != nil {
		return nil, err
	}

	followeeA, err := ById(c, followee)
	if err != nil {
		return nil, err
	}

	authE := followerA.Authorize(c)
	if authE != nil {
		return nil, authE
	}

	if unfollow {
		followerA.Following = stringset.Remove(followerA.Following, followee)
		followeeA.Followers = stringset.Remove(followeeA.Followers, follower)
	} else {
		followerA.Following = stringset.Add(followerA.Following, followee)
		followeeA.Followers = stringset.Add(followeeA.Followers, follower)
	}

	trErr := datastore.RunInTransaction(c, func(c appengine.Context) error {
		_, err := datastore.Put(c, NewKey(c, follower), followerA)
		if err != nil {
			return err
		}
		_, err = datastore.Put(c, NewKey(c, followee), followeeA)
		if err != nil {
			return err
		}
		return nil
	}, &datastore.TransactionOptions{XG: true})

	if trErr != nil {
		return nil, &handler.Error{trErr, "could not set relationship", http.StatusInternalServerError}
	}

	return followerA, nil
}

func updateAccount(c appengine.Context, w http.ResponseWriter, r *http.Request, v map[string]string) (interface{}, *handler.Error) {
	accountId := v["accountId"]
	_ = accountId
	return nil, &handler.Error{errors.New("updating account not supported"), "", http.StatusMethodNotAllowed}
}

func getAccountByEmail(c appengine.Context, email string) (*Account, error) {
	q := datastore.NewQuery(kind).Filter("Email =", email)

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

func whoami(c appengine.Context, w http.ResponseWriter, r *http.Request, v map[string]string) (interface{}, *handler.Error) {
	u := user.Current(c)

	if u == nil {
		url, err := user.LoginURL(c, "?redir=true")
		if err != nil {
			return nil, &handler.Error{err, "Error getting login URL", http.StatusInternalServerError}
		}
		return LoginData{LoginURL: url}, nil
	}

	url, err := user.LogoutURL(c, "")
	if err != nil {
		return nil, &handler.Error{err, "Error getting logout URL", http.StatusInternalServerError}
	}

	account, err := getAccountByEmail(c, u.Email)
	if err != nil {
		return nil, &handler.Error{err, "Error while getting account", http.StatusInternalServerError}
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

func md5hex(src string) string {
	hasher := md5.New()
	io.WriteString(hasher, src)
	return hex.EncodeToString(hasher.Sum(nil))
}
