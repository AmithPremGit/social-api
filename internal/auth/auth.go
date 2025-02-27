package auth

import (
	"errors"
	"time"
)

// Common authentication errors
var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("expired token")
)

// Authenticator defines the interface for authentication operations
type Authenticator interface {
	// GenerateToken generates a token for a user
	GenerateToken(userID int64, expiry time.Duration) (string, error)

	// ValidateToken validates a token and returns the user ID
	ValidateToken(token string) (int64, error)
}
