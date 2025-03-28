package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Service interface {
	Health() map[string]string
	User() UserService
	DealFlow() DealFlowService
	Founder() FounderService
	Investor() InvestorService
	Investment() InvestmentService
}

type service struct {
	db         *mongo.Client
	user       UserService
	dealFlow   DealFlowService
	founder    FounderService
	investor   InvestorService
	investment InvestmentService
}

var (
	host = os.Getenv("BLUEPRINT_DB_HOST")
	port = os.Getenv("BLUEPRINT_DB_PORT")

// database = os.Getenv("BLUEPRINT_DB_DATABASE")
)

func New() Service {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s", host, port)))
	if err != nil {
		log.Fatal(err)
	}
	
	dbName := os.Getenv("BLUEPRINT_DB_DATABASE")
	if dbName == "" {
		dbName = "ddb" // Fallback name
	}
	
	db := client.Database(dbName)
	
	return &service{
		db:         client,
		user:       NewUserService(client),
		investor:   NewInvestorService(client),
		founder:    NewFounderService(client),
		dealFlow:   NewDealFlowService(client),
		investment: NewInvestmentService(db),
	}
}

func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := s.db.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("db down: %v", err)
	}

	return map[string]string{
		"message": "It's healthy",
	}
}

func (s *service) User() UserService {
	return s.user
}

func (s *service) DealFlow() DealFlowService {
	return s.dealFlow
}

func (s *service) Investor() InvestorService {
	return s.investor
}

func (s *service) Founder() FounderService {
	return s.founder
}

func (s *service) Investment() InvestmentService {
	return s.investment
}
