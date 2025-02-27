package middleware

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"social-api/internal/model"
)

// RateLimiter defines the interface for rate limiting
type RateLimiter interface {
	Allow(ip string) (bool, time.Duration)
}

// FixedWindowRateLimiter implements a simple fixed window rate limiter
type FixedWindowRateLimiter struct {
	mu      sync.RWMutex
	clients map[string]*clientLimit
	limit   int
	window  time.Duration
}

// clientLimit tracks request count and window expiry for a client
type clientLimit struct {
	count  int
	expiry time.Time
}

// NewFixedWindowRateLimiter creates a new fixed window rate limiter
func NewFixedWindowRateLimiter(limit int, window time.Duration) *FixedWindowRateLimiter {
	return &FixedWindowRateLimiter{
		clients: make(map[string]*clientLimit),
		limit:   limit,
		window:  window,
	}
}

// Allow checks if a request from an IP should be allowed
func (rl *FixedWindowRateLimiter) Allow(ip string) (bool, time.Duration) {
	rl.mu.RLock()
	client, exists := rl.clients[ip]
	rl.mu.RUnlock()

	now := time.Now()

	// If client exists and window hasn't expired
	if exists && now.Before(client.expiry) {
		rl.mu.Lock()
		defer rl.mu.Unlock()

		// Recheck in case it changed while acquiring lock
		client, exists = rl.clients[ip]
		if exists && now.Before(client.expiry) {
			// If under limit, increment and allow
			if client.count < rl.limit {
				client.count++
				return true, 0
			}

			// Over limit, calculate retry-after
			retryAfter := client.expiry.Sub(now)
			return false, retryAfter
		}
	}

	// Create new window or reset expired window
	rl.mu.Lock()
	defer rl.mu.Unlock()

	rl.clients[ip] = &clientLimit{
		count:  1,
		expiry: now.Add(rl.window),
	}

	return true, 0
}

// RateLimiterMiddleware is middleware for rate limiting requests
func RateLimiterMiddleware(limiter RateLimiter, logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get client IP
			ip := getClientIP(r)

			// Check if request is allowed
			allowed, retryAfter := limiter.Allow(ip)
			if !allowed {
				// Set retry-after header
				w.Header().Set("Retry-After", fmt.Sprintf("%.0f", retryAfter.Seconds()))

				// Log rate limit exceeded
				logger.Printf("Rate limit exceeded for IP: %s", ip)

				// Return rate limit error
				model.WriteJSON(w, http.StatusTooManyRequests, model.ErrorResponse{
					Error: fmt.Sprintf("Rate limit exceeded. Try again in %.0f seconds", retryAfter.Seconds()),
				})
				return
			}

			// Process the request
			next.ServeHTTP(w, r)
		})
	}
}

// getClientIP extracts the client IP from the request
func getClientIP(r *http.Request) string {
	// Check for X-Forwarded-For header
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		return forwarded
	}

	// Fall back to RemoteAddr
	return r.RemoteAddr
}
