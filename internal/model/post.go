package model

import "time"

// PostInput represents input for post creation
type PostInput struct {
	Title   string `json:"title" validate:"required,min=3,max=200"`
	Content string `json:"content" validate:"required,min=10"`
}

// PostUpdateInput represents input for post update
type PostUpdateInput struct {
	Title   *string `json:"title" validate:"omitempty,min=3,max=200"`
	Content *string `json:"content" validate:"omitempty,min=10"`
}

// PostResponse represents a post in responses
type PostResponse struct {
	ID        int64        `json:"id"`
	Title     string       `json:"title"`
	Content   string       `json:"content"`
	User      UserResponse `json:"user"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

// PostFilter represents filters for post queries
type PostFilter struct {
	UserID   *int64     `json:"user_id,omitempty"`
	Title    *string    `json:"title,omitempty"`
	Content  *string    `json:"content,omitempty"`
	FromDate *time.Time `json:"from_date,omitempty"`
	ToDate   *time.Time `json:"to_date,omitempty"`
}
