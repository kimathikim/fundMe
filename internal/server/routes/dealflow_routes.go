package routes

import (
	"DBackend/internal/database"
	"DBackend/internal/server/handlers"
	"DBackend/internal/server/middleware"
	"github.com/gofiber/fiber/v2"
)

func DealFlowRoutes(api fiber.Router, db database.Service) {
	handler := handlers.NewDealFlowHandler(db)
	dealflow := api.Group("/dealflow", middleware.JWTMiddleware(db))

	dealflow.Post("/", handler.AddDealFlowHandler)
	dealflow.Get("/:id", handler.GetDealFlowByIDHandler)
	dealflow.Get("/", handler.ListAllDealFlowHandler)
	dealflow.Put("/:id", handler.UpdateDealFlowHandler)
	dealflow.Delete("/:id", handler.DeleteDealFlowHandler)
	dealflow.Post("/:id/invest", middleware.RequireRole("investor"), handler.InvestInStartupHandler)
	dealflow.Post("/:id/meetings", handler.AddMeetingHandler)
	dealflow.Post("/:id/documents", handler.AddDocumentHandler)
	dealflow.Post("/:id/tasks", handler.AddTaskHandler)
	dealflow.Patch("/:id/tasks/:taskID", handler.UpdateTaskStatusHandler)
	dealflow.Patch("/:id/stage", handler.UpdateDealStageHandler)
	dealflow.Patch("/:id/status", handler.UpdateDealStatusHandler)

	// Remove duplicate route
	// dealflow.Post("/:id/meeting", middleware.RequireRole("investor"), handler.AddMeetingHandler)
}
