package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"social-api/internal/model"
	"social-api/internal/store"
)

// contextKey is a custom type for context keys
type contextKey string

// userContextKey is the key for storing the user in the request context
const userContextKey = contextKey("user")

// Middleware is the authentication middleware
func Middleware(authenticator *JWTAuthenticator, userStore store.UserStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				model.WriteJSON(w, http.StatusUnauthorized, model.ErrorResponse{
					Error: "authorization header is required",
				})
				return
			}

			// Check for Bearer token
			headerParts := strings.Split(authHeader, " ")
			if len(headerParts) != 2 || headerParts[0] != "Bearer" {
				model.WriteJSON(w, http.StatusUnauthorized, model.ErrorResponse{
					Error: "authorization header format must be Bearer {token}",
				})
				return
			}

			// Get token
			token := headerParts[1]

			// Validate token
			userID, err := authenticator.ValidateToken(token)
			if err != nil {
				switch {
				case errors.Is(err, ErrExpiredToken):
					model.WriteJSON(w, http.StatusUnauthorized, model.ErrorResponse{
						Error: "token has expired",
					})
				default:
					model.WriteJSON(w, http.StatusUnauthorized, model.ErrorResponse{
						Error: "invalid authentication token",
					})
				}
				return
			}

			// Get user from store
			user, err := userStore.GetByID(r.Context(), userID)
			if err != nil {
				if errors.Is(err, store.ErrNotFound) {
					model.WriteJSON(w, http.StatusUnauthorized, model.ErrorResponse{
						Error: "user not found",
					})
				} else {
					model.WriteJSON(w, http.StatusInternalServerError, model.ErrorResponse{
						Error: "internal server error",
					})
				}
				return
			}

			// Check if user is active
			if !user.IsActive {
				model.WriteJSON(w, http.StatusForbidden, model.ErrorResponse{
					Error: "user account is inactive",
				})
				return
			}

			// Add user to context
			ctx := context.WithValue(r.Context(), userContextKey, user)

			// Call next handler with updated context
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserFromContext retrieves the user from the context
func GetUserFromContext(ctx context.Context) (*store.User, bool) {
	user, ok := ctx.Value(userContextKey).(*store.User)
	return user, ok
}

// RequireUser is a middleware that requires a user to be authenticated
func RequireUser(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, ok := GetUserFromContext(r.Context())
		if !ok {
			model.WriteJSON(w, http.StatusUnauthorized, model.ErrorResponse{
				Error: "unauthorized",
			})
			return
		}
		next(w, r)
	}
}
