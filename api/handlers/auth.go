package handlers

import (
	"encoding/json"
	"errors"
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

func GetToken(r *http.Request) (*uuid.UUID, error) {
	tokenCookie, err := r.Cookie("sessiontoken")
	if err != nil {
		return nil, errors.New("failed to read authentication token")
	}
	tokenString := tokenCookie.Value
	token, err := uuid.Parse(tokenString)
	if err != nil {
		return nil, errors.New("failed to read authentication token")
	}
	return &token, nil
}

func (a *AuthHandler) GetUser(r *http.Request) (*uuid.UUID, error) {
	token, err := GetToken(r)
	if err != nil {
		return nil, err
	}
	userId, err := a.DBService.GetLoggedInUser(*token)
	if err != nil {
		return nil, err
	}
	return userId, nil
}

func (a *AuthHandler) login(w http.ResponseWriter, r *http.Request) {
	var credentials db.Credentials
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("bad request"))
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
	token, err := GetToken(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(err.Error()))
	}
	a.DBService.InvalidateToken(*token)
}
