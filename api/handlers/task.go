package handlers

import (
	"net/http"
)

type TaskHandler struct{}

func (t *TaskHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is tasks endpoint"))
}
