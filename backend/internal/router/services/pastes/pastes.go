package pastes

import (
	"TextVault/internal/lib/jwt"
	"TextVault/internal/lib/log/sl"
	"TextVault/internal/middleware"
	"TextVault/internal/storage"
	"TextVault/internal/storage/models"
	"TextVault/pkg/random"
	"context"
	"encoding/json"
	"errors"
	"log/slog"

	"github.com/gofiber/fiber/v2"
)

type Service struct {
	pasteSaver    PasteSaver
	pasteGetter   PasteGetter
	pasteProvider PasteProvider
	cacheProvider CacheProvider

	log *slog.Logger
}

// PasteSaver is an interface that provides methods for saving and deleting pastes to the database.
type PasteSaver interface {
	SavePaste(ctx context.Context, paste *models.Paste, content []byte) error
	DeletePaste(ctx context.Context, hash string) error
}

// PasteProvider is an interface that provides methods for uploading, downloading, and deleting pastes from s3 storage.
type PasteProvider interface {
	UploadPaste(ctx context.Context, objectKey string, content []byte) error
	GetPasteContent(ctx context.Context, objectKey string) ([]byte, error)
	DeletePaste(ctx context.Context, objectKey string) error
}

// PasteGetter is an interface that provides methods for getting pastes from the database.
type PasteGetter interface {
	GetPaste(ctx context.Context, hash string) (models.Paste, error)
}

type CacheProvider interface {
	Set(ctx context.Context, key string, value string) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
}

// pasteBody is a struct that represents the request body for saving a new paste.
type pasteBody struct {
	Title    string `json:"title"`
	Language string `json:"language"`
	Content  string `json:"content"`
}

// New creates a new paste service.
func New(log *slog.Logger,
	pasteSaver PasteSaver,
	pasteGetter PasteGetter,
	pasteProvider PasteProvider,
	cacheProvider CacheProvider,
) *Service {
	return &Service{
		pasteSaver:    pasteSaver,
		pasteGetter:   pasteGetter,
		pasteProvider: pasteProvider,
		cacheProvider: cacheProvider,
		log:           log,
	}
}

// SavePaste saves a new paste to the database and upload content to s3 storage. If the request contains a valid
// authorization token, the paste's author ID is set to the user ID extracted from
// the token. Otherwise, the author ID is set to 0 (anonymous user). The response
// body contains the hash of the saved paste.
func (s *Service) SavePaste(c *fiber.Ctx) error {
	const prefix = "internal.router.services.paste.SavePaste"

	tokenString, err := middleware.ExtractToken(c)

	p := new(pasteBody)

	if err := c.BodyParser(p); err != nil {
		s.log.Error("Failed to parse save request", sl.Err(err))

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "content is required",
		})
	}

	log := s.log.With(
		slog.String("op", prefix),
		slog.String("title", p.Title),
	)

	log.Info("Attempting to save paste")

	log.Debug("Paste content", slog.String("content", p.Content))

	var AuthorID int64 = 0
	if err == nil {
		userID, err := middleware.GetUserIDFromToken(tokenString)
		if err == nil {
			AuthorID = userID
			log.Info("User ID extracted from token", slog.Int64("user_id", userID))
		} else {
			log.Warn("Failed to extract user ID from token", sl.Err(err))
		}
	}

	pasteHash := random.String(16)
	log.Info("Saving paste", slog.String("hash", pasteHash), slog.String("title", p.Title))

	pasteModel := &models.Paste{
		Title:    p.Title,
		Hash:     pasteHash,
		Language: p.Language,
		AuthorID: AuthorID, // If token is not valid, AuthorID will be 0
	}

	err = s.pasteSaver.SavePaste(c.Context(), pasteModel, []byte(p.Content))
	if err != nil {
		log.Error("Failed to save paste", sl.Err(err))

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "failed to save paste",
		})
	}

	err = s.pasteProvider.UploadPaste(c.Context(), pasteHash, []byte(p.Content))
	if err != nil {
		log.Error("Failed to upload paste", sl.Err(err))

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "failed to upload paste",
		})
	}

	log.Info("Paste saved successfully", slog.String("hash", pasteHash))

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"hash": pasteHash,
	})
}

