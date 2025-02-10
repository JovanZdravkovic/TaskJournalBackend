package handlers

import (
	"net/http"

	"github.com/JovanZdravkovic/TaskJournalBackend/db"
)

type UserHandler struct {
	DBService *db.DatabaseService
}

func (u *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	preflight := EnableCORS(w, r)
	if preflight {
		return
	}
	w.Write([]byte("This is users endpoint"))
}
