package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/willfreit4s/chainPulse/configs"
	"github.com/willfreit4s/chainPulse/pkg/logger"
)

func InitDatabase(ctx context.Context, cfg *configs.Config) (*pgxpool.Pool, error) {
	log := logger.FromContext(ctx)

	log.Info().Msg("initializing postgres database")
	log.Info().Msgf("Postgres config: maxConn=%d minConn=%d", cfg.MaxConn, cfg.MinConn)

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		cfg.DBUser,
		cfg.DBPass,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	config.MaxConns = int32(cfg.MaxConn)
	config.MinConns = int32(cfg.MinConn)

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	// ping
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	log.Info().Msg("postgres database initialized successfully")

	return pool, nil
}
