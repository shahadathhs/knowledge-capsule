package models

import (
	"time"
)

type Capsule struct {
	ID        string    `json:"id" gorm:"primaryKey;type:varchar(36)"`
	UserID    string    `json:"user_id" gorm:"index;not null"`
	Title     string    `json:"title" gorm:"not null"`
	Content   string    `json:"content"`
	Topic     string    `json:"topic"`
	Tags      Tags      `json:"tags" gorm:"type:jsonb"`
	IsPrivate bool      `json:"is_private" gorm:"default:false"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Capsule) TableName() string { return "capsules" }
