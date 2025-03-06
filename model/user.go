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
	UserID       primitive.ObjectID `bson:"user_id"`
	StartupName  string             `bson:"startup_name"`
	FundingRaised float64           `bson:"funding_raised"`
}

type Investor struct {
	UserID              primitive.ObjectID `bson:"user_id"`
	InvestmentPortfolio []string           `bson:"investment_portfolio"`
	TotalInvested       float64            `bson:"total_invested"`
}

// Admin-specific details (if any)
type Admin struct {
	UserID primitive.ObjectID `bson:"user_id"`
}

