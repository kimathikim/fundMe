package database

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"DBackend/model"

	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// DealFlowService interface
type DealFlowService interface {
	Service
	AddStartupToDealFlow(ctx context.Context, deal model.DealFlow) (*mongo.InsertOneResult, error)
	GetDealFlowByID(ctx context.Context, id primitive.ObjectID) (*model.DealFlow, error)
	ListAllDealFlow(ctx context.Context) ([]bson.M, error)
	UpdateDealFlow(ctx context.Context, id primitive.ObjectID, updateFields bson.M) (*mongo.UpdateResult, error)
	DeleteDealFlow(ctx context.Context, id primitive.ObjectID) (*mongo.DeleteResult, error)
	AddMeeting(ctx context.Context, dealID primitive.ObjectID, meeting model.Meeting) (*mongo.UpdateResult, error)
	AddDocument(ctx context.Context, dealID primitive.ObjectID, document model.Document) (*mongo.UpdateResult, error)
	AddTask(ctx context.Context, dealID primitive.ObjectID, task model.Task) (*mongo.UpdateResult, error)
	UpdateTaskCompletionStatus(ctx context.Context, dealID primitive.ObjectID, taskID primitive.ObjectID, completed bool) (*mongo.UpdateResult, error)
	AddNote(ctx context.Context, dealID primitive.ObjectID, note model.Note) error
	UpdateDealStage(ctx context.Context, objID primitive.ObjectID, stage string) (any, error)
}

type dealFlowService struct {
	dealFlowCollection     *mongo.Collection
	notificationCollection *mongo.Collection
	activityCollection     *mongo.Collection
}

// UpdateTaskCompletionStatus implements DealFlowService.
func (d *dealFlowService) UpdateTaskCompletionStatus(ctx context.Context, dealID primitive.ObjectID, taskID primitive.ObjectID, completed bool) (*mongo.UpdateResult, error) {
	panic("unimplemented")
}

// AddDocument implements DealFlowService.
func (d *dealFlowService) AddDocument(ctx context.Context, dealID primitive.ObjectID, document model.Document) (*mongo.UpdateResult, error) {
	panic("unimplemented")
}

// AddMeeting implements DealFlowService.
func (d *dealFlowService) AddMeeting(ctx context.Context, dealID primitive.ObjectID, meeting model.Meeting) (*mongo.UpdateResult, error) {
	meeting.ID = primitive.NewObjectID()
	update := bson.M{"$push": bson.M{"meetings": meeting}}
	result, err := d.dealFlowCollection.UpdateOne(ctx, bson.M{"_id": dealID}, update)

	if err == nil && result.ModifiedCount > 0 {
		// Get the deal to access investor ID
		var deal model.DealFlow
		err = d.dealFlowCollection.FindOne(ctx, bson.M{"_id": dealID}).Decode(&deal)
		if err == nil {
			// Add activity record
			d.addActivity(ctx, deal.InvestorID, "meeting", fmt.Sprintf("Meeting scheduled: %s", meeting.Title))
		}
	}

	return result, err
}

// AddNote implements DealFlowService.
func (d *dealFlowService) AddNote(ctx context.Context, dealID primitive.ObjectID, note model.Note) error {
	update := bson.M{"$push": bson.M{"notes": note}}
	_, err := d.dealFlowCollection.UpdateOne(ctx, bson.M{"_id": dealID}, update)
	return err
}

// addActivity creates an activity record for an investor
func (d *dealFlowService) addActivity(ctx context.Context, investorID primitive.ObjectID, activityType, description string) error {
	activity := model.Activity{
		ID:          primitive.NewObjectID(),
		InvestorID:  investorID,
		Type:        activityType,
		Description: description,
		Date:        time.Now(),
	}

	// Print activity details to console
	fmt.Printf("New activity added: [%s] %s for investor %s\n",
		activityType,
		description,
		investorID.Hex())

	_, err := d.activityCollection.InsertOne(ctx, activity)
	return err
}

