package server

import (
	"context"
	"fmt"
	"log"
	"time"

	"DBackend/internal/database"
	"DBackend/internal/server/routes"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func (s *FiberServer) RegisterFiberRoutes() {
	// Apply CORS middleware
	api := s.Group("/api/v1")
	api.Get("/health", s.healthHandler)
	api.Get("/websocket", websocket.New(s.websocketHandler))
}

func (s *FiberServer) healthHandler(c *fiber.Ctx) error {
	return c.JSON(s.db.Health())
}

func SetupRoutes(app *fiber.App, db database.Service, prefix string) {
	api := app.Group("/" + prefix)
	routes.UserRoutes(api, db)
	routes.AuthRoutes(api, db)
  routes.ProtectedRoutes(api, db)
//	NotFoundRoute(app)
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
