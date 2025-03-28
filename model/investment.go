package model

import (
    "time"
    
    "go.mongodb.org/mongo-driver/bson/primitive"
)
// Investment represents an investment made by an investor in a startup
type Investment struct {
    ID             primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
    DealID         primitive.ObjectID `bson:"deal_id" json:"dealId"`
    InvestorID     primitive.ObjectID `bson:"investor_id" json:"investorId"`
    FounderID      primitive.ObjectID `bson:"founder_id" json:"founderId"`
    Amount         float64            `bson:"amount" json:"amount"`
    InvestmentDate time.Time          `bson:"investment_date" json:"investmentDate"`
    CreatedAt      time.Time          `bson:"created_at" json:"createdAt"`
    UpdatedAt      time.Time          `bson:"updated_at" json:"updatedAt"`
}