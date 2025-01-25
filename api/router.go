package api

import "net/http"

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
	r.mux.Handle("/")
	r.mux.Handle("/task")
	r.mux.Handle("/task/")
	r.mux.Handle("/tasks")
	r.mux.Handle("/tasks/")
	r.mux.Handle("/task_history")
	r.mux.Handle("/task_history/")
	r.mux.Handle("/user")
	r.mux.Handle("/user/")
	r.mux.Handle("/auth")
	r.mux.Handle("/auth/")
}

func (r *Router) ListenAndServe() {
	http.ListenAndServe(r.address, r.mux)
}
