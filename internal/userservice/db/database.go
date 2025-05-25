package db

import (
	"context"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// interface for pgx
type DBops interface {
	Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error    // for pgxscan
	Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error // for pgxscan
	ExecQueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row           // for one return value
	Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error)
	GetPool(ctx context.Context) *pgxpool.Pool
}

type Database struct {
	pool *pgxpool.Pool
}

func newDatabase(pool *pgxpool.Pool) *Database {
	return &Database{pool: pool}
}

func (db *Database) GetPool(ctx context.Context) *pgxpool.Pool {
	return db.pool
}

func (db *Database) Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return pgxscan.Get(ctx, db.pool, dest, query, args...)
}

func (db *Database) Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return pgxscan.Select(ctx, db.pool, dest, query, args...)
}

func (db *Database) ExecQueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row {
	return db.pool.QueryRow(ctx, query, args...)
}

func (db *Database) Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error) {
	return db.pool.Exec(ctx, query, args...)
}
