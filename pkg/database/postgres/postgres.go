// pkg/database/postgres/postgres.go
package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// NewPool creates new postgres connection pool.
func NewPool(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
	// TODO: Parse database URL

	// TODO: Create pgxpool config

	// TODO: Set max connections

	// TODO: Create pool with retry logic

	// TODO: Test connection with Ping()

	return nil, fmt.Errorf("not implemented")
}

// Tx represents database transaction interface for mocking.
type Tx interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

// BeginTx starts new transaction.
func BeginTx(ctx context.Context, pool *pgxpool.Pool) (Tx, error) {
	// TODO: Start transaction

	return nil, fmt.Errorf("not implemented")
}
