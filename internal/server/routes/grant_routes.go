package routes

import (
    "DBackend/internal/database"
    "DBackend/internal/server/handlers"
    "DBackend/internal/server/middleware"

    "github.com/gofiber/fiber/v2"
)

// GrantRoutes sets up all grant-related routes
func GrantRoutes(api fiber.Router, db database.Service) {
    grant := api.Group("/grants")
    
    // Public routes
    grant.Get("/", handlers.NewFounderHandler(db).GetGrantsHandler)
    grant.Get("/:id", handlers.GetGrantByID(db))
    
    // Protected routes
    grant.Use(middleware.JWTMiddleware(db))
    grant.Post("/", handlers.CreateGrant(db))
    grant.Put("/:id", handlers.UpdateGrant(db))
    grant.Delete("/:id", handlers.DeleteGrant(db))
    
    // Grant application routes
    grant.Post("/apply", handlers.ApplyForGrant(db))
    grant.Get("/applications", handlers.GetGrantApplications(db))
    grant.Get("/applications/:id", handlers.GetGrantApplicationByID(db))
    grant.Put("/applications/:id", handlers.UpdateGrantApplication(db))
}
