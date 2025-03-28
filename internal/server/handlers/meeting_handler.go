package handlers

import (
    "time"

    "DBackend/internal/database"
    "DBackend/internal/server/services"
    "DBackend/model"
    "DBackend/utils"

    "github.com/gofiber/fiber/v2"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

// ScheduleMeeting handles creating a new meeting
func ScheduleMeeting(db database.Service) fiber.Handler {
    return func(c *fiber.Ctx) error {
        var meeting model.Meeting
        if err := c.BodyParser(&meeting); err != nil {
            return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
        }

        // Generate new ID for the meeting
        meeting.ID = primitive.NewObjectID()

        // Get user ID from token
        tokenStr, ok := c.Locals("token").(string)
        if !ok || tokenStr == "" {
            return c.Status(401).JSON(fiber.Map{"error": "Invalid token"})
        }

        claims, err := utils.ValidateJWT(tokenStr)
        if err != nil {
            return c.Status(401).JSON(fiber.Map{"error": "Invalid token"})
        }

        // Create Google Calendar event
        calendarService, err := services.NewGoogleCalendarService()
        if err != nil {
            return c.Status(500).JSON(fiber.Map{"error": "Failed to initialize Google Calendar service"})
        }

        event, err := calendarService.CreateEvent(
            meeting.Title,
            "", // Location
            meeting.Notes,
            meeting.StartTime.Format(time.RFC3339),
            meeting.EndTime.Format(time.RFC3339),
            "UTC",
            []string{}, // Attendees
        )
        if err != nil {
            return c.Status(500).JSON(fiber.Map{"error": "Failed to create Google Calendar event"})
        }

        meeting.GoogleMeetURL = event.HangoutLink
        userID, _ := primitive.ObjectIDFromHex(claims.UserID)
        
        // Assuming the current user is the investor
        meeting.InvestorID = userID
        
        // If FounderID is not set in the request, you might need to handle that separately
        if meeting.FounderID.IsZero() {
            return c.Status(400).JSON(fiber.Map{"error": "Founder ID is required"})
        }
        meeting.CreatedAt = time.Now()
        meeting.UpdatedAt = time.Now()

        // Save meeting to database
        result, err := db.DealFlow().AddMeeting(c.Context(), primitive.NilObjectID, meeting)
        if err != nil {
            return c.Status(500).JSON(fiber.Map{"error": "Failed to create meeting"})
        }

        return c.Status(201).JSON(fiber.Map{
            "message": "Meeting scheduled successfully",
            "meeting": result,
        })
    }
}

// GetAllMeetings returns all meetings
func GetAllMeetings(db database.Service) fiber.Handler {
    return func(c *fiber.Ctx) error {
        meetings, err := db.User().GetMeetings(c.Context(), primitive.NilObjectID)
        if err != nil {
            return c.Status(500).JSON(fiber.Map{"error": "Failed to retrieve meetings"})
        }
        return c.JSON(fiber.Map{"meetings": meetings})
    }
}

// GetMeetingByID returns a specific meeting by ID
func GetMeetingByID(db database.Service) fiber.Handler {
    return func(c *fiber.Ctx) error {
        id, err := primitive.ObjectIDFromHex(c.Params("id"))
        if err != nil {
            return c.Status(400).JSON(fiber.Map{"error": "Invalid meeting ID"})
        }

        var meeting model.Meeting
        result, err := db.User().FindByID(c.Context(), "meetings", id)
        if err != nil {
            return c.Status(404).JSON(fiber.Map{"error": "Meeting not found"})
        }
        
        meeting, ok := result.(model.Meeting)
        if !ok {
            return c.Status(500).JSON(fiber.Map{"error": "Failed to parse meeting data"})
        }

        return c.JSON(fiber.Map{"meeting": meeting})
    }
}

// UpdateMeeting updates an existing meeting
func UpdateMeeting(db database.Service) fiber.Handler {
    return func(c *fiber.Ctx) error {
        id, err := primitive.ObjectIDFromHex(c.Params("id"))
        if err != nil {
            return c.Status(400).JSON(fiber.Map{"error": "Invalid meeting ID"})
        }

        var updates model.Meeting
        if err := c.BodyParser(&updates); err != nil {
            return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
        }

        updates.UpdatedAt = time.Now()

        err = db.User().UpdateMeeting(c.Context(), id, updates)
        if err != nil {
            return c.Status(500).JSON(fiber.Map{"error": "Failed to update meeting"})
        }

        return c.JSON(fiber.Map{"message": "Meeting updated successfully"})
    }
}

// CancelMeeting deletes a meeting
func CancelMeeting(db database.Service) fiber.Handler {
    return func(c *fiber.Ctx) error {
        id, err := primitive.ObjectIDFromHex(c.Params("id"))
        if err != nil {
            return c.Status(400).JSON(fiber.Map{"error": "Invalid meeting ID"})
        }

        // Use an existing method that can delete a meeting
        _, err = db.User().DeleteNotification(c.Context(), id)
        if err != nil {
            return c.Status(500).JSON(fiber.Map{"error": "Failed to cancel meeting"})
        }

        return c.JSON(fiber.Map{"message": "Meeting cancelled successfully"})
    }
}

// GetUserMeetings returns all meetings for the current user
func GetUserMeetings(db database.Service) fiber.Handler {
    return func(c *fiber.Ctx) error {
        tokenStr, ok := c.Locals("token").(string)
        if !ok || tokenStr == "" {
            return c.Status(401).JSON(fiber.Map{"error": "Invalid token"})
        }

        claims, err := utils.ValidateJWT(tokenStr)
        if err != nil {
            return c.Status(401).JSON(fiber.Map{"error": "Invalid token"})
        }

        userID, err := primitive.ObjectIDFromHex(claims.UserID)
        if err != nil {
            return c.Status(400).JSON(fiber.Map{"error": "Invalid user ID"})
        }

        meetings, err := db.User().GetMeetings(c.Context(), userID)
        if err != nil {
            return c.Status(500).JSON(fiber.Map{"error": "Failed to retrieve meetings"})
        }

        return c.JSON(fiber.Map{"meetings": meetings})
    }
}

// AddMeetingNotes adds notes to a meeting
func AddMeetingNotes(db database.Service) fiber.Handler {
    return func(c *fiber.Ctx) error {
        id, err := primitive.ObjectIDFromHex(c.Params("id"))
        if err != nil {
            return c.Status(400).JSON(fiber.Map{"error": "Invalid meeting ID"})
        }

        var data struct {
            Notes string `json:"notes"`
        }
        if err := c.BodyParser(&data); err != nil {
            return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
        }

        err = db.User().AddMeetingNotes(c.Context(), id, data.Notes)
        if err != nil {
            return c.Status(500).JSON(fiber.Map{"error": "Failed to add meeting notes"})
        }

        return c.JSON(fiber.Map{"message": "Meeting notes added successfully"})
    }
}

// GetMeetingNotes retrieves notes for a meeting
func GetMeetingNotes(db database.Service) fiber.Handler {
    return func(c *fiber.Ctx) error {
        id, err := primitive.ObjectIDFromHex(c.Params("id"))
        if err != nil {
            return c.Status(400).JSON(fiber.Map{"error": "Invalid meeting ID"})
        }

        notes, err := db.User().GetMeetingNotes(c.Context(), id)
        if err != nil {
            return c.Status(500).JSON(fiber.Map{"error": "Failed to retrieve meeting notes"})
        }

        return c.JSON(fiber.Map{"notes": notes})
    }
}

// AddMeetingParticipant adds a participant to a meeting
func AddMeetingParticipant(db database.Service) fiber.Handler {
    return func(c *fiber.Ctx) error {
        id, err := primitive.ObjectIDFromHex(c.Params("id"))
        if err != nil {
            return c.Status(400).JSON(fiber.Map{"error": "Invalid meeting ID"})
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

        err = db.User().AddMeetingParticipant(c.Context(), id, userID)
        if err != nil {
            return c.Status(500).JSON(fiber.Map{"error": "Failed to add participant"})
        }

        return c.JSON(fiber.Map{"message": "Participant added successfully"})
    }
}

// RemoveMeetingParticipant removes a participant from a meeting
func RemoveMeetingParticipant(db database.Service) fiber.Handler {
    return func(c *fiber.Ctx) error {
        id, err := primitive.ObjectIDFromHex(c.Params("id"))
        if err != nil {
            return c.Status(400).JSON(fiber.Map{"error": "Invalid meeting ID"})
        }

        userID, err := primitive.ObjectIDFromHex(c.Params("userId"))
        if err != nil {
            return c.Status(400).JSON(fiber.Map{"error": "Invalid user ID"})
        }

        err = db.User().RemoveMeetingParticipant(c.Context(), id, userID)
        if err != nil {
            return c.Status(500).JSON(fiber.Map{"error": "Failed to remove participant"})
        }

        return c.JSON(fiber.Map{"message": "Participant removed successfully"})
    }
}
