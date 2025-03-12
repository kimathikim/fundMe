package utils

import (
	"fmt"
	"strings"
	"time"

	//  "encoding/base64"

	"github.com/golang-jwt/jwt/v4"
)

// Secret key for signing JWT (Use an environment variable in production!)
var jwtSecret = []byte("your-secure-secret-key")

// Claims struct for JWT payload
type Claims struct {
	UserID string   `json:"user_id"`
	Roles  []string `json:"roles"`
	jwt.RegisteredClaims
}

// GenerateJWT creates a signed JWT token for a user
func GenerateJWT(user_id string, roles []string) (string, error) {
	// Set token expiration time (e.g., 15 minutes)
	expirationTime := time.Now().Add(60 * time.Minute)

	// Define claims (payload)
	claims := &Claims{
		UserID: user_id,
		Roles:  roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime), // Expiry
		},
	}

	// Create token with claims & sign it
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ExtractBearerToken removes "Bearer " prefix from Authorization header
func ExtractBearerToken(authHeader string) string {
	if strings.HasPrefix(authHeader, "Bearer ") {
		return strings.TrimPrefix(authHeader, "Bearer ")
	}
	return ""
}

// ValidateJWT verifies and extracts claims from a JWT token
func ValidateJWT(tokenString string) (*Claims, error) {
	claims := &Claims{}

	// Parse the token and validate it
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	// Return error if token is invalid
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token: %v", err)
	}

	return claims, nil
}
