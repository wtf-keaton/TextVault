package postgres

import (
	"TextVault/internal/storage"
	"TextVault/internal/storage/models"
	"context"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
)

func (s *Storage) GetPaste(ctx context.Context, hash string) (models.Paste, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	stmt := "SELECT ID, title, hash, language, authorid FROM Pastes WHERE hash = $1"

	var paste models.Paste
	err := pgxscan.Get(ctx, s.conn, &paste, stmt, hash)
	if err != nil {
		if err == pgx.ErrNoRows {
			return models.Paste{}, storage.ErrPasteNotFound
		}

		return models.Paste{}, err
	}

	return paste, nil
}

func (s *Storage) SavePaste(ctx context.Context, paste *models.Paste, content []byte) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	stmt := "INSERT INTO Pastes (title, hash, language, authorid) VALUES ($1, $2, $3, $4) RETURNING id"

	var id int64
	err := s.conn.QueryRow(ctx, stmt, paste.Title, paste.Hash, paste.Language, paste.AuthorID).Scan(&id)
	if err != nil {
		return err
	}

	return nil
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