// AddStartupToDealFlow implements DealFlowService.
func (d *dealFlowService) AddStartupToDealFlow(ctx context.Context, deal model.DealFlow) (*mongo.InsertOneResult, error) {
	// Check if deal already exists
	var existingDeal model.DealFlow
	err := d.dealFlowCollection.FindOne(ctx, bson.M{
		"investor_id": deal.InvestorID,
		"startup_id":  deal.StartupID,
	}).Decode(&existingDeal)

	if err == nil {
		// Deal already exists
		return nil, errors.New("deal already exists for this startup and investor")
	}

	// Only proceed if error is "no documents found"
	if err != mongo.ErrNoDocuments {
		return nil, err
	}

	deal.CreatedAt = time.Now()
	deal.UpdatedAt = time.Now()
	result, err := d.dealFlowCollection.InsertOne(ctx, deal)
	if err != nil {
		return nil, err
	}

	// Add activity record
	d.addActivity(ctx, deal.InvestorID, "deal", "Added new startup to deal flow")

	// Call notification function
	err = d.notifyFounder(ctx, deal.StartupID, "Your deal has been added by an investor.")
	if err != nil {
		return nil, err
	}
	return result, nil
}

// AddTask implements DealFlowService.
func (d *dealFlowService) AddTask(ctx context.Context, dealID primitive.ObjectID, task model.Task) (*mongo.UpdateResult, error) {
	task.ID = primitive.NewObjectID()
	update := bson.M{"$push": bson.M{"tasks": task}}
	return d.dealFlowCollection.UpdateOne(ctx, bson.M{"_id": dealID}, update)
}

// DealFlow implements DealFlowService.
func (d *dealFlowService) DealFlow() DealFlowService {
	return nil 
}

// DeleteDealFlow implements DealFlowService.
func (d *dealFlowService) DeleteDealFlow(ctx context.Context, id primitive.ObjectID) (*mongo.DeleteResult, error) {
	return d.dealFlowCollection.DeleteOne(ctx, bson.M{"_id": id})
}

// Founder implements DealFlowService.
func (d *dealFlowService) Founder() FounderService {
	// Get the MongoDB client from the existing collection
	client := d.dealFlowCollection.Database().Client()

	// Create and return a new FounderService instance using the same client
	return NewFounderService(client)
}

// GetDealFlowByID implements DealFlowService.
func (d *dealFlowService) GetDealFlowByID(ctx context.Context, id primitive.ObjectID) (*model.DealFlow, error) {
	var deal model.DealFlow
	err := d.dealFlowCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&deal)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("deal not found")
		}
		return nil, err
	}
	return &deal, nil
}

// Health implements DealFlowService.
func (d *dealFlowService) Health() map[string]string {
	return map[string]string{
		"status":  "healthy",
		"service": "dealflow",
	}
}

// Investor implements DealFlowService.
func (d *dealFlowService) Investor() InvestorService {
	// This should return an InvestorService implementation
	// Since we don't have the implementation details, returning nil for now
	return nil
}

// ListAllDealFlow implements DealFlowService.
func (d *dealFlowService) ListAllDealFlow(ctx context.Context) ([]bson.M, error) {
	// Create a pipeline to join founders and users collections
	pipeline := mongo.Pipeline{
		// First lookup to join with founders collection
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "founders"},
			{Key: "localField", Value: "founder_id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "startup"},
		}}},
		// Unwind the startup array (from the lookup)
		{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$startup"},
			{Key: "preserveNullAndEmptyArrays", Value: true},
		}}},
		// Second lookup to join with users collection using the user_id from founders
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "users"},
			{Key: "localField", Value: "startup.user_id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "founder"},
		}}},
		// Unwind the founder array (from the second lookup)
		{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$founder"},
			{Key: "preserveNullAndEmptyArrays", Value: true},
		}}},
		// Project to reshape the output
		{{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 1},
			{Key: "investor_id", Value: 1},
			{Key: "founder_id", Value: 1},
			{Key: "stage", Value: 1},
			{Key: "status", Value: 1},
			{Key: "match_score", Value: 1},
			{Key: "created_at", Value: 1},
			{Key: "updated_at", Value: 1},
			{Key: "last_activity", Value: 1},
			{Key: "startup", Value: 1},
			{Key: "founder_name", Value: bson.D{
				{Key: "$concat", Value: bson.A{"$founder.first_name", " ", "$founder.last_name"}},
			}},
			{Key: "founder_email", Value: "$founder.email"},
		}}},
	}

	// Execute the aggregation pipeline
	cursor, err := d.dealFlowCollection.Aggregate(ctx, pipeline)
	if err != nil {
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

	return results, nil
}

