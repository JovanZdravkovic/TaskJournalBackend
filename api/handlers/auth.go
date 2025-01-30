package handlers

import (
	"net/http"

	"github.com/JovanZdravkovic/TaskJournalBackend/db"
)

type AuthHandler struct {
	DBService *db.DatabaseService
}

func (a *AuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is auth endpoint"))
}
