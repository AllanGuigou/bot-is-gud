package db

import (
	"context"
	"guigou/bot-is-gud/env"
	"log"

	"github.com/jackc/pgx/v4"
)

func New(ctx context.Context) *pgx.Conn {
	if env.DATABASE_URL == "" {
		return nil
	}

	conn, err := pgx.Connect(ctx, env.DATABASE_URL)

	if err != nil {
		log.Fatal(err)
	}

	return conn
}
