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

// DealFlowHandler struct to handle deal flow requests
type DealFlowHandler struct {
	db database.Service
}

// NewDealFlowHandler creates a new instance of DealFlowHandler
func NewDealFlowHandler(db database.Service) *DealFlowHandler {
	return &DealFlowHandler{db: db}
}

// AddDealFlowHandler - Add a startup to deal flow
func (h *DealFlowHandler) AddDealFlowHandler(c *fiber.Ctx) error {
	deal := new(struct {
		UserID           string   `json:"UserID"`
		Name             string   `json:"Name"`
		Email            string   `json:"Email"`
		Avatar           string   `json:"Avatar"`
		StartupName      string   `json:"StartupName"`
		Industry         string   `json:"Industry"`
		FundingStage     string   `json:"FundingStage"`
		Location         string   `json:"Location"`
		FundRequired     int      `json:"FundRequired"`
		MissionStatement string   `json:"MissionStatement"`
		BusinessModel    string   `json:"BusinessModel"`
		RevenueStreams   string   `json:"RevenueStreams"`
		Traction         string   `json:"Traction"`
		ScalingPotential string   `json:"ScalingPotential"`
		Competition      string   `json:"Competition"`
		LeadershipTeam   string   `json:"LeadershipTeam"`
		TeamSize         string   `json:"TeamSize"`
		StartupWebsite   string   `json:"StartupWebsite"`
		PitchDeck        string   `json:"PitchDeck"`
		FundAllocation   string   `json:"FundAllocation"`
		Founded          string   `json:"Founded"`
		MatchScore       float64  `json:"MatchScore"`
		Tags             []string `json:"Tags"`
		Bookmarked       bool     `json:"Bookmarked"`
	})

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

	// Parse request body
	if err := c.BodyParser(&deal); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Validate required fields
	if deal.UserID == "" || deal.FundingStage == "" {
		return c.Status(400).JSON(fiber.Map{"error": "StartupID, InvestorID, Stage, and Status are required"})
	}

	startUpId, err := primitive.ObjectIDFromHex(deal.UserID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid Startup ID"})
	}
	// Create a new DealFlow instance
	newDeal := model.DealFlow{
		ID:         primitive.NewObjectID(),
		StartupID:  startUpId,
		InvestorID: investorID,
		MatchScore: deal.MatchScore,
		AddedDate:  time.Now(),
		Meetings:   []model.Meeting{},
		Documents:  []model.Document{},
		Tasks:      []model.Task{},
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	existingDeal, err := h.db.User().GetDealFlowByStartupID(c.Context(), startUpId)
	if err == nil && existingDeal.ID != primitive.NilObjectID {
		return c.Status(400).JSON(fiber.Map{"error": "Deal already exists"})
	}

	// Insert into database
	insertResult, err := h.db.User().AddStartupToDealFlow(c.Context(), newDeal)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to add deal to deal flow"})
	}

	return c.JSON(fiber.Map{"message": "Deal added successfully", "id": insertResult.InsertedID})
} // GetDealFlowByIDHandler - Retrieve a specific deal flow entry
func (h *DealFlowHandler) GetDealFlowByIDHandler(c *fiber.Ctx) error {
	id, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	deal, err := h.db.User().GetDealFlowByID(c.Context(), id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Deal not found"})
	}

	return c.JSON(deal)
}

// ListAllDealFlowHandler - Retrieve all deal flow entries with founder details
func (h *DealFlowHandler) ListAllDealFlowHandler(c *fiber.Ctx) error {
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

	deals, err := h.db.User().ListDealsByInvestorID(c.Context(), investorID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch deal flow entries"})
	}

	// Enrich deals with founder information and ensure required fields exist
	for i, deal := range deals {
		founder, err := h.db.User().FindByID(c.Context(), "users", deal["startup"].(bson.M)["user_id"].(primitive.ObjectID))
		if err == nil {
			founderUser, ok := founder.(*model.User)
			if ok {
				deals[i]["founder_name"] = founderUser.FirstName + " " + founderUser.SecondName
			}
			deals[i]["founder_email"] = founderUser.Email
		}

		// Ensure documents field exists (even if empty)
		if _, exists := deal["documents"]; !exists {
			deals[i]["documents"] = []interface{}{}
		}

		// Ensure other required fields exist
		requiredFields := []string{"meetings", "tasks", "last_activity"}
		for _, field := range requiredFields {
			if _, exists := deal[field]; !exists {
				deals[i][field] = []interface{}{}
			}
		}
	}

	return c.JSON(deals)
}

