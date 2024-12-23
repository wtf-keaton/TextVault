package postgres

import (
	"TextVault/internal/storage"
	"TextVault/internal/storage/models"
	"context"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
)

// SaveUser creates a new user in the database and returns the ID of the newly created user.
// If any error occurs during the execution of the function, it is wrapped in a custom error
// with the prefix "storage.postgresql.SaveUser" and returned.
// The context is used to set a timeout for the execution of the function, which is set to the
// value of the timeout field of the receiver.
func (s *Storage) SaveUser(ctx context.Context, username, email, password string) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	stmt := "INSERT INTO Users (username, email, passwordhash) VALUES ($1, $2, $3) RETURNING id"

	var id int64
	err := s.conn.QueryRow(ctx, stmt, username, email, password).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// GetUser retrieves a User from the database based on the provided username.
// If the username is not found, the function returns ErrUserNotFound.
// If any other error occurs, the function returns the error wrapped in a custom error
// with the prefix "storage.postgresql.GetUser".
func (s *Storage) GetUser(ctx context.Context, username string) (models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	stmt := "SELECT ID, username, email, passwordhash FROM Users WHERE username = $1 OR email = $1"

	var user models.User
	err := pgxscan.Get(ctx, s.conn, &user, stmt, username)
	if err != nil {
		if err == pgx.ErrNoRows {
			return models.User{}, storage.ErrUserNotFound
		}

		return models.User{}, err
	}

	return user, nil
}

func (s *Storage) GetUserPastes(ctx context.Context, userID int64) ([]models.Paste, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	stmt := "SELECT ID, title, hash, language, authorid FROM Pastes WHERE authorid = $1"

	var pastes []models.Paste
	err := pgxscan.Select(ctx, s.conn, &pastes, stmt, userID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return []models.Paste{}, storage.ErrUserDontHavePastes
		}

		return nil, err
	}

	return pastes, nil
}

func (s *Storage) UpdateUser(ctx context.Context, user *models.User) error {
	return nil
}
