package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// CORSMiddleware configures CORS for the application
func CORSMiddleware() fiber.Handler {
	return cors.New(cors.Config{
    AllowOrigins:     "http://localhost:3000", // Change this to specific frontend domains in production
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders:     "Accept, Authorization, Content-Type",
		AllowCredentials: true, // Enables cookies (for authentication)
		MaxAge:           3200, // Cache the preflight response for 1 hour
	})
}

