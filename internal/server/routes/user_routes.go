package routes

import (
	"DBackend/internal/database"
	"DBackend/internal/server/handlers"
	"github.com/gofiber/fiber/v2"
)

// UserRoutes registers user-related routes
func UserRoutes(api fiber.Router, db database.Service) {
	apiV1 := api.Group("/user")
	userHandler := handlers.NewUserHandler(db)
	apiV1.Post("/register", userHandler.RegisterHandler)
}
