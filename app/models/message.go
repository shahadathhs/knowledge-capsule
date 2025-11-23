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
	ID         string      `json:"id"`
	SenderID   string      `json:"sender_id"`
	ReceiverID string      `json:"receiver_id"`
	Content    string      `json:"content,omitempty"` // Text content
	Type       MessageType `json:"type"`
	FileURL    string      `json:"file_url,omitempty"`
	CreatedAt  time.Time   `json:"created_at"`
}
