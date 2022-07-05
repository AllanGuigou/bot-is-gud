package db

import (
	"context"
	"guigou/bot-is-gud/env"
	"log"

	"github.com/jackc/pgx/v4"
)

func New() *pgx.Conn {
	ctx := context.Background()

	conn, err := pgx.Connect(ctx, env.DATABASE_URL)

	if err != nil {
		log.Fatal(err)
	}

	return conn
}