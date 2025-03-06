package routes

import (
	"DBackend/internal/database"
	"DBackend/internal/server/middleware"
	"github.com/gofiber/fiber/v2"
)

// ProtectedRoutes registers protected routes (role-based)
func ProtectedRoutes(api fiber.Router, db database.Service) {
	protected := api.Group("/protected", middleware.JWTMiddleware(db))

	// admin-only routes
	protected.Get("/admin", middleware.RequireRole("admin"), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Admin route"})
	})

	protected.Get("/investor", middleware.RequireRole("investor"), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Investor route"})
	})

	protected.Get("/founder", middleware.RequireRole("founder"), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Founder route"})
	})
}
