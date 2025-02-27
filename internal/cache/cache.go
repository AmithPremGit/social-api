package cache

import (
	"context"
	"time"
)

// Cache defines the interface for cache operations
type Cache interface {
	// Get retrieves a value from the cache
	Get(ctx context.Context, key string, destination interface{}) error

	// Set stores a value in the cache with an expiration
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error

	// Delete removes a value from the cache
	Delete(ctx context.Context, key string) error
}

// Helper functions for generating cache keys

// UserCacheKey generates a cache key for a user
func UserCacheKey(userID int64) string {
	return "user:" + string(userID)
}

// PostCacheKey generates a cache key for a post
func PostCacheKey(postID int64) string {
	return "post:" + string(postID)
}
