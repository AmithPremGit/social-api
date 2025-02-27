package store

import (
	"context"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Common errors for user operations
var (
	ErrNotFound          = errors.New("record not found")
	ErrDuplicateEmail    = errors.New("email address already in use")
	ErrDuplicateUsername = errors.New("username already in use")
)

// User represents a user in the system
type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  Password  `json:"-"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

// Password is a wrapper for user passwords
type Password struct {
	Plaintext *string `json:"-"` // Stores the plaintext password temporarily
	Hash      []byte  `json:"-"` // Stores the hashed password
}

// Set sets a password by hashing the plaintext
func (p *Password) Set(plaintext string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintext), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	p.Plaintext = &plaintext
	p.Hash = hash
	return nil
}

// Matches checks if a plaintext password matches the hash
func (p *Password) Matches(plaintext string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.Hash, []byte(plaintext))
	if err != nil {
		switch {
		case err == bcrypt.ErrMismatchedHashAndPassword:
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}

// UserStore defines the interface for user operations
type UserStore interface {
	// Create creates a new user
	Create(ctx context.Context, user *User) error

	// GetByID retrieves a user by ID
	GetByID(ctx context.Context, id int64) (*User, error)

	// GetByEmail retrieves a user by email
	GetByEmail(ctx context.Context, email string) (*User, error)

	// GetByUsername retrieves a user by username
	GetByUsername(ctx context.Context, username string) (*User, error)

	// Update updates a user
	Update(ctx context.Context, user *User) error
}
