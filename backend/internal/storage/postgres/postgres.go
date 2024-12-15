package postgres

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	conn *pgxpool.Pool
}

// New returns a new Storage instance based on the DATABASE_URL environment
// variable. The function will panic if the environment variable is not set or
// if the connection to the database cannot be established.
func New(ctx context.Context) *Storage {
	poolConfig, err := pgxpool.ParseConfig(os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}

	conn, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		panic(err)
	}

	return &Storage{conn: conn}
}
