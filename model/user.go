package model

import "go.mongodb.org/mongo-driver/bson/primitive"

// Unified user profile
type User struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	FirstName  string             `bson:"first_name"`
	SecondName string             `bson:"second_name"`
	Email      string             `bson:"email"`
	Password   string             `bson:"password"`
	Roles      []string           `bson:"roles"`
}

// role specific user profiles
type Founder struct {
	UserID            primitive.ObjectID `bson:"user_id"`
	StartupName       string             `bson:"startup_name"`
	MissionStatement  string             `bson:"mission_statement"`
	Industry          string             `bson:"industry"`
	FundingStage      string             `bson:"funding_stage"`
	FundingAllocation string             `bson:"funding_allocation"`
	BusinessModel     string             `bson:"business_model"`
	RevenueStreams    string             `bson:"revenue_streams"`
	Traction          string             `bson:"traction"`
	ScalingPotential  string             `bson:"scaling_potential"`
	Competition       string             `bson:"competition"`
	LeadershipTeam    string             `bson:"leadership_team"`
	TeamSize          int                `bson:"team_size"`
	Location          string             `bson:"location"`
	StartupWebsite    string             `bson:"startup_website"`
	PitchDeck         string             `bson:"pitch_deck"`
}

type Investor struct {
	UserID                primitive.ObjectID `bson:"user_id"`
	InvestmentPortfolio   []string           `bson:"investment_portfolio"`
	TotalInvested         float64            `bson:"total_invested"`
	InvestorType          string             `bson:"investor_type"`
	Thesis                string             `bson:"thesis"`
	PreferredFundingStage string             `bson:"preferred_funding_stage"`
	InvestmentRange       string             `bson:"investment_range"`
	InvestmentFrequency   string             `bson:"investment_frequency"`
	RiskTolerance         string             `bson:"risk_tolerance"`
	ExitStrategy          string             `bson:"exit_strategy"`
	PreferredIndustries   []string           `bson:"preferred_industries"`
	PreferredRegions      []string           `bson:"preferred_regions"`
}

// Admin-specific details (if any)
type Admin struct {
	UserID primitive.ObjectID `bson:"user_id"`
}
