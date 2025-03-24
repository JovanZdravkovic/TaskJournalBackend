package handlers

import (
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/JovanZdravkovic/TaskJournalBackend/db"
)

type SignupHandler struct {
	DBService *db.DatabaseService
}

var (
	signup = regexp.MustCompile(`^/signup/*$`)
)

func (signupHandler *SignupHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	preflight := EnableCORS(w, r)
	if preflight {
		return
	}

	switch {
	case r.Method == http.MethodPost && signup.MatchString(r.URL.Path):
		signupHandler.CreateUser(w, r)
		return
	default:
		return
	}
}

func (signupHandler *SignupHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user db.UserPost
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("bad request"))
		return
	}
	userId, err := signupHandler.DBService.CreateUser(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	userIdJson, err := json.Marshal(db.Id{Id: *userId})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error while constructing json"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(userIdJson)
}
