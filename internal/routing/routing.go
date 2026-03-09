package routing

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

type Router struct {
	R    *chi.Mux
	Port string
}

func NewRouter(port string) Router {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders: []string{"Link"},
		MaxAge:         300, // Maximum value not ignored by any of major browsers
	}))

	return Router{
		R:    r,
		Port: port,
	}
}

func (rt *Router) Listen() error {
	if rt.Port == "" {
		return fmt.Errorf("port not set")
	}
	return http.ListenAndServe(fmt.Sprintf(":%s", rt.Port), rt.R)
}
