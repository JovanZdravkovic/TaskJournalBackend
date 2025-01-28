package handlers

import (
	"net/http"

	"github.com/JovanZdravkovic/TaskJournalBackend/db"
)

type TaskHandler struct {
	DBService *db.DatabaseService
}

func (t *TaskHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.DBService.GetTasks()
	w.Write([]byte("This is tasks endpoint"))
}
