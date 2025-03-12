package routes

import (
	"DBackend/internal/database"
	"DBackend/internal/server/handlers"
	"DBackend/internal/server/middleware"
	"github.com/gofiber/fiber/v2"
)

// fousnderRoutes registers founder routes (role-based)
func InvestorRoutes(api fiber.Router, db database.Service) {
	investor := api.Group("/investor", middleware.JWTMiddleware(db))
	investorHandler := handlers.NewInvestorHandler(db)
	userHandler := handlers.NewUserHandler(db)
	investor.Patch("/profile", middleware.RequireRole("investor"), investorHandler.UpdateInvestorHandler)
	investor.Get("/profile", middleware.RequireRole("investor"), investorHandler.GetInvestorDetailsHandler)
	investor.Get("/details", userHandler.GetUserDetailsHandler)
}
