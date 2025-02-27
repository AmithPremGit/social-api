package store

import (
	"context"
	"time"

	"social-api/internal/model"
)

// Post represents a post in the system
type Post struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	UserID    int64     `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	User      *User     `json:"-"`
}

// PostStore defines the interface for post operations
type PostStore interface {
	// Create creates a new post
	Create(ctx context.Context, post *Post) error

	// GetByID retrieves a post by ID
	GetByID(ctx context.Context, id int64) (*Post, error)

	// Update updates a post
	Update(ctx context.Context, post *Post) error

	// Delete deletes a post
	Delete(ctx context.Context, id int64) error

	// List retrieves a list of posts
	List(ctx context.Context, pagination model.Pagination, filter model.PostFilter) ([]*Post, int, error)
}
