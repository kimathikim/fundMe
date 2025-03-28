package handlers

import (
	"time"

	"DBackend/internal/database"
	"DBackend/model"
	"DBackend/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateTask handles creating a new task
func CreateTask(db database.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var task model.Task
		if err := c.BodyParser(&task); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
		}

		// Generate new ID for the task
		task.ID = primitive.NewObjectID()
		task.CreatedAt = time.Now()
		task.UpdatedAt = time.Now()

		// Get user ID from token
		tokenStr, ok := c.Locals("token").(string)
		if !ok || tokenStr == "" {
			return c.Status(401).JSON(fiber.Map{"error": "Invalid token"})
		}

		claims, err := utils.ValidateJWT(tokenStr)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "Invalid token"})
		}

		userID, _ := primitive.ObjectIDFromHex(claims.UserID)
		task.CreatedBy = userID

		// Save task to database
		// Since this is a standalone task, we'll use primitive.NilObjectID for dealflow
		result, err := db.User().AddTask(c.Context(), primitive.NilObjectID, task)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to create task"})
		}

		return c.Status(201).JSON(fiber.Map{
			"message": "Task created successfully",
			"task":    result,
		})
	}
}

// GetAllTasks returns all tasks
func GetAllTasks(db database.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tasks, err := db.User().GetAllTasks(c.Context())
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to retrieve tasks"})
		}
		return c.JSON(fiber.Map{"tasks": tasks})
	}
}

// GetTaskByID returns a specific task by ID
func GetTaskByID(db database.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := primitive.ObjectIDFromHex(c.Params("id"))
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid task ID"})
		}

		task, err := db.User().GetTaskByID(c.Context(), id)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "Task not found"})
		}

		return c.JSON(fiber.Map{"task": task})
	}
}

// UpdateTask updates an existing task
func UpdateTask(db database.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := primitive.ObjectIDFromHex(c.Params("id"))
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid task ID"})
		}

		var updates model.Task
		if err := c.BodyParser(&updates); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
		}

		updates.UpdatedAt = time.Now()

		err = db.User().UpdateTask(c.Context(), id, updates)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to update task"})
		}

		return c.JSON(fiber.Map{"message": "Task updated successfully"})
	}
}

// DeleteTask deletes a task
func DeleteTask(db database.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := primitive.ObjectIDFromHex(c.Params("id"))
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid task ID"})
		}

		err = db.User().DeleteTask(c.Context(), id)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to delete task"})
		}

		return c.JSON(fiber.Map{"message": "Task deleted successfully"})
	}
}

// GetTasksByUser returns all tasks for a specific user
func GetTasksByUser(db database.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, err := primitive.ObjectIDFromHex(c.Params("userId"))
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid user ID"})
		}

		tasks, err := db.User().GetTasksByUser(c.Context(), userID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to retrieve tasks"})
		}

		return c.JSON(fiber.Map{"tasks": tasks})
	}
}

// UpdateTaskStatus updates the status of a task
func UpdateTaskStatus(db database.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := primitive.ObjectIDFromHex(c.Params("id"))
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid task ID"})
		}

		var data struct {
			Completed bool `json:"completed"`
		}
		if err := c.BodyParser(&data); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
		}

		// First try to get the task to check if it exists
		_, err = db.User().GetTaskByID(c.Context(), id)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "Task not found"})
		}

		// Update task completion status
		_, err = db.User().UpdateTaskStatus(c.Context(), primitive.NilObjectID, id, data.Completed)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to update task status"})
		}

		return c.JSON(fiber.Map{
			"message": "Task status updated successfully",
			"task": fiber.Map{
				"id":        id.Hex(),
				"completed": data.Completed,
			},
		})
	}
}

// AssignTask assigns a task to a user
func AssignTask(db database.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		taskID, err := primitive.ObjectIDFromHex(c.Params("id"))
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid task ID"})
		}

		var data struct {
			UserID string `json:"user_id"`
		}
		if err := c.BodyParser(&data); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
		}

		userID, err := primitive.ObjectIDFromHex(data.UserID)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid user ID"})
		}

		err = db.User().AssignTask(c.Context(), taskID, userID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to assign task"})
		}

		return c.JSON(fiber.Map{"message": "Task assigned successfully"})
	}
}
