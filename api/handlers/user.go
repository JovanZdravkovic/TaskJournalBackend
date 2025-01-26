package handlers

import (
	"net/http"
)

type UserHandler struct{}

func (u *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is users endpoint"))
}
