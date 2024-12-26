package postgres

import (
	"TextVault/internal/storage"
	"TextVault/internal/storage/models"
	"context"
	"errors"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
)

func (s *Storage) GetPaste(ctx context.Context, id string) (models.Paste, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	stmt := "SELECT * FROM Pastes WHERE id = $1"

	var paste models.Paste
	err := pgxscan.Get(ctx, s.conn, &paste, stmt, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Paste{}, storage.ErrPasteNotFound
		}

		return models.Paste{}, err
	}

	return paste, nil
}

func (s *Storage) SavePaste(ctx context.Context, paste *models.Paste) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	stmt := "INSERT INTO Pastes (title, language, authorid) VALUES ($1, $2, $3) RETURNING id"

	var id string
	err := s.conn.QueryRow(ctx, stmt, paste.Title, paste.Language, paste.AuthorID).Scan(&id)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (s *Storage) DeletePaste(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	stmt := "DELETE FROM Pastes WHERE id = $1"

	_, err := s.conn.Exec(ctx, stmt, id)
	if err != nil {
		return err
	}

	return nil
}
