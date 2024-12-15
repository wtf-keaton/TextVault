package storage

import "errors"

var (
	ErrPasteNotFound = errors.New("paste not found")
	ErrUserNotFound  = errors.New("user not found")
	ErrIncorrectPass = errors.New("incorrect password")
)
