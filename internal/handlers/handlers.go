package handlers

import (
	"net/http"

	"github.com/axschech/rockbot-backend/internal/user"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func Ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Pong"))
}

func UserHandler(w http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	user := user.GetUserByID(id)

	render.JSON(w, req, &user)
}
