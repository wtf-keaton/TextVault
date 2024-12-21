package user

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
	userSaver  UserSaver
	userGetter UserProvider
	log        *slog.Logger
}

type UserSaver interface {
	SaveUser(ctx context.Context, username, email, password string) (int64, error)
}

type UserProvider interface {
	GetUser(ctx context.Context, username string) (models.User, error)
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

func New(log *slog.Logger, userSaver UserSaver, userGetter UserProvider) *Service {
	return &Service{
		userSaver:  userSaver,
		userGetter: userGetter,
		log:        log,
	}
}

func (s *Service) Login(c *fiber.Ctx) error {
	const prefix = "internal.router.services.user.Login"

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

	user, err := s.userGetter.GetUser(c.Context(), p.Username)
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

func (s *Service) Register(c *fiber.Ctx) error {
	const prefix = "internal.router.services.user.Register"

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

	id, err := s.userSaver.SaveUser(c.Context(), p.Username, p.Mail, passwordHash)
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

func (s *Service) ValidateToken(ctx *fiber.Ctx) error {
	tokenString, _ := middleware.ExtractToken(ctx)

	token, err := jwt.ValidateToken(tokenString)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid token")
	}

	claims, err := jwt.ExtractUserClaims(token)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid token claims")
	}

	return ctx.JSON(fiber.Map{
		"user_id": claims.ID,
		"email":   claims.Email,
	})
}
