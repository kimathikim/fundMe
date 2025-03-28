package model

import (
    "time"
    
    "go.mongodb.org/mongo-driver/bson/primitive"
)

// PortfolioValuation represents a point-in-time valuation of an investor's portfolio
type PortfolioValuation struct {
    ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    InvestorID primitive.ObjectID `bson:"investor_id" json:"investor_id"`
    Date       time.Time          `bson:"date" json:"date"`
    Value      float64            `bson:"value" json:"value"`
    Assets     []PortfolioAsset   `bson:"assets" json:"assets"`
}

// PortfolioAsset represents a single asset in an investor's portfolio
type PortfolioAsset struct {
    StartupID    primitive.ObjectID `bson:"startup_id" json:"startup_id"`
    StartupName  string             `bson:"startup_name" json:"startup_name"`
    Value        float64            `bson:"value" json:"value"`
    Equity       float64            `bson:"equity" json:"equity"`
    AcquisitionDate time.Time       `bson:"acquisition_date" json:"acquisition_date"`
    LastValuationDate time.Time     `bson:"last_valuation_date" json:"last_valuation_date"`
}

// PipelineSummary represents a summary of an investor's deal pipeline
type PipelineSummary struct {
    TotalDeals     int           `json:"totalDeals"`
    PendingDeals   int           `json:"pendingDeals"`
    CompletedDeals int           `json:"completedDeals"`
    RecentDeals    []PipelineDeal `json:"recentDeals"`
}

// PipelineDeal represents a deal in an investor's pipeline
type PipelineDeal struct {
    ID              primitive.ObjectID `json:"id"`
    FounderID       primitive.ObjectID `json:"founder_id"`
    StartupName     string             `json:"startup_name"`
    Industry        string             `json:"industry"`
    Status          string             `json:"status"`
    Stage           string             `json:"stage"`
    MatchPercentage float64            `json:"match_percentage"`
    UpdatedAt       time.Time          `json:"updated_at"`
}

// PerformanceMetrics represents performance metrics for an investor's portfolio
type PerformanceMetrics struct {
    TotalReturn      float64 `json:"totalReturn"`
    AnnualizedReturn float64 `json:"annualizedReturn"`
    Volatility       float64 `json:"volatility"`
}