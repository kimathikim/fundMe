package routes

import (
	"DBackend/internal/database"
	"DBackend/internal/server/handlers"
	"github.com/gofiber/fiber/v2"
)

// AuthRoutes registers authentication routes
func AuthRoutes(api fiber.Router, db database.Service) {
	authHandler := handlers.NewAuthHandler(db)
	api.Post("/login", authHandler.LoginHandler)
	api.Post("/logout", authHandler.LogoutHandler)
}

