package database

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"DBackend/model"

	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var dbName = os.Getenv("BLUEPRINT_DB_DATABASE")

// UserService interface
type UserService interface {
	GetUserCount(ctx context.Context) (int64, error)
	AddStartupToDealFlow(ctx context.Context, deal model.DealFlow) (*mongo.InsertOneResult, error)
	GetDealFlowByID(ctx context.Context, id primitive.ObjectID) (*model.DealFlow, error)
	ListAllDealFlow(ctx context.Context) ([]bson.M, error)
	UpdateDealFlow(ctx context.Context, id primitive.ObjectID, updateFields bson.M) (*mongo.UpdateResult, error)
	UpdateStartupInvestment(ctx context.Context, startupID primitive.ObjectID, amount float64) (*mongo.UpdateResult, error)
	UpdateInvestorPortfolio(ctx context.Context, investorID, startupID primitive.ObjectID, amount float64) (*mongo.UpdateResult, error)
	DeleteDealFlow(ctx context.Context, id primitive.ObjectID) (*mongo.DeleteResult, error)
	AddDocument(ctx context.Context, dealID primitive.ObjectID, document model.Document) (*mongo.UpdateResult, error)
	AddTask(ctx context.Context, dealID primitive.ObjectID, task model.Task) (*mongo.UpdateResult, error)
	UpdateTaskStatus(ctx context.Context, dealID primitive.ObjectID, taskID primitive.ObjectID, completed bool) (*mongo.UpdateResult, error)
	AddNote(ctx context.Context, dealID primitive.ObjectID, note model.Note) error
	//UpdateDealStatus(ctx context.Context, objID primitive.ObjectID, status string) (any, error)
	UpdateDealStage(ctx context.Context, objID primitive.ObjectID, stage string) (any, error)
	// notification
	notifyFounder(ctx context.Context, founderID primitive.ObjectID, message string) error

	GetNotificationsByFounder(ctx context.Context, founderID primitive.ObjectID) ([]model.Notification, error)
	UpdateNotification(ctx context.Context, notificationID primitive.ObjectID, updateData bson.M) error
	DeleteNotification(ctx context.Context, notificationID primitive.ObjectID) (*mongo.DeleteResult, error)
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	// find by id of any model
	GetDealFlowByStartupID(ctx context.Context, startupID primitive.ObjectID) (model.DealFlow, error)

	GetMeetings(ctx context.Context, userID primitive.ObjectID) ([]model.Meeting, error)
	FindByID(ctx context.Context, collectionName string, id primitive.ObjectID) (interface{}, error)

	ListDealsByInvestorID(ctx context.Context, investorID primitive.ObjectID) ([]bson.M, error)
	UpdateRoles(ctx context.Context, email string, roles []string) (*mongo.UpdateResult, error)
	CreateUser(ctx context.Context, user model.User) (*mongo.InsertOneResult, error)
	//	GetFounderByUserID(ctx context.Context, userID primitive.ObjectID) (model.Founder, error)
	UpdateFounder(ctx context.Context, userID primitive.ObjectID, founder model.Founder) (*mongo.UpdateResult, error)
	UpdateInvestor(ctx context.Context, userID primitive.ObjectID, investor model.Investor) (*mongo.UpdateResult, error)
	CreateRoleData(ctx context.Context, userID primitive.ObjectID, role string) error
	GetFounderProfileWithMatch(ctx context.Context, userID primitive.ObjectID) (bson.M, error)
	BlacklistToken(ctx context.Context, token string) error
	IsTokenBlacklisted(ctx context.Context, token string) (bool, error)
	GetStartupDetails(ctx context.Context) ([]model.Founder, error)
	GetInvestorDetails(ctx context.Context) ([]model.Investor, error)
	// notification
	GetAllNotificationsByFounder(ctx context.Context, founderID primitive.ObjectID) ([]model.Notification, error)

	// match
	UpdateMatch(ctx context.Context, matchID primitive.ObjectID, updateFields bson.M) (*mongo.UpdateResult, error)
	AddMatch(ctx context.Context, match model.MatchInvestorFounder) (*mongo.InsertOneResult, error)

	//UpdateStartupInvestment(ctx context.Context, startupID primitive.ObjectID, amount float64) (*mongo.UpdateResult, error)

	GetFounderByUserID(ctx context.Context, userID primitive.ObjectID) (model.Founder, error)
	//UpdateInvestorPortfolio(ctx context.Context, investorID, startupID primitive.ObjectID, amount float64) (*mongo.UpdateResult, error)
	UpdateMeeting(ctx context.Context, id primitive.ObjectID, updates model.Meeting) error

	GetAllTasks(ctx context.Context) ([]model.Task, error)
	GetTaskByID(ctx context.Context, id primitive.ObjectID) (model.Task, error)
	UpdateTask(ctx context.Context, id primitive.ObjectID, updates model.Task) error
	DeleteTask(ctx context.Context, id primitive.ObjectID) error
	GetTasksByUser(ctx context.Context, userID primitive.ObjectID) ([]model.Task, error)
	//UpdateTaskStatus(ctx context.Context, dealID primitive.ObjectID, taskID primitive.ObjectID, completed bool) (*mongo.UpdateResult, error)

	AssignTask(ctx context.Context, taskID primitive.ObjectID, userID primitive.ObjectID) error

	AddMeetingNotes(ctx context.Context, meetingID primitive.ObjectID, notes string) error
	GetMeetingNotes(ctx context.Context, meetingID primitive.ObjectID) (string, error)

	UpdateDealStatus(ctx *fasthttp.RequestCtx, objID primitive.ObjectID, status string) (any, error)
	AddMeetingParticipant(ctx context.Context, meetingID, userID primitive.ObjectID) error
	RemoveMeetingParticipant(ctx context.Context, meetingID, userID primitive.ObjectID) error
	UpdateDealFundRequired(ctx *fasthttp.RequestCtx, id primitive.ObjectID, f float64) (any, error)


}
// userService struct
type userService struct {
	userCollection                 *mongo.Collection
	founderCollection              *mongo.Collection
	notificationCollection         *mongo.Collection
	investorCollection             *mongo.Collection
	adminCollection                *mongo.Collection
	dealFlowCollection             *mongo.Collection
	meetingCollection              *mongo.Collection
	matchFounderInvestorCollection *mongo.Collection
	blacklistCollection            *mongo.Collection
	taskCollection                 *mongo.Collection
}

