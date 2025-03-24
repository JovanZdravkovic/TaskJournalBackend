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
	TaskUpdateID = regexp.MustCompile(`^/task/update/([a-fA-F0-9\-]{36})$`)
	Tasks        = regexp.MustCompile(`^/tasks/*$`)
)

type TaskHandler struct {
	DBService *db.DatabaseService
}

func (t *TaskHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
	case r.Method == http.MethodGet && TaskID.MatchString(r.URL.Path):
		t.GetTask(w, r, token)
		return
	case r.Method == http.MethodPut && TaskID.MatchString(r.URL.Path):
		t.CompleteTask(w, r, token)
		return
	case r.Method == http.MethodPut && TaskUpdateID.MatchString(r.URL.Path):
		t.UpdateTask(w, r, token)
		return
	case r.Method == http.MethodDelete && TaskID.MatchString(r.URL.Path):
		t.DeleteTask(w, r, token)
		return
	case r.Method == http.MethodPost && Tasks.MatchString(r.URL.Path):
		t.CreateTask(w, r, token)
		return
	case r.Method == http.MethodGet && Tasks.MatchString(r.URL.Path):
		t.GetTasks(w, r, token)
		return
	default:
		return
	}
}

func (t *TaskHandler) GetTasks(w http.ResponseWriter, r *http.Request, userId uuid.UUID) {
	searchName := r.URL.Query().Get("searchName")
	searchIcons := r.URL.Query()["searchIcons"]
	searchOrderBy := r.URL.Query().Get("searchOrderBy")
	tasks, err := t.DBService.GetTasks(userId, &searchName, searchIcons, &searchOrderBy)
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

func (t *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request, userId uuid.UUID) {
	taskId, err := uuid.Parse(path.Base(r.URL.Path))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error proccessing the uuid"))
		return
	}
	var task db.TaskPut
	err = json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("bad request"))
		return
	}
	err = t.DBService.UpdateTask(taskId, task, userId)
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

func (t *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request, userId uuid.UUID) {
	taskId, err := uuid.Parse(path.Base(r.URL.Path))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error proccessing the uuid"))
		return
	}
	err = t.DBService.DeleteTask(taskId, userId)
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
