package routes

import (
	"DBackend/internal/database"
	"DBackend/internal/server/handlers"
	"github.com/gofiber/fiber/v2"
  "DBackend/internal/server/middleware"
)

func MatchRoutes(api fiber.Router, db database.Service) {
	matchGroup := api.Group("/match", middleware.JWTMiddleware(db))
	// This GET route will call the MatchHandler
	MatchHandler := handlers.NewMatHandler(db)

	matchGroup.Get("/data/:userID", MatchHandler.MatchHandler)
}
