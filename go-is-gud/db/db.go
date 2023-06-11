package db

import (
	"context"
	"guigou/bot-is-gud/env"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func New(logger *zap.SugaredLogger, ctx context.Context) *pgxpool.Pool {
	if env.DATABASE_URL == "" {
		return nil
	}

	config, err := pgxpool.ParseConfig(env.DATABASE_URL)
	if err != nil {
		logger.Panic(err)
	}

	config.AfterConnect = func(ctx context.Context, c *pgx.Conn) error {
		logger.Info("new db connection opened")
		return nil
	}

	conn, err := pgxpool.NewWithConfig(ctx, config)

	if err != nil {
		logger.Panic(err)
	}

	err = conn.Ping(ctx)
	if err != nil {
		logger.Panic(err)
	}

	return conn
}
