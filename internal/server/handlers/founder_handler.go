package handlers

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"DBackend/internal/database"
	"DBackend/model"
	"DBackend/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/bson"
)

type FounderHandler struct {
	db database.Service
}

func NewFounderHandler(db database.Service) *FounderHandler {
	return &FounderHandler{db: db}
}

func (h *FounderHandler) UpdateFounderHandler(c *fiber.Ctx) error {
	// Parse JSON data from form-data
	data := new(struct {
		StartupName       string `json:"startup_name"`
		MissionStatement  string `json:"mission_statement"`
		Industry          string `json:"industry"`
		FundingStage      string `json:"funding_stage"`
		FundingAllocation string `json:"funding_allocation"`
		BussinessModel    string `json:"bussiness_model"`
		RevenueStreams    string `json:"revenue_streams"`
		Traction          string `json:"traction"`
		ScalingPotential  string `json:"scaling_potential"`
		TotalInvested     int    `json:"total_invested"`
		FundRequired      string `json:"fund_required"`
		Competition       string `json:"competition"`
		LeadershipTeam    string `json:"leadership_team"`
		TeamSize          string `json:"team_size"`
		Location          string `json:"location"`
		StartupWebsite    string `json:"startup_website"`
	})

	// Extract JSON data from form field "data"
	jsonData := c.FormValue("data")
	if jsonData == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Missing JSON data"})
	}

	// Unmarshal JSON into struct
	if err := json.Unmarshal([]byte(jsonData), data); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid JSON format"})
	}

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

	// Find the user based on the ID in the claims
	userInterface, err := h.db.User().FindByID(c.Context(), "founders", id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	// Assert correct type for the user object
	user, ok := userInterface.(*model.Founder)
	if !ok {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to parse user data"})
	}

	// Convert FundRequired from string to int
	fundRequired, err := strconv.Atoi(data.FundRequired)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid fund required value"})
	}

	// Handle pitch deck file upload
	file, err := c.FormFile("pitch_deck")
	if err == nil { // File is provided
		// Define storage path
		folderPath := "./Pitch"
		filePath := fmt.Sprintf("%s/%s_%s", folderPath, id.Hex(), file.Filename)

		// Ensure the folder exists
		if err := os.MkdirAll(folderPath, os.ModePerm); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to create storage directory"})
		}

		// Save the file
		if err := c.SaveFile(file, filePath); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to save pitch deck file"})
		}

		// Update the pitch deck file name in the user struct
		user.PitchDeck = filePath
	}

	// Add validation for required fields
	if data.StartupName == "" || data.Industry == "" || data.FundingStage == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Missing required fields"})
	}

	// Ensure FundRequired is never nil by providing a default
	if data.FundRequired == "" {
		data.FundRequired = "0"
	}

	// Ensure TotalInvested is never nil
	if data.TotalInvested < 0 {
		data.TotalInvested = 0
	}

	// Update user's profile with provided data
	user.StartupName = data.StartupName
	user.TotalInvested = data.TotalInvested
	user.FundRequired = fundRequired // Use the converted value
	user.MissionStatement = data.MissionStatement
	user.Industry = data.Industry
	user.FundingStage = data.FundingStage
	user.FundingAllocation = data.FundingAllocation
	user.BussinessModel = data.BussinessModel
	user.RevenueStreams = data.RevenueStreams
	user.Traction = data.Traction
	user.ScalingPotential = data.ScalingPotential
	user.Competition = data.Competition
	user.LeadershipTeam = data.LeadershipTeam
	user.TeamSize = data.TeamSize
	user.Location = data.Location
	user.StartupWebsite = data.StartupWebsite

	// Save the updated data to MongoDB
	if _, err := h.db.User().UpdateFounder(c.Context(), id, *user); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update profile"})
	}
	datad := fiber.Map{"message": "Founder profile updated successfully", "pitch_deck": *user}
	return c.JSON(datad)
}

