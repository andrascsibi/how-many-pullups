package app

import (
	"github.com/andrascsibi/how-many-pullups/account"
	"github.com/andrascsibi/how-many-pullups/challenge"

  "github.com/gorilla/mux"
  "net/http"

)


func init() {

  r := mux.NewRouter()
  account.RegisterHandlers(r)
  challenge.RegisterHandlers(r)

  http.Handle("/", r)
}
