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
	TaskID = regexp.MustCompile(`^/task/([a-fA-F0-9\-]{36})$`)
	Tasks  = regexp.MustCompile(`^/tasks/*$`)
)

type TaskHandler struct {
	DBService   *db.DatabaseService
	AuthService *AuthHandler
}

func (t *TaskHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	preflight := EnableCORS(w, r)
	if preflight {
		return
	}

	userId, err := t.AuthService.GetUser(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(err.Error()))
		return
	}

	switch {
	case r.Method == http.MethodGet && TaskID.MatchString(r.URL.Path):
		t.GetTask(w, r, *userId)
		return
	case r.Method == http.MethodPut && TaskID.MatchString(r.URL.Path):
		t.CompleteTask(w, r, *userId)
		return
	case r.Method == http.MethodDelete && TaskID.MatchString(r.URL.Path):
		t.DeleteTask(w, r)
		return
	case r.Method == http.MethodPost && Tasks.MatchString(r.URL.Path):
		t.CreateTask(w, r, *userId)
		return
	case r.Method == http.MethodGet && Tasks.MatchString(r.URL.Path):
		t.GetTasks(w, r, *userId)
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

func (t *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request, userId uuid.UUID) {
	taskId, err := uuid.Parse(path.Base(r.URL.Path))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error proccessing the uuid"))
		return
	}
	task, err := t.DBService.GetTask(taskId, userId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("task with given id doesn't exist"))
		return
	}
	taskJson, err := json.Marshal(task)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal server error occured"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(taskJson)
}

func (t *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request, userId uuid.UUID) {
	var task db.TaskPost
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("bad request"))
		return
	}
	task.CreatedBy = userId
	taskId, err := t.DBService.CreateTask(task)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	taskIdJson, err := json.Marshal(db.Id{Id: *taskId})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error while constructing json"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(taskIdJson)
}

func (t *TaskHandler) CompleteTask(w http.ResponseWriter, r *http.Request, userId uuid.UUID) {
	taskId, err := uuid.Parse(path.Base(r.URL.Path))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error proccessing the uuid"))
		return
	}
	_, err = t.DBService.CompleteTask(taskId, userId)
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

func (t *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("This is delete single task endpoint"))
}