// NewUserService initializes collections
func NewUserService(client *mongo.Client) UserService {
	// Use environment variable with fallback
	if dbName == "" {
		dbName = "ddb" // Fallback name
	}

	return &userService{
		userCollection:                 client.Database(dbName).Collection("users"),
		notificationCollection:         client.Database(dbName).Collection("notifications"),
		founderCollection:              client.Database(dbName).Collection("founders"),
		investorCollection:             client.Database(dbName).Collection("investors"),
		taskCollection:                 client.Database(dbName).Collection("tasks"),
		adminCollection:                client.Database(dbName).Collection("admins"),
		blacklistCollection:            client.Database(dbName).Collection("blacklist_tokens"),
		dealFlowCollection:             client.Database(dbName).Collection("deal_flow"),
		meetingCollection:              client.Database(dbName).Collection("meetings"),
		matchFounderInvestorCollection: client.Database(dbName).Collection("match_founder_investor"),
	}
}

// BlacklistToken stores a JWT token in the blacklist collection
func (s *userService) BlacklistToken(ctx context.Context, token string) error {
	_, err := s.blacklistCollection.InsertOne(ctx, bson.M{
		"token":     token,
		"expiredAt": time.Now().Add(24 * time.Hour), // Set expiry for 24h
	})

	return err
}

// IsTokenBlacklisted checks if a token is blacklisted
func (s *userService) IsTokenBlacklisted(ctx context.Context, token string) (bool, error) {
	count, err := s.blacklistCollection.CountDocuments(ctx, bson.M{"token": token})
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// Find user by email
func (s *userService) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := s.userCollection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// FindByID retrieves an object from a specified collection based on the provided ID
func (s *userService) FindByID(ctx context.Context, collectionName string, id primitive.ObjectID) (interface{}, error) {
	var collection *mongo.Collection
	var result interface{}

	// Map the collectionName to the correct collection and model
	switch collectionName {
	case "users":
		collection = s.userCollection
		result = &model.User{}
	case "founders":
		collection = s.founderCollection
		result = &model.Founder{}
	case "investors":
		collection = s.investorCollection
		result = &model.Investor{}
	case "admins":
		collection = s.adminCollection
		result = &model.Admin{}
	default:
		return nil, errors.New("invalid collection name")
	}

	// Query the specified collection
	filter := bson.M{"_id": id}
	if collectionName != "users" {
		filter = bson.M{"user_id": id}
	}

	err := collection.FindOne(ctx, filter).Decode(result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("object not found")
		}
		return nil, err
	}

	return result, nil
}

