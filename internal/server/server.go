package server

import (
	"github.com/gofiber/fiber/v2"

	"DBackend/internal/database"
)

type FiberServer struct {
	*fiber.App

	db database.Service
}

func New() *FiberServer {
	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: "DBackend",
			AppName:      "DBackend",
		}),

		db: database.New(),
	}

	return server
}
