package database

import (
	"encoding/json"
	"fmt"
	"os"
	"go.mongodb.org/mongo-driver/bson"
)

// Match and store data in JSON
func GenerateMatchmakingData(founders []bson.M, investors []bson.M) {
	var matches []map[string]interface{}

	for _, founder := range founders {
		for _, investor := range investors {
			// Match if industry & funding stage align
			if contains(investor["preferred_industries"].(bson.A), founder["industry"].(string)) &&
				founder["funding_stage"] == investor["preferred_funding_stage"] {

				match := map[string]interface{}{
					"founder_id":       founder["user_id"],
					"investor_id":      investor["user_id"],
					"fund_required":    founder["fund_required"],
					"total_invested":   investor["total_invested"],
					"industry":         founder["industry"],
					"funding_stage":    founder["funding_stage"],
					"risk_tolerance":   investor["risk_tolerance"],
					"match_percentage": 80, // Placeholder, replace with ML model later
				}
				matches = append(matches, match)
			}
		}
	}

	// Save to JSON
	file, _ := json.MarshalIndent(matches, "", "  ")
	_ = os.WriteFile("matchmaking_data.json", file, 0o644)
	fmt.Println("✅ Matchmaking data saved as matchmaking_data.json")
}

// Helper function to check if an array contains a value
func contains(array bson.A, value string) bool {
	for _, v := range array {
		if v == value {
			return true
		}
	}
	return false
}
