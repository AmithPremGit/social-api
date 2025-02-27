package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"social-api/internal/model"
	"social-api/internal/store"
)

// PostStore implements store.PostStore using PostgreSQL
type PostStore struct {
	db *sql.DB
}

// NewPostStore creates a new PostgreSQL post store
func NewPostStore(db *sql.DB) *PostStore {
	return &PostStore{
		db: db,
	}
}

// Create creates a new post
func (s *PostStore) Create(ctx context.Context, post *store.Post) error {
	// SQL query to insert a new post
	query := `
		INSERT INTO posts (title, content, user_id)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Execute query
	err := s.db.QueryRowContext(
		ctx,
		query,
		post.Title,
		post.Content,
		post.UserID,
	).Scan(
		&post.ID,
		&post.CreatedAt,
		&post.UpdatedAt,
	)

	// Check for errors
	if err != nil {
		return err
	}

	return nil
}

// GetByID retrieves a post by ID
func (s *PostStore) GetByID(ctx context.Context, id int64) (*store.Post, error) {
	// SQL query to get a post by ID with user information
	query := `
		SELECT 
			p.id, p.title, p.content, p.user_id, p.created_at, p.updated_at,
			u.id, u.username, u.email, u.is_active, u.created_at
		FROM posts p
		JOIN users u ON p.user_id = u.id
		WHERE p.id = $1
	`

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Post and user to store the result
	var post store.Post
	var user store.User

	// Execute query
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		&post.UserID,
		&post.CreatedAt,
		&post.UpdatedAt,
		&user.ID,
		&user.Username,
		&user.Email,
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

	// Set user
	post.User = &user

	return &post, nil
}

// Update updates a post
func (s *PostStore) Update(ctx context.Context, post *store.Post) error {
	// SQL query to update a post
	query := `
		UPDATE posts
		SET title = $1, content = $2
		WHERE id = $3 AND user_id = $4
		RETURNING updated_at
	`

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Execute query
	err := s.db.QueryRowContext(
		ctx,
		query,
		post.Title,
		post.Content,
		post.ID,
		post.UserID,
	).Scan(&post.UpdatedAt)

	// Check for errors
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return store.ErrNotFound
		}
		return err
	}

	return nil
}

// Delete deletes a post
func (s *PostStore) Delete(ctx context.Context, id int64) error {
	// SQL query to delete a post
	query := `DELETE FROM posts WHERE id = $1`

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Execute query
	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
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

// List retrieves a list of posts
func (s *PostStore) List(ctx context.Context, pagination model.Pagination, filter model.PostFilter) ([]*store.Post, int, error) {
	// Base query for listing posts
	baseQuery := `
		SELECT 
			p.id, p.title, p.content, p.user_id, p.created_at, p.updated_at,
			u.id, u.username, u.email, u.is_active, u.created_at
		FROM posts p
		JOIN users u ON p.user_id = u.id
	`

	// Count query for total records
	countQuery := `
		SELECT COUNT(*) FROM posts p
	`

	// Build where clause
	whereClause, args := s.buildWhereClause(filter)
	if whereClause != "" {
		baseQuery += " WHERE " + whereClause
		countQuery += " WHERE " + whereClause
	}

	// Add order by clause
	baseQuery += fmt.Sprintf(" ORDER BY p.%s %s", pagination.SortBy, pagination.Sort)

	// Add pagination
	baseQuery += " LIMIT $" + fmt.Sprint(len(args)+1) + " OFFSET $" + fmt.Sprint(len(args)+2)
	args = append(args, pagination.PageSize, pagination.GetOffset())

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Query for total count
	var totalCount int
	err := s.db.QueryRowContext(ctx, countQuery, args[:len(args)-2]...).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	// Query for posts
	rows, err := s.db.QueryContext(ctx, baseQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	// Process rows
	posts := []*store.Post{}
	for rows.Next() {
		var post store.Post
		var user store.User

		err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.Content,
			&post.UserID,
			&post.CreatedAt,
			&post.UpdatedAt,
			&user.ID,
			&user.Username,
			&user.Email,
			&user.IsActive,
			&user.CreatedAt,
		)
		if err != nil {
			return nil, 0, err
		}

		// Set user
		post.User = &user
		posts = append(posts, &post)
	}

	// Check for errors in row iteration
	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return posts, totalCount, nil
}

// buildWhereClause builds a WHERE clause based on filters
func (s *PostStore) buildWhereClause(filter model.PostFilter) (string, []interface{}) {
	var clauses []string
	var args []interface{}
	var paramCount int

	// Add user ID filter
	if filter.UserID != nil {
		paramCount++
		clauses = append(clauses, fmt.Sprintf("p.user_id = $%d", paramCount))
		args = append(args, *filter.UserID)
	}

	// Add title filter (using ILIKE for case-insensitive search)
	if filter.Title != nil && *filter.Title != "" {
		paramCount++
		clauses = append(clauses, fmt.Sprintf("p.title ILIKE $%d", paramCount))
		args = append(args, "%"+*filter.Title+"%")
	}

	// Add content filter (using ILIKE for case-insensitive search)
	if filter.Content != nil && *filter.Content != "" {
		paramCount++
		clauses = append(clauses, fmt.Sprintf("p.content ILIKE $%d", paramCount))
		args = append(args, "%"+*filter.Content+"%")
	}

	// Add from date filter
	if filter.FromDate != nil {
		paramCount++
		clauses = append(clauses, fmt.Sprintf("p.created_at >= $%d", paramCount))
		args = append(args, *filter.FromDate)
	}

	// Add to date filter
	if filter.ToDate != nil {
		paramCount++
		clauses = append(clauses, fmt.Sprintf("p.created_at <= $%d", paramCount))
		args = append(args, *filter.ToDate)
	}

	// Join all clauses with AND
	whereClause := strings.Join(clauses, " AND ")

	return whereClause, args
}
