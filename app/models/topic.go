package models

import (
	"time"
)

// TopicInput request body for POST /api/topics and PUT /api/topics/{id}
type TopicInput struct {
	Name        string `json:"name" example:"Golang" gorm:"not null"`
	Description string `json:"description" example:"Go programming language"`
}

// Topic extends TopicInput with ID and timestamps
type Topic struct {
	ID        string    `json:"id" gorm:"primaryKey;type:varchar(36)"`
	TopicInput
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Topic) TableName() string { return "topics" }
