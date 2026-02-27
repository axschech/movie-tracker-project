package routing

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Router struct {
	R    *chi.Mux
	Port string
}

func NewRouter(port string) Router {
	r := chi.NewRouter()
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
