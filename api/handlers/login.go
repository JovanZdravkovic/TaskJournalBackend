package handlers

import (
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/JovanZdravkovic/TaskJournalBackend/db"
)

type LoginHandler struct {
	DBService *db.DatabaseService
}

var (
	login = regexp.MustCompile(`^/login/*$`)
)

func (loginHandler *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	preflight := EnableCORS(w, r)
	if preflight {
		return
	}

	switch {
	case r.Method == http.MethodPost && login.MatchString(r.URL.Path):
		loginHandler.Login(w, r)
		return
	default:
		return
	}
}

func (loginHandler *LoginHandler) Login(w http.ResponseWriter, r *http.Request) {
	var credentials db.Credentials
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("bad request"))
		return
	}

	authRow, err := loginHandler.DBService.CreateToken(credentials)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(err.Error()))
		return
	}

	cookie := http.Cookie{
		Name:     "sessiontoken",
		Value:    authRow.Id.String(),
		Path:     "/",
		Expires:  authRow.ExpiresAt,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   true}
	http.SetCookie(w, &cookie)
	w.WriteHeader(http.StatusNoContent)
}
