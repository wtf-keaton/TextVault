package middleware

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

var (
	ErrMissingAuthorizationHeader = errors.New("missing authorization header")
	ErrInvalidAuthorizationFormat = errors.New("invalid authorization format")
)

func ExtractToken(c *fiber.Ctx) (string, error) {
	authHeader := c.Get("Authorization")

	if authHeader == "" {
		return "", ErrMissingAuthorizationHeader
	}

	if len(authHeader) <= 7 || authHeader[:7] != "Bearer " {
		return "", ErrInvalidAuthorizationFormat
	}

	tokenString := authHeader[7:]

	return tokenString, nil
}