// UpdateDealFlowHandler - Update deal flow entry (stage, status, match score)
func (h *DealFlowHandler) UpdateDealFlowHandler(c *fiber.Ctx) error {
	id, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	fmt.Println(id)
	var updateFields map[string]interface{}
	if err := c.BodyParser(&updateFields); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Update database
	updateResult, err := h.db.User().UpdateDealFlow(c.Context(), id, updateFields)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update deal flow"})
	}

	return c.JSON(fiber.Map{"message": "Deal flow updated successfully", "modifiedCount": updateResult.ModifiedCount})
}

// DeleteDealFlowHandler - Remove startup from deal flow
func (h *DealFlowHandler) DeleteDealFlowHandler(c *fiber.Ctx) error {
	id, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	deleteResult, err := h.db.User().DeleteDealFlow(c.Context(), id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete deal flow entry"})
	}

	return c.JSON(fiber.Map{"message": "Deal flow entry deleted", "deletedCount": deleteResult.DeletedCount})
}

// AddMeetingHandler - Add a meeting to a deal flow entry
func (h *DealFlowHandler) AddMeetingHandler(c *fiber.Ctx) error {
	id, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid deal flow ID"})
	}

	var meeting model.Meeting
	if err := c.BodyParser(&meeting); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

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

// AddDocumentHandler - Add a document to a deal flow entry
func (h *DealFlowHandler) AddDocumentHandler(c *fiber.Ctx) error {
	id, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid deal flow ID"})
	}

	var document model.Document
	if err := c.BodyParser(&document); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	document.Date = time.Now()

	updateResult, err := h.db.User().AddDocument(c.Context(), id, document)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to add document"})
	}

	return c.JSON(fiber.Map{"message": "Document added successfully", "modifiedCount": updateResult.ModifiedCount})
}

// AddTaskHandler - Add a task to a deal flow entry
func (h *DealFlowHandler) AddTaskHandler(c *fiber.Ctx) error {
	id, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid deal flow ID"})
	}

	var task model.Task
	if err := c.BodyParser(&task); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	task.ID = primitive.NewObjectID()
	task.Completed = false // Default to incomplete

	updateResult, err := h.db.User().AddTask(c.Context(), id, task)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to add task"})
	}

	return c.JSON(fiber.Map{"message": "Task added successfully", "modifiedCount": updateResult.ModifiedCount})
}

// UpdateTaskStatusHandler - Update a task's status in a deal flow
func (h *DealFlowHandler) UpdateTaskStatusHandler(c *fiber.Ctx) error {
	dealID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid deal flow ID"})
	}

	taskID, err := primitive.ObjectIDFromHex(c.Params("taskID"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid task ID"})
	}

	var updateData struct {
		Completed bool `json:"completed"`
	}
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Implement the database update
	// This assumes you have a method to update a task's status
	updateResult, err := h.db.User().UpdateTaskStatus(c.Context(), dealID, taskID, updateData.Completed)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update task status"})
	}

	return c.JSON(fiber.Map{
		"message":       "Task status updated successfully",
		"modifiedCount": updateResult.ModifiedCount,
	})
}

