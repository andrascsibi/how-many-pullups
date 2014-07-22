package handler

import (
	"net/http"

	"encoding/json"

	"appengine"
)

type Error struct {
	Error   error
	Message string
	Code    int
}

type handler struct {
	hf handlerFun
}

func New(hf handlerFun) http.Handler {
	return Handler
}

type handlerFun func(w http.ResponseWriter, r *http.Request) (interface{}, *Error)

// Handler implements the http.Handler interface
func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	c := appengine.NewContext(r)
	response, err := h.hf(w, r)

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