// Update user roles
func (s *userService) UpdateRoles(ctx context.Context, email string, roles []string) (*mongo.UpdateResult, error) {
	filter := bson.M{"email": email}
	update := bson.M{"$set": bson.M{"roles": roles}}
	return s.userCollection.UpdateOne(ctx, filter, update)
}

// Create a new user
func (s *userService) CreateUser(ctx context.Context, user model.User) (*mongo.InsertOneResult, error) {
	return s.userCollection.InsertOne(ctx, user)
}

// Create role-specific data
func (s *userService) CreateRoleData(ctx context.Context, userID primitive.ObjectID, role string) error {
	switch role {
	case "founder":
		_, err := s.founderCollection.InsertOne(ctx, model.Founder{UserID: userID})
		return err
	case "investor":
		_, err := s.investorCollection.InsertOne(ctx, model.Investor{UserID: userID})
		return err
	case "admin":
		_, err := s.adminCollection.InsertOne(ctx, model.Admin{UserID: userID})
		return err
	default:
		return errors.New("invalid role")
	}
}

func (s *userService) UpdateFounder(ctx context.Context, userID primitive.ObjectID, founder model.Founder) (*mongo.UpdateResult, error) {
	filter := bson.M{"user_id": userID}
	updateFields := bson.M{
		"startup_name":       founder.StartupName,
		"mission_statement":  founder.MissionStatement,
		"industry":           founder.Industry,
		"funding_stage":      founder.FundingStage,
		"funding_allocation": founder.FundingAllocation,
		"bussiness_model":    founder.BussinessModel,
		"revenue_streams":    founder.RevenueStreams,
		"traction":           founder.Traction,
		"scaling_potential":  founder.ScalingPotential,
		"total_invested":     founder.TotalInvested,
		"fund_required":      founder.FundRequired,
		"competition":        founder.Competition,
		"leadership_team":    founder.LeadershipTeam,
		"team_size":          founder.TeamSize,
		"location":           founder.Location,
		"startup_website":    founder.StartupWebsite,
	}
	update := bson.M{"$set": updateFields}
	return s.founderCollection.UpdateOne(ctx, filter, update)
}

func (s *userService) UpdateInvestor(ctx context.Context, userID primitive.ObjectID, investor model.Investor) (*mongo.UpdateResult, error) {
	filter := bson.M{"user_id": userID}
	updateFields := bson.M{
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
	}
	update := bson.M{"$set": updateFields}
	return s.investorCollection.UpdateOne(ctx, filter, update)
}

// get all the registered founders
func (s *userService) GetStartupDetails(ctx context.Context) ([]model.Founder, error) {
	var founders []model.Founder
	cursor, err := s.founderCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var founder model.Founder
		if err := cursor.Decode(&founder); err != nil {
			return nil, err
		}
		founders = append(founders, founder)
	}

	return founders, nil
}

// get all registered investors
func (s *userService) GetInvestorDetails(ctx context.Context) ([]model.Investor, error) {
	var investors []model.Investor
	cursor, err := s.investorCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var investor model.Investor
		if err := cursor.Decode(&investor); err != nil {
			return nil, err
		}
		investors = append(investors, investor)
	}
	return investors, nil
}

