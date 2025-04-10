package handlers

import (
	"errors"
	"net/http"

	"github.com/JovanZdravkovic/TaskJournalBackend/db"
	"github.com/google/uuid"
)

func EnableCORS(w http.ResponseWriter, r *http.Request) bool {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return true
	}
	return false
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

func GetUser(r *http.Request, dbService db.DatabaseService) (*uuid.UUID, error) {
	token, err := GetToken(r)
	if err != nil {
		return nil, err
	}
	userId, err := dbService.GetLoggedInUser(*token)
	if err != nil {
		return nil, err
	}
	return userId, nil
}

func AuthMiddleware(next http.Handler, dbService db.DatabaseService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		preflight := EnableCORS(w, r)
		if preflight {
			return
		}
		token, err := GetUser(r, dbService)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(err.Error()))
			return
		}

		r.Header.Set("X-Auth-Token", token.String())

		next.ServeHTTP(w, r)
	})
}

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		preflight := EnableCORS(w, r)
		if preflight {
			return
		}

		next.ServeHTTP(w, r)
	})
}
