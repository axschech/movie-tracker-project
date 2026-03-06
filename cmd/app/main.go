package main

import (
	"context"
	"net/http"

	"github.com/axschech/rockbot-backend/internal/config"
	"github.com/axschech/rockbot-backend/internal/database"
	"github.com/axschech/rockbot-backend/internal/database/repository"
	"github.com/axschech/rockbot-backend/internal/routing"
	"github.com/axschech/rockbot-backend/internal/service"
)

func main() {
	cfg, err := config.MakeConfig()
	if err != nil {
		panic(err)
	}
	ctx := context.Background()

	db, err := database.NewDatabase(ctx, cfg.DB)
	if err != nil {
		panic(err)
	}

	r := repository.NewRepository(ctx, db)

	router := routing.NewRouter(cfg.Port)

	httpClient := &http.Client{}

	service := service.NewService(
		cfg,
		*r,
		router,
		httpClient,
	)

	err = service.Run()
	if err != nil {
		panic(err)
	}
}
