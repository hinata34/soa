package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewDB(ctx context.Context, dbUrl string) DBops {
	config, err := pgxpool.ParseConfig(dbUrl)
	if err != nil {
		log.Fatalf("failed to parsed pgxpool config: %v", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Fatalf("failed to create pgx pool: %v", err)
	}
	return newDatabase(pool)
}
