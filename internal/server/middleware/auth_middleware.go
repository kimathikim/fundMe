package middleware

import (
	"fmt"
	"strings"

	"DBackend/internal/database"
	"DBackend/utils"

	"github.com/gofiber/fiber/v2"
)

// JWTMiddleware verifies JWT tokens and checks if they are blacklisted
func JWTMiddleware(db database.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get token from Authorization header or cookie
		token := c.Cookies("jwt")
		if token == "" {
			authHeader := c.Get("Authorization")
			if strings.HasPrefix(authHeader, "Bearer ") {
				token = utils.ExtractBearerToken(authHeader)
			}
		}

		if token == "" {
			return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
		}

		// Check if token is blacklisted
		isBlacklisted, err := db.User().IsTokenBlacklisted(c.Context(), token)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Database error"})
		}

		if isBlacklisted {
			return c.Status(401).JSON(fiber.Map{"error": "Token expired, please log in again"})
		}

		// Validate token
		claims, err := utils.ValidateJWT(token)
		if err != nil {
      fmt.Println(err)
			return c.Status(401).JSON(fiber.Map{"error": "Invalid token"})
		}
		// Store user details in context
		c.Locals("email", claims.Email)
		c.Locals("roles", claims.Roles)

		return c.Next()
	}
}


func RequireRole(role string) fiber.Handler {
  return func(c *fiber.Ctx) error {
    roles, ok := c.Locals("user_roles").([]string)
    if !ok {
      return c.Status(403).JSON(fiber.Map{"error": "Forbidden"})
    }
    for _, r := range roles {
      if r == role {
        return c.Next()
      }
    }
    return c.Status(403).JSON(fiber.Map{"error": "access denied"})
  }
}
