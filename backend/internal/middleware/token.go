package middleware

import (
	"TextVault/internal/lib/jwt"
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

func GetUserIDFromToken(tokenString string) (int64, error) {
	token, err := jwt.ValidateToken(tokenString)
	if err != nil {
		return 0, err
	}

	claims, err := jwt.ExtractUserClaims(token)
	if err != nil {
		return 0, err
	}

	return claims.ID, nil
}
