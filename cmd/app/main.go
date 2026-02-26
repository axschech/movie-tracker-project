package main

import (
	"github.com/axschech/rockbot-backend/internal/config"
	"github.com/axschech/rockbot-backend/internal/service"
)

func main() {
	cfg, err := config.MakeConfig()
	if err != nil {
		panic(err)
	}

	service := service.NewService(cfg)

	err = service.Run()
	if err != nil {
		panic(err)
	}
}
