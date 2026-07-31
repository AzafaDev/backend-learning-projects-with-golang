package database

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectDB(ctx context.Context, dbUrl string) (*pgxpool.Pool, error) {
	pgCfg, err := pgxpool.ParseConfig(dbUrl)
	if err != nil {
		return nil, err
	}
	pgCfg.MaxConns = 10
	pgCfg.MinConns = 2
	pgCfg.MaxConnLifetime = time.Hour
	pgCfg.MaxConnIdleTime = 30 * time.Minute
	pgCfg.HealthCheckPeriod = time.Minute
	pool, err := pgxpool.NewWithConfig(ctx, pgCfg)
	if err != nil {
		return nil, err
	}
	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}
return pool, nil
}
