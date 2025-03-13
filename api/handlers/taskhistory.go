package handlers

import (
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/JovanZdravkovic/TaskJournalBackend/db"
	"github.com/google/uuid"
)

var (
	TaskHistoryID = regexp.MustCompile(`^/task_history/([a-fA-F0-9\-]{36})$`)
	TasksHistory  = regexp.MustCompile(`^/tasks_history/*$`)
)

type TaskHistoryHandler struct {
	DBService   *db.DatabaseService
	AuthService *AuthHandler
}

func (th *TaskHistoryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	preflight := EnableCORS(w, r)
	if preflight {
		return
	}

	userId, err := th.AuthService.GetUser(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(err.Error()))
		return
	}

	switch {
	case r.Method == http.MethodGet && TaskHistoryID.MatchString(r.URL.Path):
		th.GetTaskHistory(w, r, *userId)
		return
	case r.Method == http.MethodGet && TasksHistory.MatchString(r.URL.Path):
		th.GetTasksHistory(w, r, *userId)
		return
	default:
		return
	}
}

func (th *TaskHistoryHandler) GetTaskHistory(w http.ResponseWriter, r *http.Request, userId uuid.UUID) {
}

func (th *TaskHistoryHandler) GetTasksHistory(w http.ResponseWriter, r *http.Request, userId uuid.UUID) {
	searchName := r.URL.Query().Get("searchName")
	searchIcons := r.URL.Query()["searchIcons"]
	searchRating := r.URL.Query().Get("searchRating")
	tasksHistory, err := th.DBService.GetTasksHistory(userId, &searchName, searchIcons, &searchRating)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	tasksHistoryJson, err := json.Marshal(tasksHistory)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(tasksHistoryJson)
}
