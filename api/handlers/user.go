package handlers

import (
	"net/http"

	"github.com/JovanZdravkovic/TaskJournalBackend/db"
)

type UserHandler struct {
	DBService *db.DatabaseService
}

func (u *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is users endpoint"))
}