// GetFounderProfileWithMatch retrieves a founder profile with nested user details,
// adds a "match" field as a URL with the founder id as a parameter, and sets "tags" and "bookmark" empty.
func (s *userService) GetFounderProfileWithMatch(ctx context.Context, userID primitive.ObjectID) (bson.M, error) {
	// Build the aggregation pipeline

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.D{{Key: "user_id", Value: userID}}}},
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "users"},
			{Key: "localField", Value: "user_id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "user"},
		}}},
		// Stage 3: Unwind the user array.
		{{Key: "$unwind", Value: "$user"}},
		// Stage 4: Add computed fields: match URL, empty tags, and bookmark false.
		{{Key: "$addFields", Value: bson.D{
			{Key: "tags", Value: bson.A{}},
			{Key: "match", Value: "http://localhost:8080/api/v1/match/data/" + userID.Hex()},
			{Key: "bookmark", Value: false},
		}}},
		// Stage 5: Project the desired fields.
		{{Key: "$project", Value: bson.D{
			{Key: "startup_name", Value: 1},
			{Key: "mission_statement", Value: 1},
			{Key: "industry", Value: 1},
			{Key: "funding_stage", Value: 1},
			{Key: "funding_allocation", Value: 1},
			{Key: "bussiness_model", Value: 1},
			{Key: "revenue_streams", Value: 1},
			{Key: "traction", Value: 1},
			{Key: "total_invested", Value: 1},
			{Key: "fund_required", Value: 1},
			{Key: "year_founded", Value: 1},
			{Key: "scaling_potential", Value: 1},
			{Key: "competition", Value: 1},
			{Key: "leadership_team", Value: 1},
			{Key: "team_size", Value: 1},
			{Key: "avatar", Value: 1},
			{Key: "founded", Value: 1},
			{Key: "location", Value: 1},
			{Key: "startup_website", Value: 1},
			{Key: "pitch_deck", Value: 1},
			{Key: "created_at", Value: 1},
			{Key: "updated_at", Value: 1},
			{Key: "match", Value: 1},
			{Key: "tags", Value: 1},
			{Key: "bookmark", Value: 1},
			{Key: "user.first_name", Value: 1},
			{Key: "user.second_name", Value: 1},
			{Key: "user.email", Value: 1},
			{Key: "user.avatar", Value: 1},
			{Key: "user.created_at", Value: 1},
			{Key: "user.updated_at", Value: 1},
		}}},
	}

	cursor, err := s.founderCollection.Aggregate(ctx, pipeline)
	if err != nil {
		log.Println("Error running aggregation:", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, mongo.ErrNoDocuments
	}

	return results[0], nil
}

// AddMatch adds a new match between an investor and a founder
func (s *userService) AddMatch(ctx context.Context, match model.MatchInvestorFounder) (*mongo.InsertOneResult, error) {
	match.CreatedAt = time.Now()
	match.UpdatedAt = time.Now()
	return s.matchFounderInvestorCollection.InsertOne(ctx, match)
}

// UpdateMatch updates an existing match between an investor and a founder
func (s *userService) UpdateMatch(ctx context.Context, matchID primitive.ObjectID, updateFields bson.M) (*mongo.UpdateResult, error) {
	updateFields["updated_at"] = time.Now()
	update := bson.M{"$set": updateFields}
	return s.matchFounderInvestorCollection.UpdateOne(ctx, bson.M{"_id": matchID}, update)
}

