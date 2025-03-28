package routes

import (
	"DBackend/internal/database"
	"DBackend/internal/server/handlers"
	"DBackend/internal/server/middleware"
	"github.com/gofiber/fiber/v2"
)

// InvestorRoutes registers investor routes (role-based)
func InvestorRoutes(api fiber.Router, db database.Service) {
	investor := api.Group("/investor", middleware.JWTMiddleware(db))
	investorHandler := handlers.NewInvestorHandler(db)
	userHandler := handlers.NewUserHandler(db)

	// Profile routes
	investor.Patch("/profile", middleware.RequireRole("investor"), investorHandler.UpdateInvestorHandler)
	investor.Get("/profile", middleware.RequireRole("investor"), investorHandler.GetInvestorDetailsHandler)
	investor.Get("/", userHandler.GetUserDetailsHandler)

	// Startup and founder routes
	investor.Get("/startups", investorHandler.GetStartupDetailsHandler)
	investor.Get("/founderProfile", investorHandler.GetFounderProfilesHandler)

	// Meeting routes
	investor.Post("/investor/:id/meeting", investorHandler.AddMeetingHandler)
	investor.Get("/investor/:id/meetings", investorHandler.GetMeetingsHandler)

	// Notification routes
	investor.Get("/notifications", middleware.RequireRole("investor"), investorHandler.GetAllNotificationsHandler)
	investor.Put("/notifications/:notificationID", middleware.RequireRole("investor"), investorHandler.UpdateNotificationHandler)
	investor.Delete("/notification/:notificationID", middleware.RequireRole("investor"), investorHandler.DeleteNotificationHandler)

	// Dashboard routes
	investor.Get("/dashboard", middleware.RequireRole("investor"), investorHandler.GetInvestorDashboardHandler)
	investor.Get("/portfolio/performance", middleware.RequireRole("investor"), investorHandler.GetPortfolioPerformanceHandler)
}
