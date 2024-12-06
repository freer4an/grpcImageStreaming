package connections

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Postgres(ctx context.Context, dbUrl string) *pgxpool.Pool {
	cfg, err := pgxpool.ParseConfig(dbUrl)
	if err != nil {
		panic(fmt.Errorf("failed to parse postgres connection string: %w", err))
	}

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		panic(fmt.Errorf("cannot connect to database: %w", err))
	}
	return pool
}