// UpdateDealFlow implements DealFlowService.
func (d *dealFlowService) UpdateDealFlow(ctx context.Context, id primitive.ObjectID, updateFields bson.M) (*mongo.UpdateResult, error) {
	updateFields["updated_at"] = time.Now()
	update := bson.M{"$set": updateFields}
	return d.dealFlowCollection.UpdateOne(ctx, bson.M{"_id": id}, update)
}

// UpdateDealStage implements DealFlowService.
func (d *dealFlowService) UpdateDealStage(ctx context.Context, objID primitive.ObjectID, stage string) (any, error) {
	filter := bson.M{"_id": objID}
	update := bson.M{"$set": bson.M{"stage": stage, "updated_at": time.Now()}}

	result, err := d.dealFlowCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	// Get the deal to access investor ID
	var deal model.DealFlow
	err = d.dealFlowCollection.FindOne(ctx, filter).Decode(&deal)
	if err == nil {
		// Add activity record
		d.addActivity(ctx, deal.InvestorID, "deal_update", fmt.Sprintf("Deal stage updated to %s", stage))
	}

	return result, nil
}

// UpdateDealStatus implements DealFlowService.
func (d *dealFlowService) UpdateDealStatus(ctx context.Context, objID primitive.ObjectID, status string) (*mongo.UpdateResult, error) {
	filter := bson.M{"_id": objID}
	update := bson.M{"$set": bson.M{"status": status, "updated_at": time.Now()}}

	result, err := d.dealFlowCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	// Get the deal to access investor ID
	var deal model.DealFlow
	err = d.dealFlowCollection.FindOne(ctx, filter).Decode(&deal)
	if err == nil {
		// Add activity record
		d.addActivity(ctx, deal.InvestorID, "deal_update", fmt.Sprintf("Deal status updated to %s", status))
	}

	return result, err
}

// UpdateInvestorPortfolio implements DealFlowService.
func (d *dealFlowService) UpdateInvestorPortfolio(ctx context.Context, investorID, startupID primitive.ObjectID, amount float64) (*mongo.UpdateResult, error) {
	filter := bson.M{"user_id": investorID}
	update := bson.M{
		"$inc": bson.M{"total_invested": amount},
		"$push": bson.M{"investment_portfolio": bson.M{
			"startup_id": startupID,
			"amount":     amount,
		}},
	}
	// Since we don't have direct access to the investor collection, we'll need to use the dealFlowCollection
	// This might need adjustment based on your actual data model
	return d.dealFlowCollection.UpdateOne(ctx, filter, update)
}

// UpdateStartupInvestment implements DealFlowService.
func (d *dealFlowService) UpdateStartupInvestment(ctx context.Context, startupID primitive.ObjectID, amount float64) (*mongo.UpdateResult, error) {
	filter := bson.M{"startup_id": startupID}
	update := bson.M{
		"$inc": bson.M{"investment_amount": amount},
		"$set": bson.M{"updated_at": time.Now()},
	}
	return d.dealFlowCollection.UpdateOne(ctx, filter, update)
}

// UpdateTaskStatus implements DealFlowService.
func (d *dealFlowService) UpdateTaskStatus(ctx context.Context, dealID primitive.ObjectID, taskID primitive.ObjectID, completed bool) (*mongo.UpdateResult, error) {
	filter := bson.M{"_id": dealID, "tasks._id": taskID}
	update := bson.M{"$set": bson.M{"tasks.$.completed": completed}}
	return d.dealFlowCollection.UpdateOne(ctx, filter, update)
}

// User implements DealFlowService.
func (d *dealFlowService) User() UserService {
	// This should return a UserService implementation
	// Since we don't have the implementation details, returning nil for now
	return nil
}

// notifyFounder implements DealFlowService.
func (d *dealFlowService) notifyFounder(ctx context.Context, founderID primitive.ObjectID, message string) error {
	notification := model.Notification{
		ID:        primitive.NewObjectID(),
		FounderID: founderID,
		Message:   message,
		CreatedAt: time.Now(),
	}
	_, err := d.notificationCollection.InsertOne(ctx, notification)
	return err
}

