package routes

import (
    "DBackend/internal/database"
    "DBackend/internal/server/handlers"
    "DBackend/internal/server/middleware"

    "github.com/gofiber/fiber/v2"
)

// MeetingRoutes sets up all meeting-related routes
func MeetingRoutes(api fiber.Router, db database.Service) {
    meeting := api.Group("/meetings")
    
    // All meeting routes are protected
    meeting.Use(middleware.JWTMiddleware(db))
    
    // Meeting CRUD operations
    meeting.Post("/", handlers.ScheduleMeeting(db))
    meeting.Get("/", handlers.GetAllMeetings(db))
    meeting.Get("/:id", handlers.GetMeetingByID(db))
    meeting.Put("/:id", handlers.UpdateMeeting(db))
    meeting.Delete("/:id", handlers.CancelMeeting(db))
    
    // User-specific meetings
    meeting.Get("/user", handlers.GetUserMeetings(db))
    
    // Meeting notes
    meeting.Post("/:id/notes", handlers.AddMeetingNotes(db))
    meeting.Get("/:id/notes", handlers.GetMeetingNotes(db))
    
    // Meeting participants
    meeting.Post("/:id/participants", handlers.AddMeetingParticipant(db))
    meeting.Delete("/:id/participants/:userId", handlers.RemoveMeetingParticipant(db))
}