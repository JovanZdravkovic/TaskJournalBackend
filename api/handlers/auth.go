package handlers

import (
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/JovanZdravkovic/TaskJournalBackend/db"
	"github.com/google/uuid"
)

type AuthHandler struct {
	DBService *db.DatabaseService
}

var (
	authenticate = regexp.MustCompile(`^/auth/*$`)
)

func (a *AuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	preflight := EnableCORS(w, r)
	if preflight {
		return
	}

	tokenString := r.Header.Get("X-Auth-Token")
	if tokenString == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	token, err := uuid.Parse(tokenString)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	switch {
	case r.Method == http.MethodGet && authenticate.MatchString(r.URL.Path):
		a.Authenticate(w, r, token)
		return
	default:
		return
	}
}

func (a *AuthHandler) Authenticate(w http.ResponseWriter, r *http.Request, userId uuid.UUID) {
	userJson, err := json.Marshal(db.Id{Id: userId})
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(userJson)
}
