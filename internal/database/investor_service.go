package database

import (
	"DBackend/model"
	"context"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"sort"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// InvestorService defines methods for investor-related operations
type InvestorService interface {
	GetPortfolioSummary(ctx context.Context, investorID primitive.ObjectID) (map[string]interface{}, error)
	GetPipelineSummary(ctx context.Context, investorID primitive.ObjectID) (map[string]interface{}, error)
	GetPerformanceData(ctx context.Context, investorID primitive.ObjectID, period string) ([]map[string]interface{}, error)
	GetPerformanceMetrics(ctx context.Context, investorID primitive.ObjectID, period string) (map[string]interface{}, error)
	GetRecentActivities(ctx context.Context, investorID primitive.ObjectID) ([]map[string]interface{}, error)
	// Other methods...
}

type investorService struct {
	investorCollection *mongo.Collection
	dealFlowCollection *mongo.Collection
	activityCollection *mongo.Collection
	founderCollection  *mongo.Collection
	matchCollection    *mongo.Collection
}

// NewInvestorService initializes the investor service
func NewInvestorService(client *mongo.Client) InvestorService {
	dbName := os.Getenv("BLUEPRINT_DB_DATABASE")
	if dbName == "" {
		dbName = "ddb" // Fallback name
	}

	return &investorService{
		investorCollection: client.Database(dbName).Collection("investors"),
		dealFlowCollection: client.Database(dbName).Collection("deal_flow"),
		activityCollection: client.Database(dbName).Collection("activities"),
		founderCollection:  client.Database(dbName).Collection("founders"),
		matchCollection:    client.Database(dbName).Collection("matches"),
	}
}

// GetPortfolioSummary returns a summary of the investor's portfolio
func (s *investorService) GetPortfolioSummary(ctx context.Context, investorID primitive.ObjectID) (map[string]interface{}, error) {
	// Query the investor collection to get the investor's portfolio
	filter := bson.M{"user_id": investorID}
	var investorData bson.M
	err := s.investorCollection.FindOne(ctx, filter).Decode(&investorData)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Return empty data if investor not found
			return map[string]interface{}{
				"totalInvested": 0,
				"totalStartups": 0,
				"avgReturn":     0,
				"topPerformers": []map[string]interface{}{},
			}, nil
		}
		return nil, err
	}

	totalInvested := 0
	totalStartups := 0
	avgReturn := 0.0
	var topPerformers []map[string]interface{}
	var performanceData []struct {
		startup      bson.M
		roi          float64
		invested     float64
		currentValue float64
	}

	// Extract total invested amount
	if val, ok := investorData["total_invested"]; ok {
		if totalInv, ok := val.(int32); ok {
			totalInvested = int(totalInv)
		} else if totalInv, ok := val.(int64); ok {
			totalInvested = int(totalInv)
		} else if totalInv, ok := val.(float64); ok {
			totalInvested = int(totalInv)
		}
	}

	// Calculate number of startups in portfolio and gather performance data
	if portfolio, ok := investorData["investment_portfolio"].(primitive.A); ok {
		totalStartups = len(portfolio)

		// Calculate total ROI for average return
		totalROI := 0.0
		validInvestments := 0

		for _, inv := range portfolio {
			if investment, ok := inv.(bson.M); ok {
				startupID, ok := investment["startup_id"].(primitive.ObjectID)
				if !ok {
					continue
				}

				// Get startup details
				var startupData bson.M
				err := s.founderCollection.FindOne(ctx, bson.M{"user_id": startupID}).Decode(&startupData)
				if err != nil {
					continue
				}

				// Get investment details
				invested := 0.0
				if amount, ok := investment["amount"].(float64); ok {
					invested = amount
				} else if amount, ok := investment["amount"].(int32); ok {
					invested = float64(amount)
				} else if amount, ok := investment["amount"].(int64); ok {
					invested = float64(amount)
				}

				// Get current value (could be from startup valuation or other metrics)
				currentValue := invested // Default to no change
				if val, ok := investment["current_value"].(float64); ok {
					currentValue = val
				} else {
					// If current_value not set, estimate based on startup growth
					growth := 0.0
					if g, ok := startupData["growth_rate"].(float64); ok {
						growth = g
					}
					currentValue = invested * (1 + growth/100)

					// Update the current value in the database
					updateFilter := bson.M{
						"user_id":                         investorID,
						"investment_portfolio.startup_id": startupID,
					}
					update := bson.M{
						"$set": bson.M{"investment_portfolio.$.current_value": currentValue},
					}
					s.investorCollection.UpdateOne(ctx, updateFilter, update)
				}

				// Calculate ROI
				roi := 0.0
				if invested > 0 {
					roi = ((currentValue - invested) / invested) * 100
					totalROI += roi
					validInvestments++
				}

				// Add to performance data
				performanceData = append(performanceData, struct {
					startup      bson.M
					roi          float64
					invested     float64
					currentValue float64
				}{
					startup:      startupData,
					roi:          roi,
					invested:     invested,
					currentValue: currentValue,
				})
			}
		}

		// Calculate average return
		if validInvestments > 0 {
			avgReturn = totalROI / float64(validInvestments)
		}
	}

	// Sort startups by ROI (descending)
	sort.Slice(performanceData, func(i, j int) bool {
		return performanceData[i].roi > performanceData[j].roi
	})

	// Take top 3 performers (or fewer if portfolio is smaller)
	maxTopPerformers := 3
	if len(performanceData) < maxTopPerformers {
		maxTopPerformers = len(performanceData)
	}

	for i := 0; i < maxTopPerformers; i++ {
		data := performanceData[i]

		// Extract relevant startup data
		startupName := ""
		if name, ok := data.startup["startup_name"].(string); ok {
			startupName = name
		}

		industry := ""
		if ind, ok := data.startup["industry"].(string); ok {
			industry = ind
		}

		performer := map[string]interface{}{
			"id":           data.startup["_id"],
			"name":         startupName,
			"industry":     industry,
			"growth":       data.roi,
			"invested":     data.invested,
			"currentValue": data.currentValue,
			"roi":          data.roi,
		}
		topPerformers = append(topPerformers, performer)
	}

	return map[string]interface{}{
		"totalInvested": totalInvested,
		"totalStartups": totalStartups,
		"avgReturn":     avgReturn,
		"topPerformers": topPerformers,
	}, nil
}

