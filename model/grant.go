package model

import (
    "time"
    
    "go.mongodb.org/mongo-driver/bson/primitive"
)

// Grant represents a funding grant opportunity
type Grant struct {
    ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Name        string             `bson:"name" json:"name"`
    Description string             `bson:"description" json:"description"`
    Amount      float64            `bson:"amount" json:"amount"`
    Category    string             `bson:"category" json:"category"`
    Region      string             `bson:"region" json:"region"`
    Deadline    time.Time          `bson:"deadline" json:"deadline"`
    Eligibility string             `bson:"eligibility" json:"eligibility"`
    CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
    UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

// GrantApplication represents an application for a grant
type GrantApplication struct {
    ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    FounderID       primitive.ObjectID `bson:"founder_id" json:"founder_id"`
    GrantID         int                `bson:"grant_id" json:"grant_id"`
    StartupName     string             `bson:"startup_name" json:"startup_name"`
    ContactEmail    string             `bson:"contact_email" json:"contact_email"`
    ContactPhone    string             `bson:"contact_phone" json:"contact_phone"`
    Description     string             `bson:"description" json:"description"`
    Website         string             `bson:"website" json:"website"`
    TeamSize        string             `bson:"team_size" json:"team_size"`
    PreviousFunding string             `bson:"previous_funding" json:"previous_funding"`
    PitchDeckPath   string             `bson:"pitch_deck_path" json:"pitch_deck_path"`
    Status          string             `bson:"status,omitempty" json:"status"`
    CreatedAt       time.Time          `bson:"created_at" json:"created_at"`
}
