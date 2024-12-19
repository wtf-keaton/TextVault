package postgres

import (
	"TextVault/internal/storage/models"
	"context"
)

func (s *Storage) GetPaste(ctx context.Context, hash string) (*models.Paste, error) {
	return nil, nil
}

func (s *Storage) SavePaste(ctx context.Context, paste *models.Paste, content []byte) error {

	return nil
}

func (s *Storage) GetPastesByUser(ctx context.Context, userID int64) ([]*models.Paste, error) {
	return nil, nil
}
