package handlers

import (
	"fmt"
	"time"

	"DBackend/internal/database"
	"DBackend/internal/server/services"
	"DBackend/model"
	"DBackend/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// InvestorHandler struct to handle investor-related operations
type InvestorHandler struct {
	db database.Service
}

// NewInvestorHandler initializes an InvestorHandler
func NewInvestorHandler(db database.Service) *InvestorHandler {
	return &InvestorHandler{db: db}
}

// UpdateInvestorHandler handles updating investor profile
func (h *InvestorHandler) UpdateInvestorHandler(c *fiber.Ctx) error {
	// Parse JSON data from form-data
	data := new(struct {
		//		InvestmentPortfolio   []string `json:"investment_portfolio"`
		TotalInvested         float64  `json:"total_invested"`
		InvestorType          string   `json:"investor_type"`
		Thesis                string   `json:"thesis"`
		PreferredFundingStage string   `json:"preferred_funding_stage"`
		InvestmentRange       string   `json:"investment_range"`
		InvestmentFrequency   string   `json:"investment_frequency"`
		RiskTolerance         string   `json:"risk_tolerance"`
		ExitStrategy          string   `json:"exit_strategy"`
		PreferredIndustries   []string `json:"preferred_industries"`
		PreferredRegions      []string `json:"preferred_regions"`
	})

	if err := c.BodyParser(data); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}
	fmt.Println("data", data)

	// Get the token from the context
	userToken := c.Locals("token")
	tokenStr, ok := userToken.(string)
	if !ok || tokenStr == "" {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid token"})
	}

	// Validate token and extract claims
	claims, err := utils.ValidateJWT(tokenStr)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid token"})
	}

	// Convert claims.UserID to primitive.ObjectID
	id, err := primitive.ObjectIDFromHex(claims.UserID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	// Find the investor in the database
	userInterface, err := h.db.User().FindByID(c.Context(), "investors", id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Investor not found"})
	}

	// Assert correct type for the investor object
	investor, ok := userInterface.(*model.Investor)
	if !ok {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to parse investor data"})
	}

	// Update the investor's profile with provided data
	investor.TotalInvested = data.TotalInvested
	//	investor.InvestmentPortfolio = convertToObjectIDArray(data.InvestmentPortfolio)
	investor.InvestorType = data.InvestorType
	investor.Thesis = data.Thesis
	investor.PreferredFundingStage = data.PreferredFundingStage
	investor.InvestmentRange = data.InvestmentRange
	investor.InvestmentFrequency = data.InvestmentFrequency
	investor.RiskTolerance = data.RiskTolerance
	investor.ExitStrategy = data.ExitStrategy
	investor.PreferredIndustries = data.PreferredIndustries
	investor.PreferredRegions = data.PreferredRegions

	// Save the updated data to MongoDB
	if _, err := h.db.User().UpdateInvestor(c.Context(), id, *investor); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update investor profile"})
	}

	return c.JSON(fiber.Map{"message": "Investor profile updated successfully"})
}

// GetInvestorDetailsHandler handles retrieving investor profile details
func (h *InvestorHandler) GetInvestorDetailsHandler(c *fiber.Ctx) error {
	// Get the token from the context
	tokenStr, ok := c.Locals("token").(string)
	if !ok || tokenStr == "" {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid token"})
	}

	// Validate token and extract claims
	claims, err := utils.ValidateJWT(tokenStr)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid token"})
	}

	// Convert claims.UserID to primitive.ObjectID
	id, err := primitive.ObjectIDFromHex(claims.UserID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	// Find the investor in the database
	userInterface, err := h.db.User().FindByID(c.Context(), "investors", id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Investor not found"})
	}

	// Assert correct type for the investor object
	investor, ok := userInterface.(*model.Investor)
	if !ok {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to parse investor data"})
	}

	// Return the investor's profile in the same format as the frontend
	return c.JSON(fiber.Map{
		//		"investment_portfolio":    investor.InvestmentPortfolio,
		"total_invested":          investor.TotalInvested,
		"investor_type":           investor.InvestorType,
		"thesis":                  investor.Thesis,
		"preferred_funding_stage": investor.PreferredFundingStage,
		"investment_range":        investor.InvestmentRange,
		"investment_frequency":    investor.InvestmentFrequency,
		"risk_tolerance":          investor.RiskTolerance,
		"exit_strategy":           investor.ExitStrategy,
		"preferred_industries":    investor.PreferredIndustries,
		"preferred_regions":       investor.PreferredRegions,
	})
}

// GetStartupDetailsHandler handles retrieving startup details from the founders
func (h *InvestorHandler) GetStartupDetailsHandler(c *fiber.Ctx) error {
	// Get all the founder details and return a JSON response
	founders, err := h.db.User().GetStartupDetails(c.Context())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to retrieve startup details"})
	}

	// Return the founders' details
	return c.JSON(fiber.Map{"founders": founders})
}

