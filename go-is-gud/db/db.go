package db

import (
	"context"
	"fmt"
	"guigou/bot-is-gud/env"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func New(ctx context.Context) *pgxpool.Pool {
	if env.DATABASE_URL == "" {
		return nil
	}

	config, err := pgxpool.ParseConfig(env.DATABASE_URL)
	if err != nil {
		panic(err)
	}

	config.AfterConnect = func(ctx context.Context, c *pgx.Conn) error {
		fmt.Println("new db connection opened")
		return nil
	}

	conn, err := pgxpool.NewWithConfig(ctx, config)

	if err != nil {
		panic(err)
	}

	err = conn.Ping(ctx)
	if err != nil {
		panic(err)
	}

	return conn
}
