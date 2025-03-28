package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Unified user profile
type User struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	FirstName  string             `bson:"first_name"`
	SecondName string             `bson:"second_name"`
	Email      string             `bson:"email"`
	Password   string             `bson:"password"`
	Roles      []string           `bson:"roles"`
	CreatedAt  time.Time          `bson:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at"`
}

// Founder profile
type Founder struct {
	ID                primitive.ObjectID `bson:"_id,omitempty"`
	UserID            primitive.ObjectID `bson:"user_id"`
	StartupName       string             `bson:"startup_name"`
	MissionStatement  string             `bson:"mission_statement"`
	Industry          string             `bson:"industry"`
	FundingStage      string             `bson:"funding_stage"`
	Avatar            string             `bson:"avatar"`
	FundingAllocation string             `bson:"funding_allocation"`
	BussinessModel    string             `bson:"bussiness_model"`
	RevenueStreams    string             `bson:"revenue_streams"`
	Traction          string             `bson:"traction"`
	TotalInvested     int                `bson:"total_invested"`
	FundRequired      int                `bson:"fund_required"`
	YearFounded       string             `bson:"year_founded"`
	Founded           string             `bson:"founded_stage"`
	ScalingPotential  string             `bson:"scaling_potential"`
	Competition       string             `bson:"competition"`
	LeadershipTeam    string             `bson:"leadership_team"`
	TeamSize          string             `bson:"team_size"`
	Location          string             `bson:"location"`
	StartupWebsite    string             `bson:"startup_website"`
	PitchDeck         string             `bson:"pitch_deck"`
	CreatedAt         time.Time          `bson:"created_at"`
	UpdatedAt         time.Time          `bson:"updated_at"`
}

// Investor profile
type Investor struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	UserID primitive.ObjectID `bson:"user_id"`
	//	InvestmentPortfolio   []primitive.ObjectID `bson:"investment_portfolio"`
	TotalInvested         float64   `bson:"total_invested"`
	InvestorType          string    `bson:"investor_type"`
	Thesis                string    `bson:"thesis"`
	PreferredFundingStage string    `bson:"preferred_funding_stage"`
	InvestmentRange       string    `bson:"investment_range"`
	InvestmentFrequency   string    `bson:"investment_frequency"`
	RiskTolerance         string    `bson:"risk_tolerance"`
	ExitStrategy          string    `bson:"exit_strategy"`
	PreferredIndustries   []string  `bson:"preferred_industries"`
	PreferredRegions      []string  `bson:"preferred_regions"`
	CreatedAt             time.Time `bson:"created_at"`
	UpdatedAt             time.Time `bson:"updated_at"`
}

// Admin-specific details
type Admin struct {
	UserID    primitive.ObjectID `bson:"user_id"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

// Investor & Founder Matching
type MatchInvestorFounder struct {
	ID              primitive.ObjectID `bson:"_id,omitempty"`
	FounderID       primitive.ObjectID `bson:"founder_id"`
	InvestorID      primitive.ObjectID `bson:"investor_id"`
	MatchPercentage float64            `bson:"match_percentage"`
	Tags            []string           `bson:"tags"`
	Bookmark        bool               `bson:"bookmark"`
	CreatedAt       time.Time          `bson:"created_at"`
	UpdatedAt       time.Time          `bson:"updated_at"`
}

type Meeting struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	InvestorID    primitive.ObjectID `bson:"investor_id" json:"investor_id"`
	FounderID     primitive.ObjectID `bson:"founder_id" json:"founder_id"`
	Title         string             `bson:"title" json:"title"`
	StartTime     time.Time          `bson:"start_time" json:"start_time"`
	EndTime       time.Time          `bson:"end_time" json:"end_time"`
	GoogleMeetURL string             `bson:"google_meet_url" json:"google_meet_url"`
	Notes         string             `bson:"notes" json:"notes"`
	CreatedAt     time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time          `bson:"updated_at" json:"updated_at"`
}

// Notification Model
type Notification struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`
	//	UserID           primitive.ObjectID `bson:"user_id"`
	FounderID        primitive.ObjectID `bson:"founder_id"`
	NotificationType string             `bson:"notification_type"`
	Title            string             `bson:"title"`
	Message          string             `bson:"message"`
	ReadStatus       bool               `bson:"read_status"`
	CreatedAt        time.Time          `bson:"created_at"`
	UpdatedAt        time.Time          `bson:"updated_at"`
}