// GetPaste retrieves a paste from the database and its content from S3 storage based on the provided hash.
// If the paste is not found, it returns a 401 Unauthorized status with an error message.
// If any other error occurs during retrieval, it returns a 500 Internal Server Error status with an error message.
// On successful retrieval, it sends the paste content as a string in the response.
func (s *Service) GetPaste(c *fiber.Ctx) error {
	const prefix = "internal.router.services.paste.GetPaste"
	hash := c.Params("hash")

	log := s.log.With(
		slog.String("op", prefix),
		slog.String("hash", hash),
	)

	log.Info("Attempting to get paste")

	cacheContent, err := s.cacheProvider.Get(c.Context(), hash)
	if err == nil {
		pasteResponse := pasteBody{}

		err := json.Unmarshal([]byte(cacheContent), &pasteResponse)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Internal server error",
			})
		}

		log.Info("Paste retrieved successfully", slog.String("hash", hash))

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"title":    pasteResponse.Title,
			"language": pasteResponse.Language,
			"content":  pasteResponse.Content,
		})
	}

	paste, err := s.pasteGetter.GetPaste(c.Context(), hash)
	if err != nil {
		if errors.Is(err, storage.ErrPasteNotFound) {
			s.log.Warn("Failed to find paste")

			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "paste not found",
			})
		}

		s.log.Error("Failed to get paste", sl.Err(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	content, err := s.pasteProvider.GetPasteContent(c.Context(), paste.Hash)
	if err != nil {
		log.Error("Failed to get paste content", sl.Err(err))

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "failed to get paste",
		})
	}

	log.Info("Paste retrieved successfully", slog.String("hash", paste.Hash))

	pasteResponse := pasteBody{
		Title:    paste.Title,
		Language: paste.Language,
		Content:  string(content),
	}

	var cacheData []byte
	cacheData, err = json.Marshal(pasteResponse)
	if err != nil {
		log.Error("Failed to marshal paste response", sl.Err(err))

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	err = s.cacheProvider.Set(c.Context(), paste.Hash, string(cacheData))
	if err != nil {
		log.Error("Failed to set cache", sl.Err(err))
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"language": paste.Language,
		"content":  pasteResponse.Content,
	})
}

// DeletePaste deletes a paste from the database and s3 storage based on the provided hash.
// If the paste is not found, it returns a 401 Unauthorized status with an error message.
// If any other error occurs during deletion, it returns a 500 Internal Server Error status with an error message.
// On successful deletion, it returns a 200 OK status with an empty response body.
func (s *Service) DeletePaste(c *fiber.Ctx) error {
	const prefix = "internal.router.services.paste.DeletePaste"
	hash := c.Params("hash")

	tokenString, err := middleware.ExtractToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	token, err := jwt.ValidateToken(tokenString)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	claims, err := jwt.ExtractUserClaims(token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	log := s.log.With(
		slog.String("op", prefix),
		slog.String("hash", hash),
		slog.String("mail", claims.Email),
	)

	log.Info("Attempting to delete paste")

	paste, err := s.pasteGetter.GetPaste(c.Context(), hash)
	if err != nil {
		if errors.Is(err, storage.ErrPasteNotFound) {
			s.log.Warn("Failed to find paste")

			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "paste not found",
			})
		}

		s.log.Error("Failed to get paste", sl.Err(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	if paste.AuthorID != claims.ID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "you are not owner of this paste",
		})
	}

	// @NOTE: Delete paste from db
	err = s.pasteSaver.DeletePaste(c.Context(), hash)
	if err != nil {
		log.Error("Failed to delete paste", sl.Err(err))

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "failed to delete paste",
		})
	}

	// @NOTE: Delete paste from s3
	err = s.pasteProvider.DeletePaste(c.Context(), hash)
	if err != nil {
		log.Error("Failed to delete paste from s3 storage", sl.Err(err))

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "failed to delete paste",
		})
	}

	return c.SendStatus(fiber.StatusOK)
}
