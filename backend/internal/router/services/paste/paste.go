package paste

import (
	"TextVault/internal/lib/jwt"
	"TextVault/internal/middleware"
	"TextVault/internal/storage/models"
	"TextVault/pkg/random"
	"context"

	"github.com/gofiber/fiber/v2"
)

type Service struct {
	pasteSaver  PasteSaver
	pasteGetter PasteGetter
}

type PasteSaver interface {
	SavePaste(ctx context.Context, paste *models.Paste, content []byte) error
}

type PasteGetter interface {
}

func New(pasteSaver PasteSaver, pasteGetter PasteGetter) *Service {
	return &Service{
		pasteSaver:  pasteSaver,
		pasteGetter: pasteGetter,
	}
}

func (s *Service) SavePaste(c *fiber.Ctx) error {
	tokenString, err := middleware.ExtractToken(c)

	// TODO: Refactor this shit
	var AuthorID int64 = 0
	if err != nil {
		AuthorID = 0
	} else {
		token, err := jwt.ValidateToken(tokenString)
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid token")
		}

		claims, err := jwt.ExtractUserClaims(token)
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid token claims")
		}

		AuthorID = claims.ID
	}

	pasteHash := random.String(16)

	pasteModel := &models.Paste{
		Title:    c.FormValue("title"),
		Hash:     pasteHash,
		AuthorID: AuthorID,
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