// GetPipelineSummary returns a summary of the investor's pipeline
func (s *investorService) GetPipelineSummary(ctx context.Context, investorID primitive.ObjectID) (map[string]interface{}, error) {
	// Query the deal flow collection to get the investor's pipeline
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.D{{Key: "investor_id", Value: investorID}}}},
		{{Key: "$facet", Value: bson.D{
			{Key: "totalDeals", Value: bson.A{
				bson.D{{Key: "$count", Value: "count"}},
			}},
			{Key: "byStatus", Value: bson.A{
				bson.D{{Key: "$group", Value: bson.D{
					{Key: "_id", Value: "$status"},
					{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
				}}},
			}},
			{Key: "byStage", Value: bson.A{
				bson.D{{Key: "$group", Value: bson.D{
					{Key: "_id", Value: "$stage"},
					{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
				}}},
			}},
			{Key: "recentDeals", Value: bson.A{
				bson.D{{Key: "$sort", Value: bson.D{{Key: "updated_at", Value: -1}}}},
				bson.D{{Key: "$limit", Value: 5}},
				bson.D{{Key: "$lookup", Value: bson.D{
					{Key: "from", Value: "founders"},
					{Key: "localField", Value: "founder_id"},
					{Key: "foreignField", Value: "user_id"},
					{Key: "as", Value: "founder"},
				}}},
				bson.D{{Key: "$project", Value: bson.D{
					{Key: "_id", Value: 1},
					{Key: "founder_id", Value: 1},
					{Key: "status", Value: 1},
					{Key: "stage", Value: 1},
					{Key: "match_percentage", Value: 1},
					{Key: "updated_at", Value: 1},
					{Key: "startup_name", Value: bson.D{{Key: "$arrayElemAt", Value: bson.A{"$founder.startup_name", 0}}}},
					{Key: "industry", Value: bson.D{{Key: "$arrayElemAt", Value: bson.A{"$founder.industry", 0}}}},
				}}},
			}},
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

	totalDeals := 0
	pendingDeals := 0
	closedDeals := 0
	var recentDeals []map[string]interface{}

	if len(results) > 0 {
		result := results[0]

		// Extract total deals count
		if totalDealsArr, ok := result["totalDeals"].(primitive.A); ok && len(totalDealsArr) > 0 {
			if countDoc, ok := totalDealsArr[0].(bson.M); ok {
				if count, ok := countDoc["count"].(int32); ok {
					totalDeals = int(count)
				}
			}
		}

		// Extract deals by stage
		if byStage, ok := result["byStage"].(primitive.A); ok {
			for _, stageGroup := range byStage {
				if group, ok := stageGroup.(bson.M); ok {
					stage, hasStage := group["_id"].(string)
					count, hasCount := group["count"].(int32)

					if hasStage && hasCount {
						if stage == "closed_won" || stage == "closed_lost" {
							closedDeals += int(count)
						} else {
							pendingDeals += int(count)
						}
					}
				}
			}
		}

		// Extract recent deals
		if recentDealsArr, ok := result["recentDeals"].(primitive.A); ok {
			for _, deal := range recentDealsArr {
				if dealDoc, ok := deal.(bson.M); ok {
					stage := dealDoc["stage"].(string)
					status := "pending"
					if stage == "Closed Won" || stage == "Closed Lost" {
						status = "closed"
					}

					recentDeal := map[string]interface{}{
						"id":               dealDoc["_id"],
						"founder_id":       dealDoc["founder_id"],
						"startup_name":     dealDoc["startup_name"],
						"industry":         dealDoc["industry"],
						"status":           status,
						"stage":            stage,
						"match_percentage": dealDoc["match_percentage"],
						"updated_at":       dealDoc["updated_at"],
					}
					recentDeals = append(recentDeals, recentDeal)
				}
			}
		}
	}

	return map[string]interface{}{
		"totalDeals":   totalDeals,
		"pendingDeals": pendingDeals,
		"closedDeals":  closedDeals,
		"recentDeals":  recentDeals,
	}, nil
}

// GetPerformanceData returns performance data for the specified period
func (s *investorService) GetPerformanceData(ctx context.Context, investorID primitive.ObjectID, period string) ([]map[string]interface{}, error) {
	// Determine date range based on period
	endDate := time.Now()
	var startDate time.Time

	switch period {
	case "1m":
		startDate = endDate.AddDate(0, -1, 0)
	case "3m":
		startDate = endDate.AddDate(0, -3, 0)
	case "6m":
		startDate = endDate.AddDate(0, -6, 0)
	case "1y":
		startDate = endDate.AddDate(-1, 0, 0)
	case "3y":
		startDate = endDate.AddDate(-3, 0, 0)
	case "5y":
		startDate = endDate.AddDate(-5, 0, 0)
	default:
		startDate = endDate.AddDate(-1, 0, 0) // Default to 1 year
	}

	// Query the investor's portfolio and valuation history
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.D{{Key: "_id", Value: investorID}}}},
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "portfolio_valuations"},
			{Key: "localField", Value: "_id"},
			{Key: "foreignField", Value: "investor_id"},
			{Key: "as", Value: "valuations"},
		}}},
		{{Key: "$project", Value: bson.D{
			{Key: "valuations", Value: bson.D{
				{Key: "$filter", Value: bson.D{
					{Key: "input", Value: "$valuations"},
					{Key: "as", Value: "valuation"},
					{Key: "cond", Value: bson.D{
						{Key: "$and", Value: bson.A{
							bson.D{{Key: "$gte", Value: bson.A{"$$valuation.date", startDate}}},
							bson.D{{Key: "$lte", Value: bson.A{"$$valuation.date", endDate}}},
						}},
					}},
				}},
			}},
		}}},
	}

	cursor, err := s.investorCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	// Process the performance data
	var performanceData []map[string]interface{}

	if len(results) > 0 && len(results[0]) > 0 {
		if valuations, ok := results[0]["valuations"].(primitive.A); ok {
			// Sort valuations by date
			type valuation struct {
				Date  time.Time
				Value float64
			}

			var valuationData []valuation
			for _, v := range valuations {
				if val, ok := v.(bson.M); ok {
					date, hasDate := val["date"].(time.Time)
					value, hasValue := val["value"].(float64)

					if hasDate && hasValue {
						valuationData = append(valuationData, valuation{
							Date:  date,
							Value: value,
						})
					}
				}
			}

			// Sort by date
			sort.Slice(valuationData, func(i, j int) bool {
				return valuationData[i].Date.Before(valuationData[j].Date)
			})

			// Convert to the required format
			for _, v := range valuationData {
				performanceData = append(performanceData, map[string]interface{}{
					"date":  v.Date,
					"value": v.Value,
				})
			}

			// If no data points exist, create synthetic data for visualization
			if len(performanceData) == 0 {
				// Generate monthly data points for the period
				currentDate := startDate
				initialValue := 100000.0 // Example initial portfolio value

				for currentDate.Before(endDate) {
					performanceData = append(performanceData, map[string]interface{}{
						"date":  currentDate,
						"value": initialValue,
					})

					// Move to next month and add some random variation
					currentDate = currentDate.AddDate(0, 1, 0)
					// This is just a placeholder - in a real implementation you'd use actual data
					initialValue = initialValue * (1 + (rand.Float64()*0.05 - 0.01))
				}
			}
		}
	}

	return performanceData, nil
}

