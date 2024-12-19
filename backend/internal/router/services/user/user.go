package user

import (
	"TextVault/internal/lib/jwt"
	"TextVault/internal/middleware"
	"TextVault/internal/storage/models"
	"TextVault/pkg/passwordhash"
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
)

type Service struct {
	userSaver  UserSaver
	userGetter UserProvider
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

func New(userSaver UserSaver, userGetter UserProvider) *Service {
	return &Service{
		userSaver:  userSaver,
		userGetter: userGetter,
	}
}

func (s *Service) Login(c *fiber.Ctx) error {
	const prefix = "internal.router.services.user.Login"

	p := new(loginRequest)

	if err := c.BodyParser(p); err != nil {
		log.Fatalln(prefix, ": Failed to parse login request: ", err)
		return err
	}

	if len(p.Password) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "password is required",
		})
	}

	user, err := s.userGetter.GetUser(c.Context(), p.Username)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "user not found",
		})
	}

	if !passwordhash.Validate(p.Password, user.PasswordHash) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid credentials",
		})
	}

	token, err := jwt.NewToken(user)
	if err != nil {
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
		log.Fatalln(prefix, ": Failed to parse register request: ", err)
		return err
	}

	if len(p.Password) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "password is required",
		})
	}

	passwordHash, err := passwordhash.New(p.Password)
	if err != nil {
		log.Fatalln(prefix, ": Failed to hash password: ", err)

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	id, err := s.userSaver.SaveUser(c.Context(), p.Username, p.Mail, passwordHash)
	if err != nil {
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
