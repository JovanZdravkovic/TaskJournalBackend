package handlers

import (
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/JovanZdravkovic/TaskJournalBackend/db"
	"github.com/google/uuid"
)

var (
	Users = regexp.MustCompile(`^/user/*$`)
)

type UserHandler struct {
	DBService   *db.DatabaseService
	AuthService *AuthHandler
}

func (u *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	preflight := EnableCORS(w, r)
	if preflight {
		return
	}

	userId, err := u.AuthService.GetUser(r)
	switch {
	case r.Method == http.MethodGet && Users.MatchString(r.URL.Path):
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(err.Error()))
			return
		}
		u.GetUser(w, r, *userId)
		return
	case r.Method == http.MethodPost && Users.MatchString(r.URL.Path):
		u.CreateUser(w, r)
		return
	default:
		return
	}
}

func (u *UserHandler) GetUser(w http.ResponseWriter, r *http.Request, userId uuid.UUID) {
	user, err := u.DBService.GetUserInfo(userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	userJson, err := json.Marshal(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(userJson)
}

func (u *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user db.UserPost
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("bad request"))
		return
	}
	userId, err := u.DBService.CreateUser(user)
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
