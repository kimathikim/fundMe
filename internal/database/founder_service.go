package database

import (
	"DBackend/model"
	"context"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// FounderService defines the methods for founder-related database operations
type FounderService interface {
	GetFundraisingSummary(ctx context.Context, founderID primitive.ObjectID) (map[string]interface{}, error)
	GetInvestorEngagement(ctx context.Context, founderID primitive.ObjectID) (map[string]interface{}, error)
	GetGrants(ctx context.Context, category, region string) ([]model.Grant, error)
	SubmitGrantApplication(ctx context.Context, application model.GrantApplication) (primitive.ObjectID, error)
	GetInvestors(ctx context.Context, industry, stage string) ([]model.Investor, error)
	SubmitInvestorApplication(ctx context.Context, application model.InvestorApplication) (primitive.ObjectID, error)
	CreateGrant(ctx context.Context, grant model.Grant) (primitive.ObjectID, error)
	GetGrantByID(ctx context.Context, id primitive.ObjectID) (model.Grant, error)
	UpdateGrant(ctx context.Context, id primitive.ObjectID, updates model.Grant) error
	DeleteGrant(ctx context.Context, id primitive.ObjectID) error
	GetAllGrantApplications(ctx context.Context) ([]model.GrantApplication, error)
	GetFounderGrantApplications(ctx context.Context, founderID primitive.ObjectID) ([]model.GrantApplication, error)
	GetGrantApplicationByID(ctx context.Context, id primitive.ObjectID) (model.GrantApplication, error)
	UpdateGrantApplication(ctx context.Context, id primitive.ObjectID, status string, remarks string) error
}

func NewFounderService(client *mongo.Client) FounderService {
	dbName := os.Getenv("BLUEPRINT_DB_DATABASE")
	if dbName == "" {
		dbName = "ddb" // Fallback name
	}

	return &founderService{
		founderCollection:     client.Database(dbName).Collection("founders"),
		grantCollection:       client.Database(dbName).Collection("grants"),
		investorCollection:    client.Database(dbName).Collection("investors"),
		applicationCollection: client.Database(dbName).Collection("applications"),
	}
}

type founderService struct {
	founderCollection *mongo.Collection
	grantCollection *mongo.Collection
	investorCollection *mongo.Collection
	applicationCollection *mongo.Collection
}

// GetFundraisingSummary returns a summary of the founder's fundraising activities
func (s *founderService) GetFundraisingSummary(ctx context.Context, founderID primitive.ObjectID) (map[string]interface{}, error) {
	// Get founder details
	founder, err := s.GetFounderByUserID(ctx, founderID)
	if err != nil {
		return nil, err
	}
	
	// Get investment data
	investments, err := s.GetFounderInvestments(ctx, founderID)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}
	
	// Calculate metrics with safe defaults
	totalRaised := 0.0
	numberOfInvestors := 0
	averageInvestment := 0.0
	
	if len(investments) > 0 {
		// Calculate total raised and other metrics
		for _, inv := range investments {
			totalRaised += inv.Amount
		}
		numberOfInvestors = len(investments)
		if numberOfInvestors > 0 {
			averageInvestment = totalRaised / float64(numberOfInvestors)
		}
	}
	
	// Ensure fundingGoal is never nil
	fundingGoal := founder.FundRequired
	if fundingGoal <= 0 {
		fundingGoal = 1 // Prevent division by zero
	}
	
	// Calculate percentage with bounds checking
	percentageComplete := 0
	if fundingGoal > 0 {
		percentageComplete = int((totalRaised * 100) / float64(fundingGoal))
	}
	
	return map[string]interface{}{
		"totalRaised":        totalRaised,
		"fundingGoal":        fundingGoal,
		"percentageComplete": percentageComplete,
		"numberOfInvestors":  numberOfInvestors,
		"averageInvestment":  averageInvestment,
	}, nil
}

