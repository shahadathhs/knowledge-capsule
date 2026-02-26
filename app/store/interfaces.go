package store

import "knowledge-capsule/app/models"

// UserStore defines user storage operations.
type UserStore interface {
	AddUser(name, email, password string) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	FindByID(id string) (*models.User, error)
}

// CapsuleStore defines capsule storage operations.
type CapsuleStore interface {
	AddCapsule(userID, title, content, topic string, tags []string, isPrivate bool) (*models.Capsule, error)
	GetCapsulesByUser(userID string) ([]models.Capsule, error)
	FindByID(id string) (*models.Capsule, error)
	UpdateCapsule(id, userID string, updated models.Capsule) (*models.Capsule, error)
	DeleteCapsule(id, userID string) error
	SearchCapsules(userID, query string) ([]models.Capsule, error)
}

// TopicStore defines topic storage operations.
type TopicStore interface {
	AddTopic(name, description string) (*models.Topic, error)
	GetAllTopics() ([]models.Topic, error)
	FindByID(id string) (*models.Topic, error)
	UpdateTopic(id, name, description string) (*models.Topic, error)
	DeleteTopic(id string) error
}

// MessageStore defines message storage operations.
type MessageStore interface {
	SaveMessage(senderID, receiverID, content string, msgType models.MessageType, fileURL string) (*models.Message, error)
	GetMessagesBetweenUsers(user1ID, user2ID string) ([]models.Message, error)
}
