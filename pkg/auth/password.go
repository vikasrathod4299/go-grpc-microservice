package auth

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

const (
	MinimumPasswordLength = 8
	MaximumPasswordLength = 72
)

var (
	ErrPasswordTooShort  = errors.New("password is too short, must be at least 8 characters")
	ErrPasswordTooLong   = errors.New("password is too long, must be at most 72 characters")
	ErrPasswordMissmatch = errors.New("password does not match the hash")
)

func HashPassword(password string) (string, error) {
	if len(password) < MinimumPasswordLength {
		return "", ErrPasswordTooShort
	}
	if len([]byte(password)) > MaximumPasswordLength {
		return "", ErrPasswordTooLong
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func ComparePassword(password, passwordHash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	if err != nil {
		return ErrPasswordMissmatch
	}
	return nil
}