func (h *FounderHandler) GetFounderDetailsHandler(c *fiber.Ctx) error {
	// Get the token from the context
	userToken := c.Locals("token")
	tokenStr, ok := userToken.(string)
	if !ok || tokenStr == "" {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid token"})
	}

	// Validate token and extract claims
	claims, err := utils.ValidateJWT(tokenStr)
	if err != nil {
		return c.Status(402).JSON(fiber.Map{"error": "Invalid token"})
	}

	// Convert claims.UserID to primitive.ObjectID
	id, err := primitive.ObjectIDFromHex(claims.UserID)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	// Find the founder in the database
	userInterface, err := h.db.User().FindByID(c.Context(), "founders", id)
	if err != nil {
		return c.Status(405).JSON(fiber.Map{"error": "Founder not found"})
	}

	// Assert correct type for the founder object
	founder, ok := userInterface.(*model.Founder)
	if !ok {
		return c.Status(501).JSON(fiber.Map{"error": "Failed to parse founder data"})
	}

	return c.JSON(fiber.Map{
		"id":                 id,
		"startup_name":       founder.StartupName,
		"mission_statement":  founder.MissionStatement,
		"industry":           founder.Industry,
		"funding_stage":      founder.FundingStage,
		"funding_allocation": founder.FundingAllocation,
		"bussiness_model":    founder.BussinessModel,
		"revenue_streams":    founder.RevenueStreams,
		"traction":           founder.Traction,
		"scaling_potential":  founder.ScalingPotential,
		"fund_required":      founder.FundRequired,
		"competition":        founder.Competition,
		"leadership_team":    founder.LeadershipTeam,
		"team_size":          founder.TeamSize,
		"location":           founder.Location,
		"startup_website":    founder.StartupWebsite,
		"pith_deck":          founder.PitchDeck,
	})
}

// GetAllNotificationsHandler retrieves all notifications for the authenticated founder.
func (h *FounderHandler) GetAllNotificationsHandler(c *fiber.Ctx) error {
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
	founderID, err := primitive.ObjectIDFromHex(claims.UserID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid user ID"})
	}
	// get the founder where user_id = founderID
	founder, err := h.db.User().GetFounderByUserID(c.Context(), founderID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to retrieve founder"})
	}

	founderID = founder.ID

	notifications, err := h.db.User().GetAllNotificationsByFounder(c.Context(), founderID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to retrieve notifications"})
	}

	return c.JSON(notifications)
}

// UpdateNotificationHandler updates a specific notification for the authenticated founder.
func (h *FounderHandler) UpdateNotificationHandler(c *fiber.Ctx) error {
	notificationID, err := primitive.ObjectIDFromHex(c.Params("notificationID"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid notification ID"})
	}

	var updateData model.Notification
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	update := bson.M{
		"title":       updateData.Title,
		"message":     updateData.Message,
		"updated_at":  time.Now(),
		"read_status": updateData.ReadStatus,
	}

	err = h.db.User().UpdateNotification(c.Context(), notificationID, update)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update notification"})
	}

	return c.JSON(fiber.Map{"message": "Notification updated successfully"})
}

// DeleteNotificationHandler deletes a specific notification for the authenticated founder.
func (h *FounderHandler) DeleteNotificationHandler(c *fiber.Ctx) error {
	notificationID, err := primitive.ObjectIDFromHex(c.Params("notificationID"))
	fmt.Println(notificationID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid notification ID"})
	}

	_, err = h.db.User().DeleteNotification(c.Context(), notificationID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete notification"})
	}

	return c.JSON(fiber.Map{"message": "Notification deleted successfully"})
}

// GetFounderDashboardHandler returns dashboard summary for founders
func (h *FounderHandler) GetFounderDashboardHandler(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}

	// Get founder data from database
	founderID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	// Get fundraising summary
	fundraisingSummary, err := h.db.Founder().GetFundraisingSummary(c.Context(), founderID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to get fundraising summary"})
	}

	// Get investor engagement
	investorEngagement, err := h.db.Founder().GetInvestorEngagement(c.Context(), founderID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to get investor engagement"})
	}

	return c.JSON(fiber.Map{
		"fundraisingSummary": fundraisingSummary,
		"investorEngagement": investorEngagement,
	})
}

