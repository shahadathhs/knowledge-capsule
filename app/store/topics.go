package store

import (
	"errors"

	"knowledge-capsule/app/models"
	"knowledge-capsule/pkg/utils"

	"gorm.io/gorm"
)

// topicStore implements topic storage with GORM.
type topicStore struct {
	DB *gorm.DB
}

// NewTopicStore returns a TopicStore backed by GORM.
func NewTopicStore(db *gorm.DB) TopicStore {
	return &topicStore{DB: db}
}

// AddTopic creates a new topic.
func (s *topicStore) AddTopic(name, description string) (*models.Topic, error) {
	var existing models.Topic
	if err := s.DB.Where("name = ?", name).First(&existing).Error; err == nil {
		return nil, errors.New("topic already exists")
	}

	topic := models.Topic{
		ID: utils.GenerateUUID(),
		TopicInput: models.TopicInput{
			Name:        name,
			Description: description,
		},
	}
	if err := s.DB.Create(&topic).Error; err != nil {
		return nil, err
	}
	return &topic, nil
}

// GetAllTopics returns topics with optional search filter.
func (s *topicStore) GetAllTopics(filters *models.TopicFilters) ([]models.Topic, error) {
	query := s.DB.Model(&models.Topic{})

	if filters != nil && filters.Q != "" {
		pattern := "%" + filters.Q + "%"
		query = query.Where("name ILIKE ? OR description ILIKE ?", pattern, pattern)
	}

	var topics []models.Topic
	if err := query.Find(&topics).Error; err != nil {
		return nil, err
	}
	return topics, nil
}

// FindByID returns a topic by its ID.
func (s *topicStore) FindByID(id string) (*models.Topic, error) {
	var topic models.Topic
	err := s.DB.First(&topic, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("topic not found")
		}
		return nil, err
	}
	return &topic, nil
}

// UpdateTopic updates a topic by ID.
func (s *topicStore) UpdateTopic(id, name, description string) (*models.Topic, error) {
	result := s.DB.Model(&models.Topic{}).Where("id = ?", id).Updates(map[string]interface{}{
		"name":        name,
		"description": description,
	})
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("topic not found")
	}
	var topic models.Topic
	s.DB.First(&topic, "id = ?", id)
	return &topic, nil
}

// DeleteTopic removes a topic by ID.
func (s *topicStore) DeleteTopic(id string) error {
	result := s.DB.Where("id = ?", id).Delete(&models.Topic{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("topic not found")
	}
	return nil
}

// SearchTopics searches topics by name or description.
func (s *topicStore) SearchTopics(query string, limit int) ([]models.Topic, error) {
	if limit <= 0 {
		limit = 20
	}
	pattern := "%" + query + "%"
	var topics []models.Topic
	err := s.DB.Where("name ILIKE ? OR description ILIKE ?", pattern, pattern).Limit(limit).Find(&topics).Error
	return topics, err
}
