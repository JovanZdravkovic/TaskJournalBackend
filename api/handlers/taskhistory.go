package handlers

import (
	"net/http"

	"github.com/JovanZdravkovic/TaskJournalBackend/db"
)

type TaskHistoryHandler struct {
	DBService *db.DatabaseService
}

func (th *TaskHistoryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is task history endpoint"))
}
