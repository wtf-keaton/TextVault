package jwt

import (
	"TextVault/internal/storage/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type UserClaims struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
}

func NewToken(user models.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = user.ID
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	tokenString, err := token.SignedString([]byte("secret_signing_key"))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte("secret_signing_key"), nil
	})
}

func ExtractUserClaims(token *jwt.Token) (*UserClaims, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, jwt.ErrInvalidKey
	}

	id, ok := claims["id"].(float64)
	if !ok {
		return nil, jwt.ErrInvalidKey
	}

	email, ok := claims["email"].(string)
	if !ok {
		return nil, jwt.ErrInvalidKey
	}

	return &UserClaims{
		ID:    int64(id),
		Email: email,
	}, nil
}
