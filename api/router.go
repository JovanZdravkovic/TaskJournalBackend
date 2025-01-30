package api

import (
	"net/http"

	"github.com/JovanZdravkovic/TaskJournalBackend/api/handlers"
	"github.com/JovanZdravkovic/TaskJournalBackend/db"
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

func (r *Router) ConfigureRoutes(dbService *db.DatabaseService) {
	r.mux.Handle("/", &handlers.HomeHandler{})
	r.mux.Handle("/task", &handlers.TaskHandler{DBService: dbService})
	r.mux.Handle("/task/", &handlers.TaskHandler{DBService: dbService})
	r.mux.Handle("/tasks", &handlers.TaskHandler{DBService: dbService})
	r.mux.Handle("/tasks/", &handlers.TaskHandler{DBService: dbService})
	r.mux.Handle("/task_history", &handlers.TaskHistoryHandler{DBService: dbService})
	r.mux.Handle("/task_history/", &handlers.TaskHistoryHandler{DBService: dbService})
	r.mux.Handle("/user", &handlers.UserHandler{DBService: dbService})
	r.mux.Handle("/user/", &handlers.UserHandler{DBService: dbService})
	r.mux.Handle("/auth", &handlers.AuthHandler{DBService: dbService})
	r.mux.Handle("/auth/", &handlers.AuthHandler{DBService: dbService})
}

func (r *Router) ListenAndServe() {
	http.ListenAndServe(r.address, r.mux)
}
