package api

import (
	"net/http"

	"github.com/JovanZdravkovic/TaskJournalBackend/api/handlers"
)

type Router struct {
	address string
	mux     *http.ServeMux
}

func NewRouter(address string) *Router {
	return &Router{
		address: address,
		mux:     http.NewServeMux(),
	}
}

func (r *Router) ConfigureRoutes() {
	r.mux.Handle("/", &handlers.HomeHandler{})
	r.mux.Handle("/task", &handlers.TaskHandler{})
	r.mux.Handle("/task/", &handlers.TaskHandler{})
	r.mux.Handle("/tasks", &handlers.TaskHandler{})
	r.mux.Handle("/tasks/", &handlers.TaskHandler{})
	r.mux.Handle("/task_history", &handlers.TaskHistoryHandler{})
	r.mux.Handle("/task_history/", &handlers.TaskHistoryHandler{})
	r.mux.Handle("/user", &handlers.UserHandler{})
	r.mux.Handle("/user/", &handlers.UserHandler{})
	r.mux.Handle("/auth", &handlers.AuthHandler{})
	r.mux.Handle("/auth/", &handlers.AuthHandler{})
}

func (r *Router) ListenAndServe() {
	http.ListenAndServe(r.address, r.mux)
}
