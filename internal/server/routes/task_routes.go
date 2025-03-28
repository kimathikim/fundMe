package routes

import (
	"DBackend/internal/database"
	"DBackend/internal/server/handlers"
	"DBackend/internal/server/middleware"

	"github.com/gofiber/fiber/v2"
)

// TaskRoutes sets up all task-related routes
func TaskRoutes(api fiber.Router, db database.Service) {
	task := api.Group("/tasks")

	// All task routes are protected
	task.Use(middleware.JWTMiddleware(db))

	// Task CRUD operations
	task.Post("/", handlers.CreateTask(db))
	task.Get("/", handlers.GetAllTasks(db))
	task.Get("/:id", handlers.GetTaskByID(db))
	task.Put("/:id", handlers.UpdateTask(db))
	task.Delete("/:id", handlers.DeleteTask(db))

	// User-specific tasks
	task.Get("/user/:userId", handlers.GetTasksByUser(db))

	// Task status updates
	task.Patch("/:id/status", handlers.UpdateTaskStatus(db))

	// Task assignment
	task.Post("/:id/assign", handlers.AssignTask(db))
}