func NewDealFlowService(client *mongo.Client) DealFlowService {
	dbName := os.Getenv("BLUEPRINT_DB_DATABASE")
	if dbName == "" {
		dbName = "ddb" // Fallback name
	}

	return &dealFlowService{
		dealFlowCollection:     client.Database(dbName).Collection("deal_flow"),
		notificationCollection: client.Database(dbName).Collection("notifications"),
		activityCollection:     client.Database(dbName).Collection("activities"),
	}
}

// UpdateDealStatus updates the status of a deal
func (s *userService) UpdateDealStatus(ctx *fasthttp.RequestCtx, objID primitive.ObjectID, status string) (any, error) {
	update := bson.M{"$set": bson.M{"status": status, "updated_at": time.Now()}}
	result, err := s.dealFlowCollection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// ListDealsByInvestorID retrieves all deal flow entries for a given investor ID
func (s *userService) ListDealsByInvestorID(ctx context.Context, investorID primitive.ObjectID) ([]bson.M, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.D{{Key: "investor_id", Value: investorID}}}},
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "founders"},
			{Key: "localField", Value: "founder_id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "startup"},
		}}},
		{{Key: "$unwind", Value: "$startup"}},
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "users"},
			{Key: "localField", Value: "startup.user_id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "founder"},
		}}},
		{{Key: "$unwind", Value: "$founder"}},
		{{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 1},
			{Key: "investor_id", Value: 1},
			{Key: "founder_id", Value: 1},
			{Key: "stage", Value: 1},
			{Key: "status", Value: 1},
			{Key: "match_score", Value: 1},
			{Key: "created_at", Value: 1},
			{Key: "updated_at", Value: 1},
			{Key: "last_activity", Value: 1},
			{Key: "startup", Value: 1},
			{Key: "founder_name", Value: bson.D{
				{Key: "$concat", Value: bson.A{"$founder.first_name", " ", "$founder.last_name"}},
			}},
			{Key: "founder_email", Value: "$founder.email"},
			{Key: "founder_avatar", Value: "$founder.avatar"},
		}}},
	}
	cursor, err := s.dealFlowCollection.Aggregate(ctx, pipeline)
	if err != nil {
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
	return results, nil
}

// // Add a startup to deal flow
func (s *userService) AddStartupToDealFlow(ctx context.Context, deal model.DealFlow) (*mongo.InsertOneResult, error) {
	deal.CreatedAt = time.Now()
	deal.UpdatedAt = time.Now()
	result, err := s.dealFlowCollection.InsertOne(ctx, deal)
	if err != nil {
		return nil, err
	}
	// Call notification function
	err = s.notifyFounder(ctx, deal.StartupID, "Your dealjhas been added by an investor.")
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Notify founder about the deal addition
func (s *userService) notifyFounder(ctx context.Context, founderID primitive.ObjectID, message string) error {
	notification := model.Notification{
		ID:        primitive.NewObjectID(),
		FounderID: founderID,
		Message:   message,
		CreatedAt: time.Now(),
	}
	_, err := s.notificationCollection.InsertOne(ctx, notification)
	return err
}

// Get a specific deal flow entry
func (s *userService) GetDealFlowByID(ctx context.Context, id primitive.ObjectID) (*model.DealFlow, error) {
	var deal model.DealFlow
	err := s.dealFlowCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&deal)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("deal not found")
		}
		return nil, err
	}
	return &deal, nil
}

func (s *userService) ListAllDealFlow(ctx context.Context) ([]bson.M, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "founders"},
			{Key: "localField", Value: "founder_id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "startup"},
		}}},
		{{Key: "$unwind", Value: "$startup"}},
	}
	cursor, err := s.dealFlowCollection.Aggregate(ctx, pipeline)
	if err != nil {
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

	return results, nil
}

// Update deal flow details
func (s *userService) UpdateDealFlow(ctx context.Context, id primitive.ObjectID, updateFields bson.M) (*mongo.UpdateResult, error) {
	updateFields["updated_at"] = time.Now()
	update := bson.M{"$set": updateFields}
	return s.dealFlowCollection.UpdateOne(ctx, bson.M{"_id": id}, update)
}

