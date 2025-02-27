package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisCache implements the Cache interface using Redis
type RedisCache struct {
	client *redis.Client
}

// NewRedisCache creates a new Redis cache
func NewRedisCache(client *redis.Client) *RedisCache {
	return &RedisCache{
		client: client,
	}
}

// Get retrieves a value from the cache
func (c *RedisCache) Get(ctx context.Context, key string, destination interface{}) error {
	// Get data from Redis
	data, err := c.client.Get(ctx, key).Result()
	if err != nil {
		// Handle Redis Nil error (key not found)
		if err == redis.Nil {
			return fmt.Errorf("key %s not found in cache", key)
		}
		return fmt.Errorf("cache get error: %w", err)
	}

	// Unmarshal data into destination
	err = json.Unmarshal([]byte(data), destination)
	if err != nil {
		return fmt.Errorf("cache unmarshal error: %w", err)
	}

	return nil
}

// Set stores a value in the cache
func (c *RedisCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	// Marshal value to JSON
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("cache marshal error: %w", err)
	}

	// Store data in Redis
	err = c.client.Set(ctx, key, data, ttl).Err()
	if err != nil {
		return fmt.Errorf("cache set error: %w", err)
	}

	return nil
}

// Delete removes a value from the cache
func (c *RedisCache) Delete(ctx context.Context, key string) error {
	// Delete key from Redis
	err := c.client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("cache delete error: %w", err)
	}

	return nil
}

// UserKey generates a cache key for a user
func UserKey(userID int64) string {
	return fmt.Sprintf("user:%d", userID)
}

// PostKey generates a cache key for a post
func PostKey(postID int64) string {
	return fmt.Sprintf("post:%d", postID)
}
