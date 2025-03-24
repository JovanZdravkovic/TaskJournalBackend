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
	taskHandler := handlers.TaskHandler{DBService: dbService}
	taskHistoryHandler := handlers.TaskHistoryHandler{DBService: dbService}
	userHandler := handlers.UserHandler{DBService: dbService}
	loginHandler := handlers.LoginHandler{DBService: dbService}
	logoutHandler := handlers.LogoutHandler{DBService: dbService}
	signupHandler := handlers.SignupHandler{DBService: dbService}
	r.mux.Handle("/", &homeHandler)
	r.mux.Handle("/task", handlers.AuthMiddleware(&taskHandler, *dbService))
	r.mux.Handle("/task/", handlers.AuthMiddleware(&taskHandler, *dbService))
	r.mux.Handle("/tasks", handlers.AuthMiddleware(&taskHandler, *dbService))
	r.mux.Handle("/tasks/", handlers.AuthMiddleware(&taskHandler, *dbService))
	r.mux.Handle("/task_history", handlers.AuthMiddleware(&taskHistoryHandler, *dbService))
	r.mux.Handle("/task_history/", handlers.AuthMiddleware(&taskHistoryHandler, *dbService))
	r.mux.Handle("/tasks_history", handlers.AuthMiddleware(&taskHistoryHandler, *dbService))
	r.mux.Handle("/tasks_history/", handlers.AuthMiddleware(&taskHistoryHandler, *dbService))
	r.mux.Handle("/user", handlers.AuthMiddleware(&userHandler, *dbService))
	r.mux.Handle("/user/", handlers.AuthMiddleware(&userHandler, *dbService))
	r.mux.Handle("/auth", handlers.AuthMiddleware(&authHandler, *dbService))
	r.mux.Handle("/auth/", handlers.AuthMiddleware(&authHandler, *dbService))
	r.mux.Handle("/login", &loginHandler)
	r.mux.Handle("/login/", &loginHandler)
	r.mux.Handle("/logout", &logoutHandler)
	r.mux.Handle("/logout/", &logoutHandler)
	r.mux.Handle("/signup", &signupHandler)
	r.mux.Handle("/signup/", &signupHandler)
}

func (r *Router) ListenAndServe() {
	http.ListenAndServe(r.address, r.mux)
}
