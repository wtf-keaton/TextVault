package postgres

import (
	"TextVault/internal/storage"
	"TextVault/internal/storage/models"
	"context"
	"fmt"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
)

func (s *Storage) SaveUser(ctx context.Context, username, email, password string) (int64, error) {
	const prefix = "storage.postgresql.SaveUser"

	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	stmt := "INSERT INTO User (name, email, password) VALUES ($1, $2, $3) RETURNING id"

	var id int64
	err := s.conn.QueryRow(ctx, stmt, username, email, password).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", prefix, err)
	}

	return id, nil
}

func (s *Storage) GetUser(ctx context.Context, username string) (models.User, error) {
	const prefix = "storage.postgresql.GetUser"

	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	stmt := "SELECT id, name, email, password FROM User WHERE name = $1 OR email = $1"

	var user models.User
	err := pgxscan.Select(ctx, s.conn, &user, stmt, username)
	if err != nil {
		if err == pgx.ErrNoRows {
			return models.User{}, fmt.Errorf("%s: %w", prefix, storage.ErrUserNotFound)
		}

		return models.User{}, fmt.Errorf("%s: %w", prefix, err)
	}

	return user, nil
}

func (s *Storage) UpdateUser(ctx context.Context, user *models.User) error {
	return nil
}
