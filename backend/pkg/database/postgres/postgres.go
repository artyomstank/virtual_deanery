package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/artyomstank/virtual_deanery/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

// New создаёт и настраивает пул соединений с PostgreSQL.
// Принимает контекст и конфигурацию, возвращает готовый пул.
// В случае ошибки возвращает её, обёрнутую в "postgres.New: %w".
func New(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, error) {
	poolConfig, err := pgxpool.ParseConfig(cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("postgres.New: %w", err)
	}

	poolConfig.MaxConns = 25
	poolConfig.MinConns = 5
	poolConfig.MaxConnLifetime = 1 * time.Hour
	poolConfig.MaxConnIdleTime = 30 * time.Minute
	poolConfig.HealthCheckPeriod = 1 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("postgres.New: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("postgres.New: %w", err)
	}

	return pool, nil
}

// Close корректно завершает работу пула соединений.
// Безопасна для вызова с nil-указателем.
func Close(pool *pgxpool.Pool) {
	if pool != nil {
		pool.Close()
	}
}
