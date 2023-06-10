package db

import (
	"context"
	"guigou/bot-is-gud/env"

	"github.com/jackc/pgx/v4"
)

func New(ctx context.Context) *pgx.Conn {
	if env.DATABASE_URL == "" {
		return nil
	}

	conn, err := pgx.Connect(ctx, env.DATABASE_URL)

	if err != nil {
		panic(err)
	}

	err = conn.Ping(ctx)
	if err != nil {
		panic(err)
	}

	return conn
}
