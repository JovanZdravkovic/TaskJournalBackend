package handlers

import (
	"net/http"
)

type AuthHandler struct{}

func (a *AuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is auth endpoint"))
}