// InvestInStartupHandler handles investment submissions
func (h *DealFlowHandler) InvestInStartupHandler(c *fiber.Ctx) error {
	// Get the deal ID from the URL parameters
	id, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid deal flow ID"})
	}

	// Parse the request body
	var request struct {
		InvestmentAmount float64 `json:"investmentAmount"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Validate the investment amount
	if request.InvestmentAmount <= 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Investment amount must be greater than zero"})
	}

	// Get token from context
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
	investorID, err := primitive.ObjectIDFromHex(claims.UserID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	// Get the deal from the database
	deal, err := h.db.DealFlow().GetDealFlowByID(c.Context(), id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Deal not found"})
	}

	// Create an investment record
	investment := model.Investment{
		ID:             primitive.NewObjectID(),
		DealID:         id,
		InvestorID:     investorID,
		FounderID:      deal.StartupID,
		Amount:         request.InvestmentAmount,
		InvestmentDate: time.Now(),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// Save the investment record
	_, err = h.db.Investment().CreateInvestment(c.Context(), investment)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to record investment"})
	}

	// Update startup's invested amount
	_, err = h.db.User().UpdateStartupInvestment(c.Context(), deal.StartupID, request.InvestmentAmount)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update startup investment"})
	}

	// Update investor's total investment and portfolio
	_, err = h.db.User().UpdateInvestorPortfolio(c.Context(), investorID, deal.StartupID, request.InvestmentAmount)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update investor investment"})
	}

	// Update the deal's fund required amount
	_, err = h.db.User().UpdateDealFundRequired(c.Context(), id, -request.InvestmentAmount)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update deal fund required"})
	}

	return c.JSON(fiber.Map{
		"message":    "Investment recorded successfully",
		"investment": investment,
	})
}

// UpdateDealStatusHandler updates the status of a deal
func (h *DealFlowHandler) UpdateDealStatusHandler(c *fiber.Ctx) error {
	dealID := c.Params("id")
	if dealID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Deal ID is required"})
	}

	// Parse request body
	data := struct {
		Status string `json:"status"`
	}{}
	if err := c.BodyParser(&data); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Validate status
	validStatuses := map[string]bool{"active": true, "paused": true, "completed": true, "cancelled": true}
	if !validStatuses[data.Status] {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid status"})
	}

	// Convert ID string to ObjectID
	objID, err := primitive.ObjectIDFromHex(dealID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid deal ID format"})
	}

	// Update deal status
	result, err := h.db.DealFlow().UpdateDealFlow(c.Context(), objID, bson.M{"status": data.Status})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update deal status"})
	}

	return c.JSON(fiber.Map{"message": "Deal status updated successfully", "result": result})
}

// UpdateDealStageHandler updates the pipeline stage of a deal
func (h *DealFlowHandler) UpdateDealStageHandler(c *fiber.Ctx) error {
	dealID := c.Params("id")
	if dealID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Deal ID is required"})
	}

	// Parse request body
	data := struct {
		Stage string `json:"stage"`
	}{}
	if err := c.BodyParser(&data); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Validate stage
	validStages := map[string]bool{"screening": true, "dueDiligence": true, "negotiation": true, "closed": true}
	if !validStages[data.Stage] {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid stage"})
	}

	// Convert dealID to ObjectID
	objID, err := primitive.ObjectIDFromHex(dealID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid deal ID"})
	}

	// Update deal stage
	_, err = h.db.DealFlow().UpdateDealStage(c.Context(), objID, data.Stage)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update deal stage"})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Deal stage updated successfully",
		"deal": fiber.Map{
			"_id":   dealID,
			"stage": data.Stage,
		},
	})
}

// AddNoteHandler adds a note to a deal
func (h *DealFlowHandler) AddNoteHandler(c *fiber.Ctx) error {
	dealID := c.Params("id")
	if dealID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Deal ID is required"})
	}

	// Parse request body
	data := struct {
		Content string `json:"content"`
	}{}
	if err := c.BodyParser(&data); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if data.Content == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Note content is required"})
	}

	// Convert dealID to ObjectID
	objID, err := primitive.ObjectIDFromHex(dealID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid deal ID"})
	}

	// Create note
	note := model.Note{
		ID:        primitive.NewObjectID(),
		Content:   data.Content,
		CreatedAt: time.Now(),
	}

	// Add note to deal
	err = h.db.DealFlow().AddNote(c.Context(), objID, note)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to add note"})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"note": fiber.Map{
			"id":         note.ID.Hex(),
			"content":    note.Content,
			"created_at": note.CreatedAt,
		},
	})
}

// UpdateTaskStatusHandler updates the status of a task
// func (h *DealFlowHandler) UpdateTaskStatusHandler(c *fiber.Ctx) error {
// 	taskID := c.Params("taskId")
// 	if taskID == "" {
// 		return c.Status(400).JSON(fiber.Map{"error": "Task ID is required"})
// 	}

// 	// Parse request body
// 	data := struct {
// 		Completed bool `json:"completed"`
// 	}{}
// 	if err := c.BodyParser(&data); err != nil {
// 		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
// 	}

// 	// Convert taskID to ObjectID
// 	objID, err := primitive.ObjectIDFromHex(taskID)
// 	if err != nil {
// 		return c.Status(400).JSON(fiber.Map{"error": "Invalid task ID"})
// 	}

// 	// Update task status
// 	_, err = h.db.User().UpdateTaskStatus(c.Context(), objID, objID, data.Completed)
// 	if err != nil {
// 		return c.Status(500).JSON(fiber.Map{"error": "Failed to update task status"})
// 	}

// 	return c.JSON(fiber.Map{
// 		"success": true,
// 		"task": fiber.Map{
// 			"id":        taskID,
// 			"completed": data.Completed,
// 		},
// 	})
// }