// Remove a startup from deal flow
func (s *userService) DeleteDealFlow(ctx context.Context, id primitive.ObjectID) (*mongo.DeleteResult, error) {
	return s.dealFlowCollection.DeleteOne(ctx, bson.M{"_id": id})
}

// Add a meeting to deal flow
func (s *userService) AddMeeting(ctx context.Context, dealID primitive.ObjectID, meeting model.Meeting) (*mongo.UpdateResult, error) {
	meeting.ID = primitive.NewObjectID()
	update := bson.M{"$push": bson.M{"meetings": meeting}}
	return s.dealFlowCollection.UpdateOne(ctx, bson.M{"_id": dealID}, update)
}

// Add a document to deal flow
func (s *userService) AddDocument(ctx context.Context, dealID primitive.ObjectID, document model.Document) (*mongo.UpdateResult, error) {
	update := bson.M{"$push": bson.M{"documents": document}}
	return s.dealFlowCollection.UpdateOne(ctx, bson.M{"_id": dealID}, update)
}

// Add a task to deal flow
func (s *userService) AddTask(ctx context.Context, dealID primitive.ObjectID, task model.Task) (*mongo.UpdateResult, error) {
	task.ID = primitive.NewObjectID()
	update := bson.M{"$push": bson.M{"tasks": task}}
	return s.dealFlowCollection.UpdateOne(ctx, bson.M{"_id": dealID}, update)
}

// Update task completion status
func (s *userService) UpdateTaskStatus(ctx context.Context, dealID primitive.ObjectID, taskID primitive.ObjectID, completed bool) (*mongo.UpdateResult, error) {
	filter := bson.M{"_id": dealID, "tasks._id": taskID}
	update := bson.M{"$set": bson.M{"tasks.$.completed": completed}}
	return s.dealFlowCollection.UpdateOne(ctx, filter, update)
}

// Update startup's total invested amount
func (s *userService) UpdateStartupInvestment(ctx context.Context, startupID primitive.ObjectID, amount float64) (*mongo.UpdateResult, error) {
	filter := bson.M{"user_id": startupID}
	update := bson.M{"$inc": bson.M{"total_invested": amount}}
	return s.founderCollection.UpdateOne(ctx, filter, update)
}

// Update investor's total investment and add to portfolio
func (s *userService) UpdateInvestorPortfolio(ctx context.Context, investorID, startupID primitive.ObjectID, amount float64) (*mongo.UpdateResult, error) {
	filter := bson.M{"user_id": investorID}
	update := bson.M{
		"$inc": bson.M{"total_invested": amount},
		"$push": bson.M{"investment_portfolio": bson.M{
			"startup_id": startupID,
			"amount":     amount,
		}},
	}
	return s.investorCollection.UpdateOne(ctx, filter, update)
}

// Add a note to a deal flow
func (s *userService) AddNote(ctx context.Context, dealID primitive.ObjectID, note model.Note) error {
	update := bson.M{"$push": bson.M{"notes": note}}
	_, err := s.dealFlowCollection.UpdateOne(ctx, bson.M{"_id": dealID}, update)
	return err
}

// Update deal stage
func (s *userService) UpdateDealStage(ctx context.Context, objID primitive.ObjectID, stage string) (any, error) {
	update := bson.M{"$set": bson.M{"stage": stage, "updated_at": time.Now()}}
	result, err := s.dealFlowCollection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// UpdateDealFundRequired updates the fund required for a deal
func (d *dealFlowService) UpdateDealFundRequired(ctx context.Context, dealID primitive.ObjectID, amountChange float64) (*mongo.UpdateResult, error) {
	filter := bson.M{"_id": dealID}
	update := bson.M{
		"$inc": bson.M{"fund_required": amountChange},
		"$set": bson.M{"updated_at": time.Now()},
	}
	return d.dealFlowCollection.UpdateOne(ctx, filter, update)
}

// Investment implements DealFlowService.
func (d *dealFlowService) Investment() InvestmentService {
	// Get the MongoDB client from the existing collection
	client := d.dealFlowCollection.Database().Client()
	
	// Create and return a new InvestmentService instance using the same database
	dbName := d.dealFlowCollection.Database().Name()
	return NewInvestmentService(client.Database(dbName))
}
