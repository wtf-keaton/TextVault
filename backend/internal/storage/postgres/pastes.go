package postgres

import (
	"TextVault/internal/storage"
	"TextVault/internal/storage/models"
	"context"
	"fmt"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
)

func (s *Storage) GetPaste(ctx context.Context, hash string) (models.Paste, error) {
	const prefix = "storage.postgresql.GetUser"

	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	stmt := "SELECT ID, title, hash, authorid FROM Pastes WHERE hash = $1"

	var paste models.Paste
	err := pgxscan.Get(ctx, s.conn, &paste, stmt, hash)
	if err != nil {
		if err == pgx.ErrNoRows {
			return models.Paste{}, fmt.Errorf("%s: %w", prefix, storage.ErrPasteNotFound)
		}

		return models.Paste{}, fmt.Errorf("%s: %w", prefix, err)
	}

	return paste, nil
}

func (s *Storage) SavePaste(ctx context.Context, paste *models.Paste, content []byte) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	stmt := "INSERT INTO Pastes (title, hash, authorid) VALUES ($1, $2, $3) RETURNING id"

	var id int64
	err := s.conn.QueryRow(ctx, stmt, paste.Title, paste.Hash, paste.AuthorID).Scan(&id)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) GetPastesByUser(ctx context.Context, userID int64) ([]models.Paste, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	stmt := "SELECT ID, title, hash, authorid FROM Pastes WHERE authorid = $1"

	var pastes []models.Paste
	err := pgxscan.Select(ctx, s.conn, &pastes, stmt, userID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return []models.Paste{}, storage.ErrUserDontHavePastes
		}

		return nil, err
	}

	return pastes, nil
}

func (s *Storage) DeletePaste(ctx context.Context, hash string) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	stmt := "DELETE FROM Pastes WHERE hash = $1"

	_, err := s.conn.Exec(ctx, stmt, hash)
	if err != nil {
		return err
	}

	return nil
}
