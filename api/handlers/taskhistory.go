package handlers

import (
	"encoding/json"
	"net/http"
	"path"
	"regexp"
	"strconv"

	"github.com/JovanZdravkovic/TaskJournalBackend/db"
	"github.com/google/uuid"
)

var (
	TaskHistoryID = regexp.MustCompile(`^/task_history/([a-fA-F0-9\-]{36})$`)
	TasksHistory  = regexp.MustCompile(`^/tasks_history/*$`)
)

type TaskHistoryHandler struct {
	DBService *db.DatabaseService
}

func (th *TaskHistoryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	preflight := EnableCORS(w, r)
	if preflight {
		return
	}

	tokenString := r.Header.Get("X-Auth-Token")
	if tokenString == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	token, err := uuid.Parse(tokenString)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	switch {
	case r.Method == http.MethodGet && TaskHistoryID.MatchString(r.URL.Path):
		th.GetTaskHistory(w, r, token)
		return
	case r.Method == http.MethodGet && TasksHistory.MatchString(r.URL.Path):
		th.GetTasksHistory(w, r, token)
		return
	case r.Method == http.MethodPut && TaskHistoryID.MatchString(r.URL.Path):
		th.UpdateTaskHistory(w, r, token)
		return
	case r.Method == http.MethodDelete && TaskHistoryID.MatchString(r.URL.Path):
		th.DeleteTaskAndHistory(w, r, token)
		return
	default:
		return
	}
}

func (th *TaskHistoryHandler) GetTaskHistory(w http.ResponseWriter, r *http.Request, userId uuid.UUID) {
	taskHistoryId, err := uuid.Parse(path.Base(r.URL.Path))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error proccessing the uuid"))
		return
	}
	taskHistory, err := th.DBService.GetTaskHistory(taskHistoryId, userId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("task history with given id doesn't exist"))
		return
	}
	taskHistoryJson, err := json.Marshal(taskHistory)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal server error occured"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(taskHistoryJson)
}

func (th *TaskHistoryHandler) GetTasksHistory(w http.ResponseWriter, r *http.Request, userId uuid.UUID) {
	searchName := r.URL.Query().Get("searchName")
	searchIcons := r.URL.Query()["searchIcons"]
	searchRatingString := r.URL.Query().Get("searchRating")
	var searchRating int
	if searchRatingString == "1" || searchRatingString == "2" || searchRatingString == "3" {
		searchRating, _ = strconv.Atoi(searchRatingString)
	}
	tasksHistory, err := th.DBService.GetTasksHistory(userId, &searchName, searchIcons, searchRating)
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

func (th *TaskHistoryHandler) UpdateTaskHistory(w http.ResponseWriter, r *http.Request, userId uuid.UUID) {
	taskHistoryId, err := uuid.Parse(path.Base(r.URL.Path))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error proccessing the uuid"))
		return
	}
	var taskHistory db.TaskHistoryPut
	err = json.NewDecoder(r.Body).Decode(&taskHistory)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("bad request"))
		return
	}
	err = th.DBService.UpdateTaskHistory(taskHistoryId, taskHistory, userId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	responseJson, err := json.Marshal(db.Success{Success: true})
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJson)
}

func (th *TaskHistoryHandler) DeleteTaskAndHistory(w http.ResponseWriter, r *http.Request, userId uuid.UUID) {
	taskHistoryId, err := uuid.Parse(path.Base(r.URL.Path))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error proccessing the uuid"))
		return
	}
	err = th.DBService.DeleteTaskAndHistory(taskHistoryId, userId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	responseJson, err := json.Marshal(db.Success{Success: true})
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJson)
}
