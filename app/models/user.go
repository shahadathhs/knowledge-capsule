package models

import (
	"time"
)

const (
	RoleUser       = "user"
	RoleAdmin      = "admin"
	RoleSuperAdmin = "superadmin"
)

type User struct {
	ID           string    `json:"id" gorm:"primaryKey;type:varchar(36)"`
	Name         string    `json:"name" gorm:"not null"`
	Email        string    `json:"email" gorm:"uniqueIndex;not null"`
	PasswordHash string    `json:"-" gorm:"column:password_hash;not null"`
	Role         string    `json:"role" gorm:"default:user;size:20"`
	AvatarURL    string    `json:"avatar_url" gorm:"column:avatar_url"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (User) TableName() string { return "users" }
