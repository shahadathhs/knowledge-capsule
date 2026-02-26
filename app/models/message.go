package models

import "time"

type MessageType string

const (
	MessageTypeText  MessageType = "text"
	MessageTypeImage MessageType = "image"
	MessageTypeAudio MessageType = "audio"
	MessageTypeFile  MessageType = "file"
)

type Message struct {
	ID         string      `json:"id" gorm:"primaryKey;type:varchar(36)"`
	SenderID   string      `json:"sender_id" gorm:"index;not null"`
	ReceiverID string      `json:"receiver_id" gorm:"index;not null"`
	Content    string      `json:"content,omitempty"`
	Type       MessageType `json:"type" gorm:"type:varchar(20);default:'text'"`
	FileURL    string      `json:"file_url,omitempty"`
	CreatedAt  time.Time   `json:"created_at"`
}

func (Message) TableName() string { return "messages" }
