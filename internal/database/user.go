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
	UpdateRoles(ctx context.Context, email string, roles []string) (*mongo.UpdateResult, error)
	CreateUser(ctx context.Context, user model.User) (*mongo.InsertOneResult, error)
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
