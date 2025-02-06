package handlers

import (
	"encoding/json"
	"net/http"
	"path"
	"regexp"

	"github.com/JovanZdravkovic/TaskJournalBackend/db"
	"github.com/google/uuid"
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
	tokenCookie, err := r.Cookie("sessiontoken")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Failed to read authentication token"))
		return
	}
	tokenString := tokenCookie.Value
	token, err := uuid.Parse(tokenString)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Failed to read authentication token"))
		return
	}
	userId, err := t.DBService.GetLoggedInUser(token)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Invalid token"))
		return
	}

	switch {
	case r.Method == http.MethodGet && TaskID.MatchString(r.URL.Path):
		t.GetTask(w, r, *userId)
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
		t.GetTasks(w, r, *userId)
		return
	case r.Method == http.MethodGet && TasksStarred.MatchString(r.URL.Path):
		t.GetStarredTasks(w, r, *userId)
		return
	default:
		return
	}
}

func (t *TaskHandler) GetTasks(w http.ResponseWriter, r *http.Request, userId uuid.UUID) {
	tasks, err := t.DBService.GetTasks(userId)
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

func (t *TaskHandler) GetStarredTasks(w http.ResponseWriter, r *http.Request, userId uuid.UUID) {
	tasks, err := t.DBService.GetStarredTasks(userId)
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

func (t *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request, userId uuid.UUID) {
	taskId, err := uuid.Parse(path.Base(r.URL.Path))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error proccessing the uuid"))
		return
	}
	task, err := t.DBService.GetTask(taskId, userId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("The task with given id doesn't exist"))
		return
	}
	taskJson, err := json.Marshal(task)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("An internal server error occured"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(taskJson)
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
