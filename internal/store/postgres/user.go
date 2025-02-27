package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"

	"social-api/internal/store"
)

// UserStore implements store.UserStore using PostgreSQL
type UserStore struct {
	db *sql.DB
}

// NewUserStore creates a new PostgreSQL user store
func NewUserStore(db *sql.DB) *UserStore {
	return &UserStore{
		db: db,
	}
}

// Create creates a new user
func (s *UserStore) Create(ctx context.Context, user *store.User) error {
	// SQL query to insert a new user
	query := `
		INSERT INTO users (username, email, password, is_active)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Execute query using the hash from the password field
	err := s.db.QueryRowContext(
		ctx,
		query,
		user.Username,
		user.Email,
		user.Password.Hash, // Use the exported Hash field
		user.IsActive,
	).Scan(
		&user.ID,
		&user.CreatedAt,
	)

	// Check for errors, particularly duplicate keys
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			// Check for unique constraint violations
			switch pqErr.Constraint {
			case "users_email_key":
				return store.ErrDuplicateEmail
			case "users_username_key":
				return store.ErrDuplicateUsername
			}
		}
		return err
	}

	return nil
}

// GetByID retrieves a user by ID
func (s *UserStore) GetByID(ctx context.Context, id int64) (*store.User, error) {
	// SQL query to get a user by ID
	query := `
		SELECT id, username, email, password, is_active, created_at
		FROM users
		WHERE id = $1
	`

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// User to store the result
	var user store.User
	var passwordHash []byte

	// Execute query
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&passwordHash,
		&user.IsActive,
		&user.CreatedAt,
	)

	// Check for errors
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, store.ErrNotFound
		}
		return nil, err
	}

	// Set password hash
	user.Password.Hash = passwordHash

	return &user, nil
}

// GetByEmail retrieves a user by email
func (s *UserStore) GetByEmail(ctx context.Context, email string) (*store.User, error) {
	// SQL query to get a user by email
	query := `
		SELECT id, username, email, password, is_active, created_at
		FROM users
		WHERE email = $1
	`

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// User to store the result
	var user store.User
	var passwordHash []byte

	// Execute query
	err := s.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&passwordHash,
		&user.IsActive,
		&user.CreatedAt,
	)

	// Check for errors
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, store.ErrNotFound
		}
		return nil, err
	}

	// Set password hash
	user.Password.Hash = passwordHash

	return &user, nil
}

// GetByUsername retrieves a user by username
func (s *UserStore) GetByUsername(ctx context.Context, username string) (*store.User, error) {
	// SQL query to get a user by username
	query := `
		SELECT id, username, email, password, is_active, created_at
		FROM users
		WHERE username = $1
	`

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// User to store the result
	var user store.User
	var passwordHash []byte

	// Execute query
	err := s.db.QueryRowContext(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&passwordHash,
		&user.IsActive,
		&user.CreatedAt,
	)

	// Check for errors
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, store.ErrNotFound
		}
		return nil, err
	}

	// Set password hash
	user.Password.Hash = passwordHash

	return &user, nil
}

// Update updates a user
func (s *UserStore) Update(ctx context.Context, user *store.User) error {
	// SQL query to update a user
	query := `
		UPDATE users
		SET username = $1, email = $2, is_active = $3
		WHERE id = $4
	`

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Execute query
	result, err := s.db.ExecContext(
		ctx,
		query,
		user.Username,
		user.Email,
		user.IsActive,
		user.ID,
	)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			// Check for unique constraint violations
			switch pqErr.Constraint {
			case "users_email_key":
				return store.ErrDuplicateEmail
			case "users_username_key":
				return store.ErrDuplicateUsername
			}
		}
		return err
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return store.ErrNotFound
	}

	return nil
}
