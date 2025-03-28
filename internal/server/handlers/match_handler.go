package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	// import ioutil
	"io/ioutil"

	"DBackend/internal/database"
	"DBackend/model"
	"DBackend/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MatHandler struct to handle Match-related operations
type MatHandler struct {
	db database.Service
}

// NewInvestorHandler initializes an InvestorHandler
func NewMatHandler(db database.Service) *MatHandler {
	return &MatHandler{db: db}
}

// MatchRequest defines the structure of the payload that our FastAPI endpoint expects.
type MatchRequest struct {
	Founder  map[string]interface{} `json:"founder"`
	Investor map[string]interface{} `json:"investor"`
}

// MatchResponse defines the structure of the response from FastAPI.
type MatchResponse struct {
	MatchProbability float64 `json:"match_probability"`
}

func (h *MatHandler) MatchHandler(c *fiber.Ctx) error {
	// Extract founder and investor IDs from request
	founderID := c.Params("userID")

	tokenStr, ok := c.Locals("token").(string)
	if !ok || tokenStr == "" {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid token"})
	}

	claims, err := utils.ValidateJWT(tokenStr)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid token"})
	}

	// Convert claims.UserID to primitive.ObjectID
	investorID, err := primitive.ObjectIDFromHex(claims.UserID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid user ID"})
	}
	fmt.Println("founderID", founderID)
	fmt.Println("investorID", investorID)
	// Convert IDs to ObjectID
	fID, err := primitive.ObjectIDFromHex(founderID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid founder ID",
		})
	}

	// Fetch founder and investor details from the database
	founder, err := h.db.User().FindByID(c.Context(), "founders", fID)
	fmt.Println(err)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Founder not found",
		})
	}
	investor, err := h.db.User().FindByID(c.Context(), "investors", investorID)
	fmt.Println(err)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Investor not found",
		})
	}

	// Type assert founder and investor to their respective models
	founderDetails, ok := founder.(*model.Founder)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}
	investorDetails, ok := investor.(*model.Investor)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	// Build the payload with default values for empty fields
	reqPayload := MatchRequest{
		Founder: map[string]interface{}{
			"fund_required": getDefaultIfZero(float64(founderDetails.FundRequired), 500000),
			"industry":      getDefaultIfEmpty(founderDetails.Industry, "Other"),
			"funding_stage": getDefaultIfEmpty(founderDetails.FundingStage, "Seed"),
		},
		Investor: map[string]interface{}{
			"total_invested":          getDefaultIfZero(investorDetails.TotalInvested, 1000000),
			"preferred_funding_stage": getDefaultIfEmpty(investorDetails.PreferredFundingStage, "Seed"),
			"risk_tolerance":          getDefaultIfEmpty(investorDetails.RiskTolerance, "Moderate"),
		},
	}

	// Call the ML service
	matchProbability, err := GetMatchProbability(reqPayload)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Error getting match probability: %v", err),
		})
	}

	match := model.MatchInvestorFounder{
		FounderID:       fID,
		InvestorID:      investorID,
		MatchPercentage: matchProbability,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	_, err = h.db.User().AddMatch(c.Context(), match)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to add match"})
	}
	fmt.Println(matchProbability)

	// Return the match probability
	return c.Status(201).JSON(fiber.Map{
		"match_probability": matchProbability,
	})
}

// GetMatchProbability sends a POST request to the FastAPI endpoint and returns the match probability.
func GetMatchProbability(req MatchRequest) (float64, error) {
	url := "http://127.0.0.1:4040/predict/" // FastAPI endpoint URL

	// Marshal the request into JSON.
	requestBody, err := json.Marshal(req)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal request: %v", err)
	}

	// Send the POST request.
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return 0, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Read and parse the response body.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read response: %v", err)
	}

	var matchResp MatchResponse
	if err := json.Unmarshal(body, &matchResp); err != nil {
		return 0, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return matchResp.MatchProbability, nil
}

// CalculateMatchHandler calculates and stores match between investor and founder
func (h *MatHandler) CalculateMatchHandler(c *fiber.Ctx) error {
	// Parse request
	var req struct {
		FounderID  string `json:"founder_id"`
		InvestorID string `json:"investor_id"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Validate IDs
	fID, err := primitive.ObjectIDFromHex(req.FounderID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid founder ID"})
	}

	investorID, err := primitive.ObjectIDFromHex(req.InvestorID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid investor ID"})
	}

	// Get founder and investor details
	founderDetails, err := h.db.User().GetFounderByUserID(c.Context(), fID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Founder not found"})
	}

	investorDetails, err := h.db.User().FindByID(c.Context(), "investors", investorID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Investor not found"})
	}

	investorObj, ok := investorDetails.(*model.Investor)
	if !ok {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to parse investor data"})
	}

	// Build the payload
	reqPayload := MatchRequest{
		Founder: map[string]interface{}{
			"fund_required": founderDetails.FundRequired,
			"industry":      founderDetails.Industry,
			"funding_stage": founderDetails.FundingStage,
		},
		Investor: map[string]interface{}{
			"total_invested":          investorObj.TotalInvested,
			"preferred_funding_stage": investorObj.PreferredFundingStage,
			"risk_tolerance":          investorObj.RiskTolerance,
		},
	}

	// Call the ML service
	matchProbability, err := GetMatchProbability(reqPayload)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Error getting match probability: %v", err),
		})
	}

	match := model.MatchInvestorFounder{
		FounderID:       fID,
		InvestorID:      investorID,
		MatchPercentage: matchProbability,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	_, err = h.db.User().AddMatch(c.Context(), match)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to add match"})
	}
	fmt.Println(matchProbability)

	// Return the match probability
	return c.Status(201).JSON(fiber.Map{
		"match_probability": matchProbability,
	})
}

// Helper functions to provide defaults for empty values
func getDefaultIfEmpty(value string, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}

func getDefaultIfZero(value float64, defaultValue float64) float64 {
	if value == 0 {
		return defaultValue
	}
	return value
}