// GetPerformanceMetrics returns performance metrics for the specified period
func (s *investorService) GetPerformanceMetrics(ctx context.Context, investorID primitive.ObjectID, period string) (map[string]interface{}, error) {
	// Get performance data for the period
	performanceData, err := s.GetPerformanceData(ctx, investorID, period)
	if err != nil {
		return nil, err
	}

	// Default values
	totalReturn := 0.0
	annualizedReturn := 0.0
	volatility := 0.0

	// Calculate metrics if we have data
	if len(performanceData) > 1 {
		// Get initial and final values
		initialValue := performanceData[0]["value"].(float64)
		finalValue := performanceData[len(performanceData)-1]["value"].(float64)

		// Calculate total return
		totalReturn = ((finalValue - initialValue) / initialValue) * 100

		// Calculate time period in years
		startDate := performanceData[0]["date"].(time.Time)
		endDate := performanceData[len(performanceData)-1]["date"].(time.Time)
		yearsDiff := endDate.Sub(startDate).Hours() / 24 / 365

		// Calculate annualized return
		if yearsDiff > 0 {
			annualizedReturn = (math.Pow(1+(totalReturn/100), 1/yearsDiff) - 1) * 100
		}

		// Calculate volatility (standard deviation of returns)
		if len(performanceData) > 2 {
			var returns []float64
			var prevValue float64

			for i, data := range performanceData {
				currentValue := data["value"].(float64)

				if i > 0 {
					periodReturn := (currentValue - prevValue) / prevValue
					returns = append(returns, periodReturn)
				}

				prevValue = currentValue
			}

			// Calculate standard deviation
			var sum float64
			var mean float64

			// Calculate mean
			for _, r := range returns {
				sum += r
			}
			mean = sum / float64(len(returns))

			// Calculate variance
			var variance float64
			for _, r := range returns {
				variance += math.Pow(r-mean, 2)
			}
			variance = variance / float64(len(returns))

			// Standard deviation
			volatility = math.Sqrt(variance) * 100
		}
	}

	return map[string]interface{}{
		"totalReturn":      totalReturn,
		"annualizedReturn": annualizedReturn,
		"volatility":       volatility,
	}, nil
}

