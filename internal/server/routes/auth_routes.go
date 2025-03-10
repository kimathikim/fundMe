package routes

import (
	"fmt"
	"strings"

	"DBackend/internal/database"
	"DBackend/internal/server/handlers"
	"DBackend/internal/server/middleware"
	"DBackend/utils"

	"github.com/gofiber/fiber/v2"
)

// AuthRoutes registers authentication routes
func AuthRoutes(api fiber.Router, db database.Service) {
	ao := api.Group("/auth", middleware.JWTMiddleware(db))
	authHandler := handlers.NewAuthHandler(db)
	api.Post("/login", authHandler.LoginHandler)
	api.Post("/logout", authHandler.LogoutHandler)
	// ao.Get("/me", authHandler.MeHandler)
	ao.Get("/me", func(c *fiber.Ctx) error {
		tokenString := c.Get("Authorization")
		if tokenString == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing Authorization header"})
		}

		// Remove "Bearer " prefix if present
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		claims, err := utils.ValidateJWT(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": fmt.Sprintf("Invalid token: %v", err)})
		}

		return c.JSON(claims)
	})
}
