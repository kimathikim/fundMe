package utils

import (
	"strings"
	"time"
  "fmt"
//  "encoding/base64"

	"github.com/golang-jwt/jwt/v4"
)

// Secret key for signing JWT (Use an environment variable in production!)
var jwtSecret = []byte("your-secure-secret-key")

// Claims struct for JWT payload
type Claims struct {
	Email string   `json:"email"`
	Roles []string `json:"roles"`
	jwt.RegisteredClaims
}

// GenerateJWT creates a signed JWT token for a user
func GenerateJWT(email string, roles []string) (string, error) {
	// Set token expiration time (e.g., 15 minutes)
	expirationTime := time.Now().Add(15 * time.Minute)

	// Define claims (payload)
	claims := &Claims{
		Email: email,
		Roles: roles,
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
