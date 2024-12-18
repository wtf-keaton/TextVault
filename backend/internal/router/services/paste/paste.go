package paste

import (
	"TextVault/internal/storage/models"
	"context"
)

type Service struct {
	pasteSaver  PasteSaver
	pasteGetter PasteGetter
}

type PasteSaver interface {
	SavePaste(ctx context.Context, paste *models.Paste) error
}

type PasteGetter interface {
}

func New(pasteSaver PasteSaver, pasteGetter PasteGetter) *Service {
	return &Service{
		pasteSaver:  pasteSaver,
		pasteGetter: pasteGetter,
	}
}
