package routes

import (
	"DBackend/internal/database"
	"DBackend/internal/server/handlers"
	"DBackend/internal/server/middleware"

	"github.com/gofiber/fiber/v2"
)

// AuthRoutes registers authentication routes
func AuthRoutes(api fiber.Router, db database.Service) {
	ao := api.Group("/auth")
	authJwT := api.Group("/get", middleware.JWTMiddleware(db))
	authHandler := handlers.NewAuthHandler(db)
	ao.Post("/login", authHandler.LoginHandler)
	ao.Post("/logout", authHandler.LogoutHandler)

	authJwT.Get("/me", authHandler.MeHandler)
}
