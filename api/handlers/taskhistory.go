package handlers

import (
	"net/http"

	"github.com/JovanZdravkovic/TaskJournalBackend/db"
)

type TaskHistoryHandler struct {
	DBService *db.DatabaseService
}

func (th *TaskHistoryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	preflight := EnableCORS(w, r)
	if preflight {
		return
	}
	w.Write([]byte("This is task history endpoint"))
}
