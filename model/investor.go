package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// InvestorApplication represents an application from a founder to an investor
type InvestorApplication struct {
    ID            primitive.ObjectID `bson:"_id,omitempty"`
    FounderID     primitive.ObjectID `bson:"founder_id"`
    InvestorID    primitive.ObjectID `bson:"investor_id"`
    FundingAmount float64            `bson:"funding_amount"`
    UseOfFunds    string             `bson:"use_of_funds"`
    Status        string             `bson:"status,omitempty"` // Pending, Approved, Rejected
    CreatedAt     time.Time          `bson:"created_at,omitempty"`
    UpdatedAt     time.Time          `bson:"updated_at,omitempty"`
}
