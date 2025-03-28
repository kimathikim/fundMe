package database

import (
    "context"
    "time"

    "DBackend/model"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
)

// InvestmentService defines methods for investment operations
type InvestmentService interface {
    CreateInvestment(ctx context.Context, investment model.Investment) (*mongo.InsertOneResult, error)
    GetInvestmentsByInvestorID(ctx context.Context, investorID primitive.ObjectID) ([]model.Investment, error)
    GetInvestmentsByFounderID(ctx context.Context, founderID primitive.ObjectID) ([]model.Investment, error)
    GetInvestmentsByDealID(ctx context.Context, dealID primitive.ObjectID) ([]model.Investment, error)
}

type investmentService struct {
    investmentCollection *mongo.Collection
}

// NewInvestmentService creates a new investment service
func NewInvestmentService(db *mongo.Database) InvestmentService {
    return &investmentService{
        investmentCollection: db.Collection("investments"),
    }
}

// CreateInvestment creates a new investment record
func (s *investmentService) CreateInvestment(ctx context.Context, investment model.Investment) (*mongo.InsertOneResult, error) {
    investment.CreatedAt = time.Now()
    investment.UpdatedAt = time.Now()
    return s.investmentCollection.InsertOne(ctx, investment)
}

// GetInvestmentsByInvestorID retrieves all investments made by an investor
func (s *investmentService) GetInvestmentsByInvestorID(ctx context.Context, investorID primitive.ObjectID) ([]model.Investment, error) {
    cursor, err := s.investmentCollection.Find(ctx, bson.M{"investor_id": investorID})
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

    var investments []model.Investment
    if err = cursor.All(ctx, &investments); err != nil {
        return nil, err
    }
    return investments, nil
}

// GetInvestmentsByFounderID retrieves all investments received by a founder
func (s *investmentService) GetInvestmentsByFounderID(ctx context.Context, founderID primitive.ObjectID) ([]model.Investment, error) {
    cursor, err := s.investmentCollection.Find(ctx, bson.M{"founder_id": founderID})
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

    var investments []model.Investment
    if err = cursor.All(ctx, &investments); err != nil {
        return nil, err
    }
    return investments, nil
}

// GetInvestmentsByDealID retrieves all investments for a specific deal
func (s *investmentService) GetInvestmentsByDealID(ctx context.Context, dealID primitive.ObjectID) ([]model.Investment, error) {
    cursor, err := s.investmentCollection.Find(ctx, bson.M{"deal_id": dealID})
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

    var investments []model.Investment
    if err = cursor.All(ctx, &investments); err != nil {
        return nil, err
    }
    return investments, nil
}