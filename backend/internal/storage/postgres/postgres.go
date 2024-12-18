package postgres

import (
	"TextVault/internal/storage/cloud"
	"context"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	conn    *pgxpool.Pool
	s3cloud *cloud.CloudStorage

	timeout time.Duration
}

// New returns a new Storage instance based on the DATABASE_URL environment
// variable. The function will panic if the environment variable is not set or
// if the connection to the database cannot be established.
func New(cloudStorage *cloud.CloudStorage, ctx context.Context) (*Storage, error) {
	const timeout = 5 * time.Second

	poolConfig, err := pgxpool.ParseConfig(os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}

	conn, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, err
	}

	return &Storage{
		conn:    conn,
		s3cloud: cloudStorage,
		timeout: timeout,
	}, nil
}

func (s *Storage) Close() {
	s.conn.Close()
}

func (s *Storage) Ping(ctx context.Context) error {
	return s.conn.Ping(ctx)
}