// GetRecentActivities returns recent activities for an investor
func (s *investorService) GetRecentActivities(ctx context.Context, investorID primitive.ObjectID) ([]map[string]interface{}, error) {
	// Query the activities collection to get recent investor activities
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.D{{Key: "investor_id", Value: investorID}}}},
		{{Key: "$sort", Value: bson.D{{Key: "date", Value: -1}}}},
		{{Key: "$limit", Value: 5}},
		{{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 1},
			{Key: "type", Value: 1},
			{Key: "description", Value: 1},
			{Key: "date", Value: 1},
		}}},
	}

	cursor, err := s.activityCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var activities []map[string]interface{}
	if err = cursor.All(ctx, &activities); err != nil {
		return nil, err
	}

	// If no activities found, return empty array
	if activities == nil {
		activities = []map[string]interface{}{}
	}

	return activities, nil
}

// GetMatches returns matches for an investor
func (s *investorService) GetMatches(ctx context.Context, investorID primitive.ObjectID) ([]map[string]interface{}, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.D{{Key: "investor_id", Value: investorID}}}},
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "founders"},
			{Key: "localField", Value: "founder_id"},
			{Key: "foreignField", Value: "user_id"},
			{Key: "as", Value: "founder"},
		}}},
		{{Key: "$unwind", Value: "$founder"}},
		{{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 1},
			{Key: "founder_id", Value: 1},
			{Key: "match_percentage", Value: 1},
			{Key: "created_at", Value: 1},
			{Key: "startup_name", Value: "$founder.startup_name"},
			{Key: "industry", Value: "$founder.industry"},
			{Key: "funding_stage", Value: "$founder.funding_stage"},
			{Key: "fund_required", Value: "$founder.fund_required"},
		}}},
		{{Key: "$sort", Value: bson.D{{Key: "match_percentage", Value: -1}}}},
	}

	cursor, err := s.matchCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var matches []map[string]interface{}
	if err = cursor.All(ctx, &matches); err != nil {
		return nil, err
	}

	return matches, nil
}