// GetGrantsHandler returns available grants for founders
func (h *FounderHandler) GetGrantsHandler(c *fiber.Ctx) error {
	// Get query params
	category := c.Query("category")
	region := c.Query("region")

	// Get grants from database
	grants, err := h.db.Founder().GetGrants(c.Context(), category, region)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to get grants"})
	}

	return c.JSON(fiber.Map{
		"grants": grants,
	})
}

// SubmitGrantApplicationHandler handles grant application submission
func (h *FounderHandler) SubmitGrantApplicationHandler(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}

	// Fiber doesn't have ParseMultipartForm, it automatically parses the form
	// when you call FormValue or FormFile

	// Get form data
	grantIDStr := c.FormValue("grantId")
	startupName := c.FormValue("startupName")
	contactEmail := c.FormValue("contactEmail")
	contactPhone := c.FormValue("contactPhone")
	description := c.FormValue("description")
	website := c.FormValue("website")
	teamSize := c.FormValue("teamSize")
	previousFunding := c.FormValue("previousFunding")

	// Validate required fields
	if grantIDStr == "" || startupName == "" || contactEmail == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Missing required fields"})
	}

	// Convert grantID to int
	grantID, err := strconv.Atoi(grantIDStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid grant ID"})
	}

	// Get pitch deck file
	file, err := c.FormFile("pitchDeck")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Pitch deck is required"})
	}

	// Save file
	filename := fmt.Sprintf("%s-%s-%s", userID, time.Now().Format("20060102150405"), file.Filename)
	if err := c.SaveFile(file, "./uploads/"+filename); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to save file"})
	}

	// Create application
	founderID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	application := model.GrantApplication{
		ID:              primitive.NewObjectID(),
		FounderID:       founderID,
		GrantID:         grantID,
		StartupName:     startupName,
		ContactEmail:    contactEmail,
		ContactPhone:    contactPhone,
		Description:     description,
		Website:         website,
		TeamSize:        teamSize,
		PreviousFunding: previousFunding,
		PitchDeckPath:   "./uploads/" + filename,
		CreatedAt:       time.Now(),
	}

	// Save application to database
	appID, err := h.db.Founder().SubmitGrantApplication(c.Context(), application)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to submit application"})
	}

	return c.JSON(fiber.Map{
		"success":       true,
		"message":       "Application submitted successfully",
		"applicationId": appID,
	})
}

// GetInvestorsHandler returns available investors for founders
func (h *FounderHandler) GetInvestorsHandler(c *fiber.Ctx) error {
	// Get query params
	industry := c.Query("industry")
	stage := c.Query("stage")

	// Get investors from database
	investors, err := h.db.Founder().GetInvestors(c.Context(), industry, stage)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to get investors"})
	}

	return c.JSON(fiber.Map{
		"investors": investors,
	})
}

// SubmitInvestorApplicationHandler handles investor application submission
func (h *FounderHandler) SubmitInvestorApplicationHandler(c *fiber.Ctx) error {
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

	// Get user ID from claims
	userID := claims.UserID

	// Parse request body
	data := new(model.InvestorApplication)
	if err := c.BodyParser(data); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Validate required fields
	if data.InvestorID.IsZero() || data.FundingAmount == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Missing required fields"})
	}

	// Convert IDs to ObjectID
	founderID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	// No conversion needed as data.InvestorID is already primitive.ObjectID
	investorID := data.InvestorID
	
	// Create application
	application := model.InvestorApplication{
		ID:            primitive.NewObjectID(),
		FounderID:     founderID,
		InvestorID:    investorID,
		FundingAmount: data.FundingAmount,
		UseOfFunds:    data.UseOfFunds,
		// Add other fields as needed
	}

	// Save application to database
	appID, err := h.db.Founder().SubmitInvestorApplication(c.Context(), application)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to submit application"})
	}

	return c.JSON(fiber.Map{
		"success":       true,
		"message":       "Application submitted successfully",
		"applicationId": appID,
	})
}
