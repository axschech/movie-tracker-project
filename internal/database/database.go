package database

import (
	"context"
	"fmt"

	"github.com/axschech/rockbot-backend/internal/config"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type IPGX interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, arguments ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, arguments ...interface{}) pgx.Row
	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
}

type Database struct {
	P IPGX
}

var _ IPGX = (*pgxpool.Pool)(nil)

func NewDatabase(ctx context.Context, cfg config.DBConfig) (*Database, error) {
	cs := fmt.Sprintf("host=%s user=%s password=%s dbname=%s",
		cfg.Host, cfg.User, cfg.Password, cfg.Database)

	p, err := pgxpool.New(ctx, cs)

	if err != nil {
		return nil, err
	}

	return &Database{P: p}, nil
}
