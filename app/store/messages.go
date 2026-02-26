package store

import (
	"knowledge-capsule/app/models"
	"knowledge-capsule/pkg/utils"
	"time"
)

type MessageStore struct {
	FileStore[models.Message]
}

// SaveMessage saves a new message.
func (s *MessageStore) SaveMessage(senderID, receiverID, content string, msgType models.MessageType, fileURL string) (*models.Message, error) {
	messages, err := s.Load()
	if err != nil {
		return nil, err
	}

	newMsg := models.Message{
		ID:         utils.GenerateUUID(),
		SenderID:   senderID,
		ReceiverID: receiverID,
		Content:    content,
		Type:       msgType,
		FileURL:    fileURL,
		CreatedAt:  time.Now(),
	}

	messages = append(messages, newMsg)
	if err := s.Save(messages); err != nil {
		return nil, err
	}
	return &newMsg, nil
}

// GetMessagesBetweenUsers retrieves conversation history between two users.
func (s *MessageStore) GetMessagesBetweenUsers(user1ID, user2ID string) ([]models.Message, error) {
	messages, err := s.Load()
	if err != nil {
		return nil, err
	}

	var conversation []models.Message
	for _, msg := range messages {
		if (msg.SenderID == user1ID && msg.ReceiverID == user2ID) ||
			(msg.SenderID == user2ID && msg.ReceiverID == user1ID) {
			conversation = append(conversation, msg)
		}
	}
	return conversation, nil
}