// Deal Flow model for tracking investment pipeline
type DealFlow struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	InvestorID   primitive.ObjectID `bson:"investor_id"`
	StartupID    primitive.ObjectID `bson:"founder_id"`
	Stage        string             `bson:"stage"`
	MatchScore   float64            `bson:"match_score"`
	Status       string             `bson:"status"`
	Priority     string             `bson:"priority"`
	LastActivity time.Time          `bson:"last_activity"`
	AddedDate    time.Time          `bson:"added_date"`
	Meetings     []Meeting          `bson:"meetings"`
	Documents    []Document         `bson:"documents"`
	Tasks        []Task             `bson:"tasks"`
	CreatedAt    time.Time          `bson:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at"`
}

// Communication model
type Communication struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	InvestorID primitive.ObjectID `bson:"investor_id"`
	FounderID  primitive.ObjectID `bson:"founder_id"`
	Status     string             `bson:"status"`
	CreatedAt  time.Time          `bson:"created_at"`
}

// // Investment model for tracking finalized investments
// type Investment struct {
// 	ID                primitive.ObjectID `bson:"_id,omitempty"`
// 	InvestorID        primitive.ObjectID `bson:"investor_id"`
// 	StartupID         primitive.ObjectID `bson:"startup_id"`
// 	InvestmentAmount  float64            `bson:"investment_amount"`
// 	EquityPercentage  float64            `bson:"equity_percentage"`
// 	CurrentValuation  float64            `bson:"current_valuation"`
// 	InitialValuation  float64            `bson:"initial_valuation"`
// 	ROI               float64            `bson:"roi"`
// 	Status            string             `bson:"status"`
// 	Performance       string             `bson:"performance"`
// 	NextMilestone     string             `bson:"next_milestone"`
// 	NextMilestoneDate time.Time          `bson:"next_milestone_date"`
// 	Metrics           Metrics            `bson:"metrics"`
// 	Documents         []Document         `bson:"documents"`
// 	CreatedAt         time.Time          `bson:"created_at"`
// 	UpdatedAt         time.Time          `bson:"updated_at"`
// }

// Metrics model for tracking startup performance
type Metrics struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	Revenue RevenueMetrics     `bson:"revenue"`
	Users   UserMetrics        `bson:"users"`
	Burn    BurnMetrics        `bson:"burn"`
	Runway  RunwayMetrics      `bson:"runway"`
}

// Revenue metrics
type RevenueMetrics struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Current  float64            `bson:"current"`
	Previous float64            `bson:"previous"`
	Growth   float64            `bson:"growth"`
}

// User metrics
type UserMetrics struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Current  int                `bson:"current"`
	Previous int                `bson:"previous"`
	Growth   float64            `bson:"growth"`
}

// Burn metrics (Expenses)
type BurnMetrics struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Current  float64            `bson:"current"`
	Previous float64            `bson:"previous"`
	Growth   float64            `bson:"growth"`
}

// Runway metrics (Months of cash remaining)
type RunwayMetrics struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Current  int                `bson:"current"`
	Previous int                `bson:"previous"`
	Growth   float64            `bson:"growth"`
}

// Document model for storing files
type Document struct {
	ID   primitive.ObjectID `bson:"_id,omitempty"`
	Name string             `bson:"name"`
	URL  string             `bson:"url"`
	Date time.Time          `bson:"date"`
}

// Task model for investment tracking
type Task struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Title     string             `bson:"title"`
	Completed bool               `bson:"completed"`
	DueDate   time.Time          `bson:"due_date"`
	Priority  string             `bson:"priority"`
	CreatedBy primitive.ObjectID `bson:"created_by"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

// Note model for deal flow notes
type Note struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Content   string             `bson:"content" json:"content"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

// Activity represents an investor activity record
type Activity struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	InvestorID  primitive.ObjectID `bson:"investor_id"`
	Type        string             `bson:"type"`
	Description string             `bson:"description"`
	Date        time.Time          `bson:"date"`
}

// Match represents a match between a founder and an investor
type Match struct {
	ID                    primitive.ObjectID `bson:"_id,omitempty"`
	FounderID             primitive.ObjectID `bson:"founder_id"`
	InvestorID            primitive.ObjectID `bson:"investor_id"`
	InvestorName          string             `bson:"investor_name"`
	MatchScore            float64            `bson:"match_score"`
	Industry              string             `bson:"industry"`
	TotalInvested         float64            `bson:"total_invested"`
	PreferredFundingStage string             `bson:"preferred_funding_stage"`
	Stage                 string             `bson:"stage"`
	CreatedAt             time.Time          `bson:"created_at"`
	UpdatedAt             time.Time          `bson:"updated_at"`
}
