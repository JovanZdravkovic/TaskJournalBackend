package handlers

import (
	"net/http"
)

type TaskHistoryHandler struct{}

func (th *TaskHistoryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is task history endpoint"))
}
