package routing

import (
	"fmt"
	"net/http"

	"github.com/axschech/rockbot-backend/internal/handlers"
	"github.com/go-chi/chi/v5"
)

type Router struct {
	r    *chi.Mux
	Port string
}

func NewRouter(port string) Router {
	r := chi.NewRouter()
	return Router{
		r:    r,
		Port: port,
	}
}

func (rt *Router) Listen() error {
	if rt.Port == "" {
		return fmt.Errorf("port not set")
	}
	return http.ListenAndServe(fmt.Sprintf(":%s", rt.Port), rt.r)
}

func (rt *Router) MakeRoutes() {
	rt.r.Get("/ping", handlers.Ping)
	rt.r.Route("/api", func(r chi.Router) {
		r.Get("/user/{id}", handlers.UserHandler)
	})
}
