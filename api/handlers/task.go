package handlers

import (
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/JovanZdravkovic/TaskJournalBackend/db"
)

var (
	TaskID       = regexp.MustCompile(`^/task/([a-fA-F0-9\-]{36})$`)
	Tasks        = regexp.MustCompile(`^/tasks/*$`)
	TasksStarred = regexp.MustCompile(`^/tasks/starred/*$`)
)

type TaskHandler struct {
	DBService *db.DatabaseService
}

func (t *TaskHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodGet && TaskID.MatchString(r.URL.Path):
		t.GetTask(w, r)
		return
	case r.Method == http.MethodPost && TaskID.MatchString(r.URL.Path):
		t.CreateTask(w, r)
		return
	case r.Method == http.MethodPut && TaskID.MatchString(r.URL.Path):
		t.UpdateTask(w, r)
		return
	case r.Method == http.MethodDelete && TaskID.MatchString(r.URL.Path):
		t.DeleteTask(w, r)
		return
	case r.Method == http.MethodGet && Tasks.MatchString(r.URL.Path):
		t.GetTasks(w, r)
		return
	case r.Method == http.MethodGet && TasksStarred.MatchString(r.URL.Path):
		t.GetStarredTasks(w, r)
		return
	default:
		return
	}
}

func (t *TaskHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
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
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(tasksJson)
}

func (t *TaskHandler) GetStarredTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := t.DBService.GetStarredTasks()
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
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(tasksJson)
}

func (t *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("This is get single task endpoint"))
}

func (t *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("This is create single task endpoint"))
}

func (t *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("This is update single task endpoint"))
}

func (t *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("This is delete single task endpoint"))
}
