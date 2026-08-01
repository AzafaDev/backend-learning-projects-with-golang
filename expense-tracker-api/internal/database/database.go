package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewDatabase(connStr string, ctx context.Context) (*pgxpool.Pool, error) {
	cfgPool, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, err
	}

	dbPool, err := pgxpool.NewWithConfig(ctx, cfgPool)
	if err != nil {
		return nil, err
	}

	return dbPool, nil
}
