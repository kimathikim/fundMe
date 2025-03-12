package database

import (
	"context"
	"errors"
	"time"

	"DBackend/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// UserService interface
type UserService interface {
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	// find by id of any model
	FindByID(ctx context.Context, collectionName string, id primitive.ObjectID) (interface{}, error)
	UpdateRoles(ctx context.Context, email string, roles []string) (*mongo.UpdateResult, error)
	CreateUser(ctx context.Context, user model.User) (*mongo.InsertOneResult, error)
	UpdateFounder(ctx context.Context, userID primitive.ObjectID, founder model.Founder) (*mongo.UpdateResult, error)
	UpdateInvestor(ctx context.Context, userID primitive.ObjectID, investor model.Investor) (*mongo.UpdateResult, error)
	CreateRoleData(ctx context.Context, userID primitive.ObjectID, role string) error
	BlacklistToken(ctx context.Context, token string) error
	IsTokenBlacklisted(ctx context.Context, token string) (bool, error)
}

// userService struct
type userService struct {
	userCollection      *mongo.Collection
	founderCollection   *mongo.Collection
	investorCollection  *mongo.Collection
	adminCollection     *mongo.Collection
	blacklistCollection *mongo.Collection
}

// NewUserService initializes collections
func NewUserService(client *mongo.Client) UserService {
	return &userService{
		userCollection:      client.Database("ddb").Collection("users"),
		founderCollection:   client.Database("ddb").Collection("founders"),
		investorCollection:  client.Database("ddb").Collection("investors"),
		adminCollection:     client.Database("ddb").Collection("admins"),
		blacklistCollection: client.Database("ddb").Collection("blacklist_tokens"),
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

// Update founder data
func (s *userService) UpdateFounder(ctx context.Context, userID primitive.ObjectID, founder model.Founder) (*mongo.UpdateResult, error) {
	filter := bson.M{"user_id": userID}
	updateFields := bson.M{
		"startup_name":       founder.StartupName,
		"mission_statement":  founder.MissionStatement,
		"industry":           founder.Industry,
		"funding_stage":      founder.FundingStage,
		"funding_allocation": founder.FundingAllocation,
		"bussiness_model":    founder.BusinessModel,
		"revenue_streams":    founder.RevenueStreams,
		"traction":           founder.Traction,
		"scaling_potential":  founder.ScalingPotential,
		"competition":        founder.Competition,
		"leadership_team":    founder.LeadershipTeam,
		"team_size":          founder.TeamSize,
		"location":           founder.Location,
		"startup_website":    founder.StartupWebsite,
		"pitch_deck":         founder.PitchDeck,
	}
	update := bson.M{"$set": updateFields}
	return s.founderCollection.UpdateOne(ctx, filter, update)
}

func (s *userService) UpdateInvestor(ctx context.Context, userID primitive.ObjectID, investor model.Investor) (*mongo.UpdateResult, error) {
	filter := bson.M{"user_id": userID}
	updateFields := bson.M{
		"investment_portfolio":    investor.InvestmentPortfolio,
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