// GetInvestorEngagement returns metrics about investor engagement with the founder
func (s *founderService) GetInvestorEngagement(ctx context.Context, founderID primitive.ObjectID) (map[string]interface{}, error) {
	// Get matches data
	matches, err := s.GetFounderMatches(ctx, founderID)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}
	
	// Initialize with safe defaults
	totalMatches := 0
	newThisMonth := 0
	inDueDiligence := 0
	topMatches := []map[string]interface{}{}
	
	// Current month for filtering
	currentMonth := time.Now().Month()
	
	if matches != nil {
		totalMatches = len(matches)
		
		// Process matches
		for _, match := range matches {
			// Count new matches this month
			if match.CreatedAt.Month() == currentMonth {
				newThisMonth++
			}
			
			// Count deals in due diligence
			if match.Stage == "Due Diligence" {
				inDueDiligence++
			}
			
			// Add to top matches if score is high enough
			if match.MatchScore > 70 {
				topMatch := map[string]interface{}{
					"investorId":           match.InvestorID.Hex(),
					"name":                 match.InvestorName,
					"matchPercentage":      match.MatchScore,
					"industry":             match.Industry,
					"totalInvested":        match.TotalInvested,
					"preferredFundingStage": match.PreferredFundingStage,
				}
				topMatches = append(topMatches, topMatch)
			}
		}
	}
	
	// Ensure we have at least an empty array, not nil
	if topMatches == nil {
		topMatches = []map[string]interface{}{}
	}
	
	return map[string]interface{}{
		"totalMatches":   totalMatches,
		"newThisMonth":   newThisMonth,
		"inDueDiligence": inDueDiligence,
		"topMatches":     topMatches,
	}, nil
}

// GetGrants returns available grants based on category and region
func (s *founderService) GetGrants(ctx context.Context, category, region string) ([]model.Grant, error) {
	// Implement grants retrieval logic
	return []model.Grant{}, nil
}

// SubmitGrantApplication submits a grant application
func (s *founderService) SubmitGrantApplication(ctx context.Context, application model.GrantApplication) (primitive.ObjectID, error) {
	// Implement grant application submission logic
	result, err := s.applicationCollection.InsertOne(ctx, application)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return result.InsertedID.(primitive.ObjectID), nil
}

// GetInvestors returns available investors based on industry and stage
func (s *founderService) GetInvestors(ctx context.Context, industry, stage string) ([]model.Investor, error) {
	// Implement investors retrieval logic
	return []model.Investor{}, nil
}

// SubmitInvestorApplication submits an investor application
func (s *founderService) SubmitInvestorApplication(ctx context.Context, application model.InvestorApplication) (primitive.ObjectID, error) {
	// Implement investor application submission logic
	result, err := s.applicationCollection.InsertOne(ctx, application)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return result.InsertedID.(primitive.ObjectID), nil
}

// CreateGrant creates a new grant
func (s *founderService) CreateGrant(ctx context.Context, grant model.Grant) (primitive.ObjectID, error) {
	result, err := s.grantCollection.InsertOne(ctx, grant)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return result.InsertedID.(primitive.ObjectID), nil
}

// GetGrantByID retrieves a grant by its ID
func (s *founderService) GetGrantByID(ctx context.Context, id primitive.ObjectID) (model.Grant, error) {
	var grant model.Grant
	err := s.grantCollection.FindOne(ctx, primitive.M{"_id": id}).Decode(&grant)
	if err != nil {
		return model.Grant{}, err
	}
	return grant, nil
}

// UpdateGrant updates an existing grant
func (s *founderService) UpdateGrant(ctx context.Context, id primitive.ObjectID, updates model.Grant) error {
	_, err := s.grantCollection.UpdateOne(
		ctx,
		primitive.M{"_id": id},
		primitive.M{"$set": updates},
	)
	return err
}

// DeleteGrant deletes a grant
func (s *founderService) DeleteGrant(ctx context.Context, id primitive.ObjectID) error {
	_, err := s.grantCollection.DeleteOne(ctx, primitive.M{"_id": id})
	return err
}

