package handlers

import (
	"encoding/json"
	"fmt"
	"os"

	"DBackend/internal/database"
	"DBackend/model"
	"DBackend/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FounderHandler struct {
	db database.Service
}

func NewFounderHandler(db database.Service) *FounderHandler {
	return &FounderHandler{db: db}
}

// UpdateFounderHandler handles updating founder profile including pitch deck uploads
func (h *FounderHandler) UpdateFounderHandler(c *fiber.Ctx) error {
	// Parse JSON data from form-data
	data := new(struct {
		StartupName       string `json:"startup_name"`
		MissionStatement  string `json:"mission_statement"`
		Industry          string `json:"industry"`
		FundingStage      string `json:"funding_stage"`
		FundingAllocation string `json:"funding_allocation"`
		BusinessModel     string `json:"business_model"`
		RevenueStreams    string `json:"revenue_streams"`
		Traction          string `json:"traction"`
		ScalingPotential  string `json:"scaling_potential"`
		Competition       string `json:"competition"`
		LeadershipTeam    string `json:"leadership_team"`
		TeamSize          int    `json:"team_size"`
		Location          string `json:"location"`
		StartupWebsite    string `json:"startup_website"`
	})

	// Extract JSON data from form field "data"
	jsonData := c.FormValue("data")
	if jsonData == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Missing JSON data"})
	}

	fmt.Println("data", jsonData)
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

	// Update user's profile with provided data
	user.StartupName = data.StartupName
	user.MissionStatement = data.MissionStatement
	user.Industry = data.Industry
	user.FundingStage = data.FundingStage
	user.FundingAllocation = data.FundingAllocation
	user.BusinessModel = data.BusinessModel
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

	return c.JSON(fiber.Map{"message": "Founder profile updated successfully", "pitch_deck": user.PitchDeck})
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
		return c.Status(401).JSON(fiber.Map{"error": "Invalid token"})
	}

	// Convert claims.UserID to primitive.ObjectID
	id, err := primitive.ObjectIDFromHex(claims.UserID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	// Find the founder in the database
	userInterface, err := h.db.User().FindByID(c.Context(), "founders", id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Founder not found"})
	}

	// Assert correct type for the founder object
	founder, ok := userInterface.(*model.Founder)
	if !ok {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to parse founder data"})
	}

	return c.JSON(founder)
}

