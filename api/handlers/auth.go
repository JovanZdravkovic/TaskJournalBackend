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
	login  = regexp.MustCompile(`^/auth/login/*$`)
	logout = regexp.MustCompile(`^/auth/logout/*$`)
)

func (a *AuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodPost && login.MatchString(r.URL.Path):
		a.login(w, r)
		return
	case r.Method == http.MethodPost && logout.MatchString(r.URL.Path):
		a.logout(w, r)
		return
	default:
		return
	}
}

func (a *AuthHandler) login(w http.ResponseWriter, r *http.Request) {
	var credentials db.Credentials
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad request"))
		return
	}
	authRow, err := a.DBService.CreateToken(credentials)
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
}
func (a *AuthHandler) logout(w http.ResponseWriter, r *http.Request) {
	tokenCookie, err := r.Cookie("sessiontoken")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Failed to read authentication token"))
		return
	}
	tokenString := tokenCookie.Name
	token, err := uuid.Parse(tokenString)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Failed to read authentication token"))
		return
	}
	a.DBService.InvalidateToken(token)
}
