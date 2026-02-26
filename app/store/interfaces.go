package store

import "knowledge-capsule/app/models"

// UserStore defines user storage operations.
type UserStore interface {
	AddUser(name, email, password string) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	FindByID(id string) (*models.User, error)
	ListUsers(q, role string, page, limit int) ([]models.User, int, error)
	UpdateProfile(userID string, name, avatarURL string) error
	UpdateUserRole(id, role string) error
	SearchUsers(query string, limit int) ([]models.User, error)
	ListAdmins(page, limit int) ([]models.User, int, error)
}

// CapsuleStore defines capsule storage operations.
type CapsuleStore interface {
	AddCapsule(userID, title, content, topic string, tags []string, isPrivate bool) (*models.Capsule, error)
	GetCapsulesByUser(userID string, filters *models.CapsuleFilters) ([]models.Capsule, error)
	FindByID(id string) (*models.Capsule, error)
	UpdateCapsule(id, userID string, updated models.Capsule) (*models.Capsule, error)
	DeleteCapsule(id, userID string) error
	SearchAllCapsules(query string, limit int) ([]models.Capsule, error)
}

// TopicStore defines topic storage operations.
type TopicStore interface {
	AddTopic(name, description string) (*models.Topic, error)
	GetAllTopics(filters *models.TopicFilters) ([]models.Topic, error)
	FindByID(id string) (*models.Topic, error)
	UpdateTopic(id, name, description string) (*models.Topic, error)
	DeleteTopic(id string) error
	SearchTopics(query string, limit int) ([]models.Topic, error)
}

// MessageStore defines message storage operations.
type MessageStore interface {
	SaveMessage(senderID, receiverID, content string, msgType models.MessageType, fileURL string) (*models.Message, error)
	GetMessagesBetweenUsers(user1ID, user2ID string) ([]models.Message, error)
}
