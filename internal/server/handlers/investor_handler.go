package handlers

import (
	"fmt"

	"DBackend/internal/database"
	"DBackend/model"
	"DBackend/utils"

	"github.com/gofiber/fiber/v2"
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
		InvestmentPortfolio   []string `json:"investment_portfolio"`
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
	investor.InvestmentPortfolio = data.InvestmentPortfolio
	investor.TotalInvested = data.TotalInvested
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

	return c.JSON(investor)
}