func (h *InvestorHandler) GetFounderProfilesHandler(c *fiber.Ctx) error {
	// Get all founder profiles
	founders, err := h.db.User().GetStartupDetails(c.Context())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to retrieve founder profiles"})
	}

	var founderProfiles []bson.M
	for _, founder := range founders {
		profile, err := h.db.User().GetFounderProfileWithMatch(c.Context(), founder.UserID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to retrieve founder profile"})
		}
		founderProfiles = append(founderProfiles, profile)
	}

	return c.JSON(fiber.Map{"founder_profiles": founderProfiles})
}

// convertToObjectIDArray converts an array of string IDs to ObjectIDs
func convertToObjectIDArray(ids []string) []primitive.ObjectID {
	var objectIDs []primitive.ObjectID
	for _, id := range ids {
		objectID, err := primitive.ObjectIDFromHex(id)
		if err == nil {
			objectIDs = append(objectIDs, objectID)
		}
	}
	return objectIDs
}

// GetMeetingsHandler - Retrieve all meetings for an investor
func (h *InvestorHandler) GetMeetingsHandler(c *fiber.Ctx) error {
	id, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid investor ID"})
	}

	meetings, err := h.db.User().GetMeetings(c.Context(), id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to retrieve meetings"})
	}

	return c.JSON(fiber.Map{"meetings": meetings})
}

// AddMeetingHandler - Add a meeting to an investor's calendar
func (h *InvestorHandler) AddMeetingHandler(c *fiber.Ctx) error {
	id, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid investor ID"})
	}

	var meeting model.Meeting
	if err := c.BodyParser(&meeting); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	fmt.Println("meeting", meeting)

	meeting.ID = primitive.NewObjectID()

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

	updateResult, err := h.db.DealFlow().AddMeeting(c.Context(), id, meeting)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to add meeting"})
	}

	return c.JSON(fiber.Map{"message": "Meeting added successfully", "modifiedCount": updateResult.ModifiedCount})
}

// GetInvestorDashboardHandler returns dashboard summary for investors
func (h *InvestorHandler) GetInvestorDashboardHandler(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}

	// Get investor data from database
	investorID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	// Get portfolio summary
	portfolioSummary, err := h.db.Investor().GetPortfolioSummary(c.Context(), investorID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to get portfolio summary"})
	}

	// Get pipeline summary
	pipelineSummary, err := h.db.Investor().GetPipelineSummary(c.Context(), investorID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to get pipeline summary"})
	}

	// Get recent activities
	recentActivities, err := h.db.Investor().GetRecentActivities(c.Context(), investorID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to get recent activities"})
	}

	return c.JSON(fiber.Map{
		"portfolioSummary": portfolioSummary,
		"pipelineSummary":  pipelineSummary,
		"recentActivities": recentActivities,
	})
}

// GetPortfolioPerformanceHandler returns portfolio performance data
func (h *InvestorHandler) GetPortfolioPerformanceHandler(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}

	// Get period from query params
	period := c.Query("period", "all")

	// Validate period
	validPeriods := map[string]bool{"1m": true, "3m": true, "6m": true, "1y": true, "all": true}
	if !validPeriods[period] {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid period"})
	}

	// Get investor data from database
	investorID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	// Get performance data
	performanceData, err := h.db.Investor().GetPerformanceData(c.Context(), investorID, period)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to get performance data"})
	}

	// Get metrics
	metrics, err := h.db.Investor().GetPerformanceMetrics(c.Context(), investorID, period)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to get performance metrics"})
	}

	return c.JSON(fiber.Map{
		"performanceData": performanceData,
		"metrics":         metrics,
	})
}

// GetAllNotificationsHandler retrieves all notifications for the authenticated investor
func (h *InvestorHandler) GetAllNotificationsHandler(c *fiber.Ctx) error {
	// Get the user ID from the context (set by JWT middleware)
	userID, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}

	// Convert string ID to ObjectID
	investorID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	// Get notifications for this investor
	notifications, err := h.db.User().GetAllNotificationsByFounder(c.Context(), investorID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to retrieve notifications"})
	}

	return c.JSON(fiber.Map{"notifications": notifications})
}

// UpdateNotificationHandler updates a specific notification for the authenticated investor
func (h *InvestorHandler) UpdateNotificationHandler(c *fiber.Ctx) error {
	notificationID, err := primitive.ObjectIDFromHex(c.Params("notificationID"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid notification ID"})
	}

	var updateData struct {
		ReadStatus bool `json:"read_status"`
	}
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Update the notification
	updateFields := bson.M{"read_status": updateData.ReadStatus}
	err = h.db.User().UpdateNotification(c.Context(), notificationID, updateFields)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update notification"})
	}

	return c.JSON(fiber.Map{"message": "Notification updated successfully"})
}

// DeleteNotificationHandler deletes a specific notification for the authenticated investor
func (h *InvestorHandler) DeleteNotificationHandler(c *fiber.Ctx) error {
	notificationID, err := primitive.ObjectIDFromHex(c.Params("notificationID"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid notification ID"})
	}

	// Delete the notification
	_, err = h.db.User().DeleteNotification(c.Context(), notificationID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete notification"})
	}

	return c.JSON(fiber.Map{"message": "Notification deleted successfully"})
}
