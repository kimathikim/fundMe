package server

import (
	"context"
	"fmt"
	"log"
	"time"

	"DBackend/internal/database"
	"DBackend/internal/server/middleware"
	"DBackend/internal/server/routes"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func (s *FiberServer) RegisterFiberRoutes() {
	api := s.Group("/api/v1", middleware.CORSMiddleware())
	api.Get("/health", s.healthHandler)
	api.Get("/websocket", websocket.New(s.websocketHandler))
	
	// Register all other routes
	routes.UserRoutes(api, s.db)
	routes.AuthRoutes(api, s.db)
	routes.FounderRoutes(api, s.db)
	routes.InvestorRoutes(api, s.db)
	routes.MatchRoutes(api, s.db)
	routes.DealFlowRoutes(api, s.db)
	
	// Register new routes
	routes.GrantRoutes(api, s.db)
	routes.TaskRoutes(api, s.db)
	routes.MeetingRoutes(api, s.db)
}

func (s *FiberServer) healthHandler(c *fiber.Ctx) error {
	return c.JSON(s.db.Health())
}

func SetupRoutes(app *fiber.App, db database.Service, prefix string) {
	api := app.Group("/" + prefix)
	routes.UserRoutes(api, db)
	routes.AuthRoutes(api, db)
	routes.FounderRoutes(api, db)
	routes.InvestorRoutes(api, db)
	routes.MatchRoutes(api, db)
	routes.DealFlowRoutes(api, db)
	
	// Register new routes
	routes.GrantRoutes(api, db)
	routes.TaskRoutes(api, db)
	routes.MeetingRoutes(api, db)
	
	NotFoundRoute(app)
}

func NotFoundRoute(app *fiber.App) {
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(404).JSON(fiber.Map{"error": "Route not found"})
	})
}

func (s *FiberServer) websocketHandler(con *websocket.Conn) {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		for {
			_, _, err := con.ReadMessage()
			if err != nil {
				cancel()
				log.Println("Receiver Closing", err)
				break
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			payload := fmt.Sprintf("server timestamp: %d", time.Now().UnixNano())
			if err := con.WriteMessage(websocket.TextMessage, []byte(payload)); err != nil {
				log.Printf("could not write to socket: %v", err)
				return
			}
			time.Sleep(time.Second * 2)
		}
	}
}
