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
	homeHandler := handlers.HomeHandler{}
	authHandler := handlers.AuthHandler{DBService: dbService}
	taskHandler := handlers.TaskHandler{DBService: dbService, AuthService: &authHandler}
	taskHistoryHandler := handlers.TaskHistoryHandler{DBService: dbService, AuthService: &authHandler}
	userHandler := handlers.UserHandler{DBService: dbService, AuthService: &authHandler}
	r.mux.Handle("/", &homeHandler)
	r.mux.Handle("/task", &taskHandler)
	r.mux.Handle("/task/", &taskHandler)
	r.mux.Handle("/tasks", &taskHandler)
	r.mux.Handle("/tasks/", &taskHandler)
	r.mux.Handle("/task_history", &taskHistoryHandler)
	r.mux.Handle("/task_history/", &taskHistoryHandler)
	r.mux.Handle("/tasks_history", &taskHistoryHandler)
	r.mux.Handle("/tasks_history/", &taskHistoryHandler)
	r.mux.Handle("/user", &userHandler)
	r.mux.Handle("/user/", &userHandler)
	r.mux.Handle("/auth", &authHandler)
	r.mux.Handle("/auth/", &authHandler)
}

func (r *Router) ListenAndServe() {
	http.ListenAndServe(r.address, r.mux)
}
