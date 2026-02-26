package database

import (
	"context"
	"fmt"

	"github.com/axschech/rockbot-backend/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	P *pgxpool.Pool
}

func NewDatabase(ctx context.Context, cfg config.DBConfig) (*Database, error) {
	cs := fmt.Sprintf("host=%s user=%s password=%s dbname=%s",
		cfg.Host, cfg.User, cfg.Password, cfg.Database)

	p, err := pgxpool.New(ctx, cs)

	if err != nil {
		return nil, err
	}

	return &Database{P: p}, nil
}
