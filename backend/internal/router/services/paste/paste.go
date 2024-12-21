package paste

import (
	"TextVault/internal/middleware"
	"TextVault/internal/storage/models"
	"TextVault/pkg/random"
	"context"
	"log/slog"

	"github.com/gofiber/fiber/v2"
)

type Service struct {
	pasteSaver  PasteSaver
	pasteGetter PasteGetter
	log         *slog.Logger
}

type PasteSaver interface {
	SavePaste(ctx context.Context, paste *models.Paste, content []byte) error
}

type PasteGetter interface {
}

func New(log *slog.Logger, pasteSaver PasteSaver, pasteGetter PasteGetter) *Service {
	return &Service{
		pasteSaver:  pasteSaver,
		pasteGetter: pasteGetter,
		log:         log,
	}
}

// SavePaste saves a new paste to the database and upload content to s3 storage. If the request contains a valid
// authorization token, the paste's author ID is set to the user ID extracted from
// the token. Otherwise, the author ID is set to 0 (anonymous user). The response
// body contains the hash of the saved paste.

func (s *Service) SavePaste(c *fiber.Ctx) error {
	tokenString, err := middleware.ExtractToken(c)

	var AuthorID int64 = -1
	if err == nil {
		userID, err := middleware.GetUserIDFromToken(tokenString)
		if err == nil {
			AuthorID = userID
		}
	}

	pasteHash := random.String(16)

	pasteModel := &models.Paste{
		Title:    c.FormValue("title"),
		Hash:     pasteHash,
		AuthorID: AuthorID, // If token is not valid, AuthorID will be -1
	}

	err = s.pasteSaver.SavePaste(c.Context(), pasteModel, []byte(c.FormValue("content")))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "failed to save paste",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"hash": pasteHash,
	})
}