// StoreMatch stores a match between investor and founder
func (s *investorService) StoreMatch(ctx context.Context, match model.MatchInvestorFounder) (*mongo.InsertOneResult, error) {
	// Check if match already exists
	filter := bson.M{
		"investor_id": match.InvestorID,
		"founder_id":  match.FounderID,
	}

	var existingMatch model.MatchInvestorFounder
	err := s.matchCollection.FindOne(ctx, filter).Decode(&existingMatch)
	if err == nil {
		// Match exists, update it
		update := bson.M{
			"$set": bson.M{
				"match_percentage": match.MatchPercentage,
				"updated_at":       time.Now(),
			},
		}
		_, err := s.matchCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			return nil, err
		}

		// Return the existing ID
		return &mongo.InsertOneResult{
			InsertedID: existingMatch.ID,
		}, nil
	}

	// Match doesn't exist, insert new one
	match.CreatedAt = time.Now()
	match.UpdatedAt = time.Now()

	result, err := s.matchCollection.InsertOne(ctx, match)
	if err != nil {
		return nil, err
	}

	// Create activity record for the investor
	activity := model.Activity{
		ID:          primitive.NewObjectID(),
		InvestorID:  match.InvestorID,
		Type:        "match",
		Description: fmt.Sprintf("New match with %.0f%% compatibility", match.MatchPercentage),
		Date:        time.Now(),
	}

	_, err = s.activityCollection.InsertOne(ctx, activity)
	if err != nil {
		// Log error but don't fail the match creation
		log.Printf("Failed to create activity record: %v", err)
	}

	return result, nil
}
