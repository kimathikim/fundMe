package server

import (
	"DBackend/internal/database"
	"DBackend/internal/server/middleware"
	"github.com/gofiber/fiber/v2"
)

// FiberServer defines the server structure
type FiberServer struct {
	*fiber.App
	db database.Service
//   dbc database.DealFlowService
}

// New creates a new server instance
func New() *FiberServer {
	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: "DBackend",
			AppName:      "DBackend",
		}),
		db: database.New(),
	}

	if server.db != nil {
		println("DB is not nil")
	}

	// Apply CORS middleware globally
	server.Use(middleware.CORSMiddleware())

	// Register routes
	SetupRoutes(server.App, server.db, "api/v1")
  

	return server
}
