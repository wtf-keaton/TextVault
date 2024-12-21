package postgres

import (
	"TextVault/internal/config"
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	conn    *pgxpool.Pool
	timeout time.Duration
}

// New returns a new Storage instance based on the DATABASE_URL environment
// variable. The function will panic if the environment variable is not set or
// if the connection to the database cannot be established.
func New(ctx context.Context, log *slog.Logger, cfg *config.PostgresConfig) (*Storage, error) {
	const timeout = 5 * time.Second

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
	log.Debug("Connecting to database", "dsn", dsn)
	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	conn, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, err
	}

	return &Storage{
		conn:    conn,
		timeout: timeout,
	}, nil
}

func (s *Storage) Close() {
	s.conn.Close()
}

func (s *Storage) Ping(ctx context.Context) error {
	return s.conn.Ping(ctx)
}
