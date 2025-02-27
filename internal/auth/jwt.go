package auth

import (
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTAuthenticator implements the Authenticator interface using JWT
type JWTAuthenticator struct {
	secret string
	iss    string
	aud    string
}

// NewJWTAuthenticator creates a new JWT authenticator
func NewJWTAuthenticator(secret, issuer, audience string) *JWTAuthenticator {
	return &JWTAuthenticator{
		secret: secret,
		iss:    issuer,
		aud:    audience,
	}
}

// GenerateToken creates a JWT token for a user
func (a *JWTAuthenticator) GenerateToken(userID int64, expiry time.Duration) (string, error) {
	// Create token claims
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(expiry).Unix(),
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
		"iss": a.iss,
		"aud": a.aud,
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token with secret
	signedToken, err := token.SignedString([]byte(a.secret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

// ValidateToken validates a JWT token and returns the user ID
func (a *JWTAuthenticator) ValidateToken(tokenString string) (int64, error) {
	// Parse token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(a.secret), nil
	},
		jwt.WithExpirationRequired(),
		jwt.WithAudience(a.aud),
		jwt.WithIssuer(a.iss),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
	)

	// Handle parsing errors
	if err != nil {
		return 0, ErrInvalidToken
	}

	// Validate token
	if !token.Valid {
		return 0, ErrInvalidToken
	}

	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, ErrInvalidToken
	}

	// Extract user ID from subject claim
	sub, ok := claims["sub"]
	if !ok {
		return 0, ErrInvalidToken
	}

	// Convert subject to string based on type
	var userIDStr string
	switch v := sub.(type) {
	case float64:
		userIDStr = strconv.FormatFloat(v, 'f', 0, 64)
	case string:
		userIDStr = v
	default:
		return 0, ErrInvalidToken
	}

	// Parse user ID as int64
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return 0, ErrInvalidToken
	}

	return userID, nil
}
