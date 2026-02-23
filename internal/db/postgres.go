package db

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Nexi77/fleetcommander-backend/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresDB struct {
	Pool *pgxpool.Pool
}

func NewPostgresConnection(ctx context.Context, cfg *config.Config) (*PostgresDB, error) {
	slog.Debug("Initializing PostgreSQL connection pool",
		"max_conns", cfg.DBMaxConns,
		"min_conns", cfg.DBMinConns,
	)

	poolConfig, err := pgxpool.ParseConfig(cfg.DBUrl)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse database URL: %w", err)
	}

	poolConfig.MaxConns = cfg.DBMaxConns
	poolConfig.MinConns = cfg.DBMinConns
	poolConfig.MaxConnLifetime = cfg.DBMaxConnLifetime
	poolConfig.MaxConnIdleTime = cfg.DBMaxConnIdleTime

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("Unable to connect to database: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("Unable to ping database: %w", err)
	}

	slog.Info("PostgreSQL connection established successfully")

	return &PostgresDB{
		Pool: pool,
	}, nil
}

func (db *PostgresDB) Close() {
	if db.Pool != nil {
		slog.Info("Closing PostgreSQL connection pool...")
		db.Pool.Close()
	}
}
