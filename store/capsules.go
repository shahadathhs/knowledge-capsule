package store

import (
	"errors"
	"strings"
	"time"

	"knowledge-capsule-api/models"
	"knowledge-capsule-api/utils"
)

type CapsuleStore struct {
	FileStore[models.Capsule]
}

// AddCapsule creates a new capsule.
func (s *CapsuleStore) AddCapsule(userID, title, content, topic string, tags []string, isPrivate bool) (*models.Capsule, error) {
	capsules, err := s.Load()
	if err != nil {
		return nil, err
	}

	newCapsule := models.Capsule{
		ID:        utils.GenerateUUID(),
		UserID:    userID,
		Title:     title,
		Content:   content,
		Topic:     topic,
		Tags:      tags,
		IsPrivate: isPrivate,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	capsules = append(capsules, newCapsule)
	if err := s.Save(capsules); err != nil {
		return nil, err
	}
	return &newCapsule, nil
}

// GetCapsulesByUser returns all capsules owned by a specific user.
func (s *CapsuleStore) GetCapsulesByUser(userID string) ([]models.Capsule, error) {
	capsules, err := s.Load()
	if err != nil {
		return nil, err
	}

	var result []models.Capsule
	for _, c := range capsules {
		if c.UserID == userID {
			result = append(result, c)
		}
	}
	return result, nil
}

// FindByID returns a capsule by its ID.
func (s *CapsuleStore) FindByID(id string) (*models.Capsule, error) {
	capsules, err := s.Load()
	if err != nil {
		return nil, err
	}

	for _, c := range capsules {
		if c.ID == id {
			return &c, nil
		}
	}
	return nil, errors.New("capsule not found")
}

// UpdateCapsule updates title, content, topic, and tags.
func (s *CapsuleStore) UpdateCapsule(id, userID string, updated models.Capsule) (*models.Capsule, error) {
	capsules, err := s.Load()
	if err != nil {
		return nil, err
	}

	for i, c := range capsules {
		if c.ID == id && c.UserID == userID {
			c.Title = updated.Title
			c.Content = updated.Content
			c.Topic = updated.Topic
			c.Tags = updated.Tags
			c.IsPrivate = updated.IsPrivate
			c.UpdatedAt = time.Now()
			capsules[i] = c
			s.Save(capsules)
			return &c, nil
		}
	}
	return nil, errors.New("capsule not found or unauthorized")
}

// DeleteCapsule removes a capsule by ID (only owner).
func (s *CapsuleStore) DeleteCapsule(id, userID string) error {
	capsules, err := s.Load()
	if err != nil {
		return err
	}

	newList := []models.Capsule{}
	found := false
	for _, c := range capsules {
		if c.ID == id && c.UserID == userID {
			found = true
			continue
		}
		newList = append(newList, c)
	}
	if !found {
		return errors.New("capsule not found or unauthorized")
	}

	return s.Save(newList)
}

// SearchCapsules performs a simple case-insensitive keyword search.
func (s *CapsuleStore) SearchCapsules(userID, query string) ([]models.Capsule, error) {
	capsules, err := s.Load()
	if err != nil {
		return nil, err
	}

	query = strings.ToLower(query)
	var result []models.Capsule

	for _, c := range capsules {
		if c.UserID != userID {
			continue
		}
		if strings.Contains(strings.ToLower(c.Title), query) ||
			strings.Contains(strings.ToLower(c.Content), query) ||
			containsTag(c.Tags, query) {
			result = append(result, c)
		}
	}
	return result, nil
}

func containsTag(tags []string, query string) bool {
	for _, t := range tags {
		if strings.Contains(strings.ToLower(t), query) {
			return true
		}
	}
	return false
}
