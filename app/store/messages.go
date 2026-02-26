package store

import (
	"knowledge-capsule/app/models"
	"knowledge-capsule/pkg/utils"

	"gorm.io/gorm"
)

// messageStore implements message storage with GORM.
type messageStore struct {
	DB *gorm.DB
}

// NewMessageStore returns a MessageStore backed by GORM.
func NewMessageStore(db *gorm.DB) MessageStore {
	return &messageStore{DB: db}
}

// SaveMessage saves a new message.
func (s *messageStore) SaveMessage(senderID, receiverID, content string, msgType models.MessageType, fileURL string) (*models.Message, error) {
	msg := models.Message{
		ID:         utils.GenerateUUID(),
		SenderID:   senderID,
		ReceiverID: receiverID,
		Content:    content,
		Type:       msgType,
		FileURL:    fileURL,
	}
	if err := s.DB.Create(&msg).Error; err != nil {
		return nil, err
	}
	return &msg, nil
}

// GetMessagesBetweenUsers retrieves conversation history between two users.
func (s *messageStore) GetMessagesBetweenUsers(user1ID, user2ID string) ([]models.Message, error) {
	var messages []models.Message
	err := s.DB.Where(
		"(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)",
		user1ID, user2ID, user2ID, user1ID,
	).Order("created_at ASC").Find(&messages).Error
	if err != nil {
		return nil, err
	}
	return messages, nil
}
