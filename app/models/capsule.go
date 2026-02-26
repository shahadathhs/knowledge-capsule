package models

import (
	"time"
)

// CapsuleInput request body for POST /api/capsules and PUT /api/capsules/{id}
type CapsuleInput struct {
	Title     string   `json:"title" example:"Interfaces in Go" gorm:"not null"`
	Content   string   `json:"content" example:"Interfaces are named collections of method signatures..."`
	Topic     string   `json:"topic" example:"Golang"`
	Tags      Tags     `json:"tags" example:"programming,go" gorm:"type:jsonb"`
	IsPrivate bool     `json:"is_private" example:"false" gorm:"default:false"`
}

// Capsule extends CapsuleInput with ID, UserID and timestamps
type Capsule struct {
	ID        string    `json:"id" gorm:"primaryKey;type:varchar(36)"`
	UserID    string    `json:"user_id" gorm:"index;not null"`
	CapsuleInput
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Capsule) TableName() string { return "capsules" }
