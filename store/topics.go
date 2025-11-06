package store

import (
	"errors"
	"time"

	"knowledge-capsule-api/models"
	"knowledge-capsule-api/utils"
)

type TopicStore struct {
	FileStore[models.Topic]
}

// AddTopic creates a new topic.
func (s *TopicStore) AddTopic(name, description string) (*models.Topic, error) {
	topics, err := s.Load()
	if err != nil {
		return nil, err
	}

	// Check duplicate name
	for _, t := range topics {
		if t.Name == name {
			return nil, errors.New("topic already exists")
		}
	}

	newTopic := models.Topic{
		ID:          utils.GenerateUUID(),
		Name:        name,
		Description: description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	topics = append(topics, newTopic)
	if err := s.Save(topics); err != nil {
		return nil, err
	}
	return &newTopic, nil
}

// GetAllTopics returns all topics.
func (s *TopicStore) GetAllTopics() ([]models.Topic, error) {
	return s.Load()
}

// FindByID returns a topic by its ID.
func (s *TopicStore) FindByID(id string) (*models.Topic, error) {
	topics, err := s.Load()
	if err != nil {
		return nil, err
	}
	for _, t := range topics {
		if t.ID == id {
			return &t, nil
		}
	}
	return nil, errors.New("topic not found")
}
