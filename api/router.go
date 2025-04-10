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
	r.mux.Handle("/task", handlers.CORSMiddleware(handlers.AuthMiddleware(&taskHandler, *dbService)))
	r.mux.Handle("/task/", handlers.CORSMiddleware(handlers.AuthMiddleware(&taskHandler, *dbService)))
	r.mux.Handle("/tasks", handlers.CORSMiddleware(handlers.AuthMiddleware(&taskHandler, *dbService)))
	r.mux.Handle("/tasks/", handlers.CORSMiddleware(handlers.AuthMiddleware(&taskHandler, *dbService)))
	r.mux.Handle("/task_history", handlers.CORSMiddleware(handlers.AuthMiddleware(&taskHistoryHandler, *dbService)))
	r.mux.Handle("/task_history/", handlers.CORSMiddleware(handlers.AuthMiddleware(&taskHistoryHandler, *dbService)))
	r.mux.Handle("/tasks_history", handlers.CORSMiddleware(handlers.AuthMiddleware(&taskHistoryHandler, *dbService)))
	r.mux.Handle("/tasks_history/", handlers.CORSMiddleware(handlers.AuthMiddleware(&taskHistoryHandler, *dbService)))
	r.mux.Handle("/user", handlers.CORSMiddleware(handlers.AuthMiddleware(&userHandler, *dbService)))
	r.mux.Handle("/user/", handlers.CORSMiddleware(handlers.AuthMiddleware(&userHandler, *dbService)))
	r.mux.Handle("/auth", handlers.CORSMiddleware(handlers.AuthMiddleware(&authHandler, *dbService)))
	r.mux.Handle("/auth/", handlers.CORSMiddleware(handlers.AuthMiddleware(&authHandler, *dbService)))
	r.mux.Handle("/login", handlers.CORSMiddleware(&loginHandler))
	r.mux.Handle("/login/", handlers.CORSMiddleware(&loginHandler))
	r.mux.Handle("/logout", handlers.CORSMiddleware(&logoutHandler))
	r.mux.Handle("/logout/", handlers.CORSMiddleware(&logoutHandler))
	r.mux.Handle("/signup", handlers.CORSMiddleware(&signupHandler))
	r.mux.Handle("/signup/", handlers.CORSMiddleware(&signupHandler))
}

func (r *Router) ListenAndServe() {
	http.ListenAndServe(r.address, r.mux)
}