// GetNotificationsByFounder retrieves notifications for a specific founder.
func (s *userService) GetNotificationsByFounder(ctx context.Context, founderID primitive.ObjectID) ([]model.Notification, error) {
	var notifications []model.Notification
	cursor, err := s.notificationCollection.Find(ctx, bson.M{"founder_id": founderID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var notification model.Notification
		if err := cursor.Decode(&notification); err != nil {
			return nil, err
		}
		notifications = append(notifications, notification)
	}

	return notifications, nil
}

// UpdateNotification updates a specific notification by its ID.
func (s *userService) UpdateNotification(ctx context.Context, notificationID primitive.ObjectID, updateData bson.M) error {
	update := bson.M{
		"$set": updateData,
	}
	_, err := s.notificationCollection.UpdateOne(ctx, bson.M{"_id": notificationID}, update)
	return err
}

// DeleteNotification deletes a specific notification.
func (s *userService) DeleteNotification(ctx context.Context, notificationID primitive.ObjectID) (*mongo.DeleteResult, error) {
	fmt.Println("notificationID", notificationID)
	return s.notificationCollection.DeleteOne(ctx, bson.M{"_id": notificationID})
}

// GetAllNotificationsByFounder retrieves all notifications for a specific founder, ordered by the latest first.
func (s *userService) GetAllNotificationsByFounder(ctx context.Context, founderID primitive.ObjectID) ([]model.Notification, error) {
	var notifications []model.Notification
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
	cursor, err := s.notificationCollection.Find(ctx, bson.M{"founder_id": founderID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var notification model.Notification
		if err := cursor.Decode(&notification); err != nil {
			return nil, err
		}
		notifications = append(notifications, notification)
	}

	return notifications, nil
}

func (s *userService) GetFounderByUserID(ctx context.Context, userID primitive.ObjectID) (model.Founder, error) {
	var founder model.Founder
	// get the founder where user_id = userID
	err := s.founderCollection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&founder)
	if err != nil {
		return model.Founder{}, err
	}
	return founder, nil
}

// GetMeetings retrieves all meetings for a given user ID
func (s *userService) GetMeetings(ctx context.Context, userID primitive.ObjectID) ([]model.Meeting, error) {
	var meetings []model.Meeting
	cursor, err := s.meetingCollection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var meeting model.Meeting
		if err := cursor.Decode(&meeting); err != nil {
			return nil, err
		}
		meetings = append(meetings, meeting)
	}

	return meetings, nil
}

func (s *userService) GetDealFlowByStartupID(ctx context.Context, startupID primitive.ObjectID) (model.DealFlow, error) {
	var deal model.DealFlow
	err := s.dealFlowCollection.FindOne(ctx, bson.M{"startup_id": startupID}).Decode(&deal)
	if err != nil {
		return model.DealFlow{}, err
	}
	return deal, nil
}

// UpdateMeeting updates an existing meeting by ID
func (s *userService) UpdateMeeting(ctx context.Context, id primitive.ObjectID, updates model.Meeting) error {
	updates.UpdatedAt = time.Now()

	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{
		"title":           updates.Title,
		"description":     updates.Notes,
		"start_time":      updates.StartTime,
		"end_time":        updates.EndTime,
		"google_meet_url": updates.GoogleMeetURL,
		"updated_at":      updates.UpdatedAt,
	}}

	result, err := s.meetingCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("meeting not found")
	}

	return nil
}

// GetAllTasks retrieves all tasks from the database
func (s *userService) GetAllTasks(ctx context.Context) ([]model.Task, error) {
	var tasks []model.Task
	cursor, err := s.dealFlowCollection.Aggregate(ctx, []bson.M{
		{"$unwind": "$tasks"},
		{"$project": bson.M{"task": "$tasks"}},
		{"$replaceRoot": bson.M{"newRoot": "$task"}},
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

// GetTaskByID retrieves a specific task by ID
func (s *userService) GetTaskByID(ctx context.Context, id primitive.ObjectID) (model.Task, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"tasks._id": id}}},
		{{Key: "$unwind", Value: "$tasks"}},
		{{Key: "$match", Value: bson.M{"tasks._id": id}}},
		{{Key: "$replaceRoot", Value: bson.M{"newRoot": "$tasks"}}},
	}

	cursor, err := s.dealFlowCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return model.Task{}, err
	}
	defer cursor.Close(ctx)

	var tasks []model.Task
	if err = cursor.All(ctx, &tasks); err != nil {
		return model.Task{}, err
	}

	if len(tasks) == 0 {
		return model.Task{}, mongo.ErrNoDocuments
	}

	return tasks[0], nil
}

// UpdateTask updates an existing task
func (s *userService) UpdateTask(ctx context.Context, id primitive.ObjectID, updates model.Task) error {
	updates.UpdatedAt = time.Now()

	filter := bson.M{"tasks._id": id}
	update := bson.M{
		"$set": bson.M{
			"tasks.$.title":      updates.Title,
			"tasks.$.completed":  updates.Completed,
			"tasks.$.due_date":   updates.DueDate,
			"tasks.$.priority":   updates.Priority,
			"tasks.$.updated_at": updates.UpdatedAt,
		},
	}

	_, err := s.dealFlowCollection.UpdateOne(ctx, filter, update)
	return err
}

