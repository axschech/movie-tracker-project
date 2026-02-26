package service

import (
	"fmt"

	"github.com/axschech/rockbot-backend/internal/config"
	"github.com/axschech/rockbot-backend/internal/routing"
)

type Service struct {
	Config config.Config
}

func NewService(cfg config.Config) *Service {
	return &Service{
		Config: cfg,
	}
}

func (s *Service) Run() error {
	fmt.Printf("Starting server on port %s\n", s.Config.Port)
	r := routing.NewRouter(s.Config.Port)
	r.MakeRoutes()

	return r.Listen()
}