// GetAllGrantApplications retrieves all grant applications
func (s *founderService) GetAllGrantApplications(ctx context.Context) ([]model.GrantApplication, error) {
	var applications []model.GrantApplication
	cursor, err := s.applicationCollection.Find(ctx, primitive.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	err = cursor.All(ctx, &applications)
	if err != nil {
		return nil, err
	}
	return applications, nil
}

// GetFounderGrantApplications retrieves grant applications for a specific founder
func (s *founderService) GetFounderGrantApplications(ctx context.Context, founderID primitive.ObjectID) ([]model.GrantApplication, error) {
	var applications []model.GrantApplication
	cursor, err := s.applicationCollection.Find(ctx, primitive.M{"founder_id": founderID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	err = cursor.All(ctx, &applications)
	if err != nil {
		return nil, err
	}
	return applications, nil
}

// GetGrantApplicationByID retrieves a specific grant application
func (s *founderService) GetGrantApplicationByID(ctx context.Context, id primitive.ObjectID) (model.GrantApplication, error) {
	var application model.GrantApplication
	err := s.applicationCollection.FindOne(ctx, primitive.M{"_id": id}).Decode(&application)
	if err != nil {
		return model.GrantApplication{}, err
	}
	return application, nil
}

// UpdateGrantApplication updates the status and remarks of a grant application
func (s *founderService) UpdateGrantApplication(ctx context.Context, id primitive.ObjectID, status string, remarks string) error {
	_, err := s.applicationCollection.UpdateOne(
		ctx,
		primitive.M{"_id": id},
		primitive.M{"$set": primitive.M{
			"status":  status,
			"remarks": remarks,
		}},
	)
	return err
}

// GetFounderMatches retrieves matches for a specific founder
func (s *founderService) GetFounderMatches(ctx context.Context, founderID primitive.ObjectID) ([]model.Match, error) {
	var matches []model.Match
	cursor, err := s.applicationCollection.Find(ctx, primitive.M{"founder_id": founderID, "type": "match"})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	err = cursor.All(ctx, &matches)
	if err != nil {
		return nil, err
	}
	return matches, nil
}

// GetFounderByID retrieves a founder by their ID
func (s *founderService) GetFounderByID(ctx context.Context, founderID primitive.ObjectID) (model.Founder, error) {
	var founder model.Founder
	err := s.founderCollection.FindOne(ctx, primitive.M{"_id": founderID}).Decode(&founder)
	if err != nil {
		return model.Founder{}, err
	}
	return founder, nil
}

// GetFounderByUserID retrieves founder details by user ID
func (s *founderService) GetFounderByUserID(ctx context.Context, userID primitive.ObjectID) (*model.Founder, error) {
	var founder model.Founder
	err := s.founderCollection.FindOne(ctx, primitive.M{"user_id": userID}).Decode(&founder)
	if err != nil {
		return nil, err
	}
	return &founder, nil
}

// GetFounderInvestments retrieves all investments made to a founder's startup
func (s *founderService) GetFounderInvestments(ctx context.Context, founderID primitive.ObjectID) ([]model.Investment, error) {
	var investments []model.Investment
	cursor, err := s.applicationCollection.Find(ctx, primitive.M{
		"founder_id": founderID,
		"type": "investment",
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	err = cursor.All(ctx, &investments)
	if err != nil {
		return nil, err
	}
	return investments, nil
}

// // GetFounderMatches retrieves all investor matches for a founder
// func (s *founderService) GetFounderMatches(ctx context.Context, founderID primitive.ObjectID) ([]model.Match, error) {
// 	var matches []model.Match
// 	cursor, err := s.applicationCollection.Find(ctx, primitive.M{
// 		"founder_id": founderID,
// 		"type": "match",
// 	})
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer cursor.Close(ctx)
	
// 	err = cursor.All(ctx, &matches)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return matches, nil
// }