// DeleteTask deletes a task
func (s *userService) DeleteTask(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{}
	update := bson.M{"$pull": bson.M{"tasks": bson.M{"_id": id}}}

	_, err := s.dealFlowCollection.UpdateMany(ctx, filter, update)
	return err
}

// GetTasksByUser retrieves all tasks for a specific user
func (s *userService) GetTasksByUser(ctx context.Context, userID primitive.ObjectID) ([]model.Task, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"tasks.created_by": userID}}},
		{{Key: "$unwind", Value: "$tasks"}},
		{{Key: "$match", Value: bson.M{"tasks.created_by": userID}}},
		{{Key: "$replaceRoot", Value: bson.M{"newRoot": "$tasks"}}},
	}

	cursor, err := s.dealFlowCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tasks []model.Task
	if err = cursor.All(ctx, &tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}

// AssignTask assigns a task to a user
func (s *userService) AssignTask(ctx context.Context, taskID primitive.ObjectID, userID primitive.ObjectID) error {
	filter := bson.M{"tasks._id": taskID}
	update := bson.M{
		"$set": bson.M{
			"tasks.$.assigned_to": userID,
			"tasks.$.updated_at":  time.Now(),
		},
	}

	_, err := s.dealFlowCollection.UpdateOne(ctx, filter, update)
	return err
}

// UpdateTaskCompletionStatus updates the completion status of a task in a deal flow
func (s *userService) UpdateTaskCompletionStatus(ctx context.Context, dealID primitive.ObjectID, taskID primitive.ObjectID, completed bool) (*mongo.UpdateResult, error) {
	filter := bson.M{"_id": dealID, "tasks._id": taskID}
	update := bson.M{"$set": bson.M{
		"tasks.$.completed":  completed,
		"tasks.$.updated_at": time.Now(),
	}}
	return s.dealFlowCollection.UpdateOne(ctx, filter, update)
}

// AddMeetingNotes adds notes to a meeting
func (s *userService) AddMeetingNotes(ctx context.Context, meetingID primitive.ObjectID, notes string) error {
	filter := bson.M{"_id": meetingID}
	update := bson.M{
		"$set": bson.M{
			"notes":      notes,
			"updated_at": time.Now(),
		},
	}

	result, err := s.meetingCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("meeting not found")
	}

	return nil
}

// GetMeetingNotes retrieves notes for a meeting
func (s *userService) GetMeetingNotes(ctx context.Context, meetingID primitive.ObjectID) (string, error) {
	filter := bson.M{"_id": meetingID}
	var meeting struct {
		Notes string `bson:"notes"`
	}

	err := s.meetingCollection.FindOne(ctx, filter).Decode(&meeting)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", errors.New("meeting not found")
		}
		return "", err
	}

	return meeting.Notes, nil
}

// AddMeetingParticipant adds a participant to a meeting
func (s *userService) AddMeetingParticipant(ctx context.Context, meetingID, userID primitive.ObjectID) error {
	filter := bson.M{"_id": meetingID}
	update := bson.M{
		"$addToSet": bson.M{"participants": userID},
		"$set":      bson.M{"updated_at": time.Now()},
	}

	result, err := s.meetingCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("meeting not found")
	}

	return nil
}

// RemoveMeetingParticipant removes a participant from a meeting
func (s *userService) RemoveMeetingParticipant(ctx context.Context, meetingID, userID primitive.ObjectID) error {
	filter := bson.M{"_id": meetingID}
	update := bson.M{
		"$pull": bson.M{"participants": userID},
		"$set":  bson.M{"updated_at": time.Now()},
	}

	result, err := s.meetingCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("meeting not found")
	}

	return nil
}

func (s *userService) GetUserCount(ctx context.Context) (int64, error) {
	count, err := s.userCollection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return 0, err
	}
	return count, nil
}

// UpdateDealFundRequired updates the fund required amount for a deal
func (s *userService) UpdateDealFundRequired(ctx *fasthttp.RequestCtx, dealID primitive.ObjectID, amountChange float64) (any, error) {
	filter := bson.M{"_id": dealID}
	update := bson.M{
		"$inc": bson.M{"fund_required": amountChange},
		"$set": bson.M{"updated_at": time.Now()},
	}
	return s.dealFlowCollection.UpdateOne(context.Background(), filter, update)
}
