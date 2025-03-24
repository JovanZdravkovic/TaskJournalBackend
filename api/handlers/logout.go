package handlers

import (
	"net/http"
	"regexp"

	"github.com/JovanZdravkovic/TaskJournalBackend/db"
)

type LogoutHandler struct {
	DBService *db.DatabaseService
}

var (
	logout = regexp.MustCompile(`^/logout/*$`)
)

func (logoutHandler *LogoutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	preflight := EnableCORS(w, r)
	if preflight {
		return
	}

	switch {
	case r.Method == http.MethodPost && logout.MatchString(r.URL.Path):
		logoutHandler.Logout(w, r)
		return
	default:
		return
	}
}

func (logoutHandler *LogoutHandler) Logout(w http.ResponseWriter, r *http.Request) {
	token, err := GetToken(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(err.Error()))
	}
	logoutHandler.DBService.InvalidateToken(*token)
	w.WriteHeader(http.StatusNoContent)
}
