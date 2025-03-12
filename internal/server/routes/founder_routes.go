package routes

import (
	"DBackend/internal/database"
	"DBackend/internal/server/handlers"
	"DBackend/internal/server/middleware"
	"github.com/gofiber/fiber/v2"
)

// fousnderRoutes registers founder routes (role-based)
func FounderRoutes(api fiber.Router, db database.Service) {
	founder := api.Group("/founder", middleware.JWTMiddleware(db))
	founderHandler := handlers.NewFounderHandler(db)
	userHandler := handlers.NewUserHandler(db)
	founder.Patch("/profile", middleware.RequireRole("founder"), founderHandler.UpdateFounderHandler)
	founder.Get("/profile", middleware.RequireRole("founder"), founderHandler.GetFounderDetailsHandler)
	founder.Get("/details", userHandler.GetUserDetailsHandler)
}
