package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/JovanZdravkovic/TaskJournalBackend/db"
)

type TaskHandler struct {
	DBService *db.DatabaseService
}

func (t *TaskHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tasks, err := t.DBService.GetTasks()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	tasksJson, err := json.Marshal(tasks)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(tasksJson)
}
