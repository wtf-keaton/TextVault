package account

import (
	"TextVault/internal/lib/jwt"
	"TextVault/internal/lib/log/sl"
	"TextVault/internal/middleware"
	"TextVault/internal/storage"
	"TextVault/internal/storage/models"
	"TextVault/pkg/passwordhash"
	"context"
	"errors"
	"log/slog"

	"github.com/gofiber/fiber/v2"
)

type Service struct {
	accountSaver  AccountSaver
	accountGetter AccountGetter
	log           *slog.Logger
}

type AccountSaver interface {
	SaveUser(ctx context.Context, username, email, password string) (int64, error)
}

type AccountGetter interface {
	GetUser(ctx context.Context, username string) (models.User, error)
	GetUserPastes(ctx context.Context, userID int64) ([]models.Paste, error)
}

type loginRequest struct {
	Username string `json:"u"`
	Password string `json:"p"`
}

type registerRequest struct {
	Username string `json:"u"`
	Mail     string `json:"m"`
	Password string `json:"p"`
}

func New(log *slog.Logger, accountSaver AccountSaver, accountGetter AccountGetter) *Service {
	return &Service{
		accountSaver:  accountSaver,
		accountGetter: accountGetter,
		log:           log,
	}
}

// Login authenticates a user by validating their username and password.
// It requires a valid username and password in the request body.
// If the request body is invalid, it returns a 400 Bad Request status with an error message.
// If the username or password is incorrect, it returns a 401 Unauthorized status with an error message.
// If any other error occurs during authentication, it returns a 500 Internal Server Error status with an error message.
// On successful authentication, it returns a 200 OK status with a JWT token in the response.
func (s *Service) Login(c *fiber.Ctx) error {
	const prefix = "internal.router.services.account.Login"

	p := new(loginRequest)

	if err := c.BodyParser(p); err != nil {
		s.log.Error("Failed to parse login request", sl.Err(err))

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "email/username/password is required",
		})
	}

	log := s.log.With(
		slog.String("op", prefix),
		slog.String("username", p.Username),
	)

	log.Info("Attempting to login user")

	if len(p.Password) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "password is required",
		})
	}

	user, err := s.accountGetter.GetUser(c.Context(), p.Username)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			s.log.Warn("Failed to find user")

			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "user not found",
			})
		}

		s.log.Error("Failed to get user", sl.Err(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	if !passwordhash.Validate(p.Password, user.PasswordHash) {
		s.log.Info("invalid credentials")

		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid credentials",
		})
	}

	log.Info("Successfully logged in user")

	token, err := jwt.NewToken(user)
	if err != nil {
		s.log.Error("failed to generate token", sl.Err(err))

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "failed to create token",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"token": token,
	})
}

// Register creates a new user in the database and returns the user ID as a JSON response.
// It requires a valid username, email and password in the request body.
// If the request body is invalid, it returns a 400 Bad Request status with an error message.
// If the password is empty, it returns a 400 Bad Request status with an error message.
// If any other error occurs during registration, it returns a 500 Internal Server Error status with an error message.
// On successful registration, it returns a 200 OK status with the user ID in the response.
func (s *Service) Register(c *fiber.Ctx) error {
	const prefix = "internal.router.services.account.Register"

	p := new(registerRequest)

	if err := c.BodyParser(p); err != nil {
		s.log.Error("Failed to parse register request", sl.Err(err))
		return err
	}

	log := s.log.With(
		slog.String("op", prefix),
		slog.String("username", p.Username),
	)

	if len(p.Password) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "password is required",
		})
	}

	passwordHash, err := passwordhash.New(p.Password)
	if err != nil {
		log.Error("Failed to hash password", sl.Err(err))

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	id, err := s.accountSaver.SaveUser(c.Context(), p.Username, p.Mail, passwordHash)
	if err != nil {
		log.Error("failed to save user", sl.Err(err))

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "failed to save user",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"id": id,
	})
}

// GetUserPastes retrieves all pastes created by a specific user from the database.
// It requires a valid authorization token to authenticate the user and extract their user ID.
// If the token is invalid or missing, it returns a 401 Unauthorized status with an error message.
// If the user does not have any pastes, it returns a 401 Unauthorized status with a specific error message.
// If any other error occurs during retrieval, it returns a 500 Internal Server Error status with an error message.
// On successful retrieval, it returns a 200 OK status with the pastes in the response.
func (s *Service) GetUserPastes(c *fiber.Ctx) error {
	const prefix = "internal.router.services.account.GetUserPastes"

	tokenString, err := middleware.ExtractToken(c)
	if err != nil {
		return s.unauthorizedResponse(c)
	}

	userID, err := middleware.GetUserIDFromToken(tokenString)
	if err != nil {
		return s.unauthorizedResponse(c)
	}

	log := s.log.With(
		slog.String("op", prefix),
		slog.Int64("user_id", userID),
	)

	log.Info("Attempting to get pastes by user")

	pastes, err := s.accountGetter.GetUserPastes(c.Context(), userID)
	if err != nil {
		return s.handleGetPastesError(c, err, log)
	}

	s.log.Info("Successfully got pastes by user", slog.Int("count", len(pastes)), slog.Int64("user_id", userID))

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"pastes": pastes,
	})
}

// unauthorizedResponse returns a 401 Unauthorized status with an error message to the client.
// It is used in various places in the service to return an error when the user is not authenticated.
// The response body contains a JSON object with a single key-value pair, where the key is "error" and the value is "unauthorized".
func (s *Service) unauthorizedResponse(c *fiber.Ctx) error {
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"error": "unauthorized",
	})
}

func (s *Service) handleGetPastesError(c *fiber.Ctx, err error, log *slog.Logger) error {
	switch {
	case errors.Is(err, storage.ErrUserDontHavePastes):
		log.Warn("Failed to find pastes")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "user doesn't have pastes",
		})
	default:
		log.Error("Failed to get pastes", sl.Err(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}
}
