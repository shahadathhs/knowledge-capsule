package store

import (
	"errors"
	"strings"

	"knowledge-capsule/app/models"
	"knowledge-capsule/pkg/utils"

	"gorm.io/gorm"
)

// capsuleStore implements capsule storage with GORM.
type capsuleStore struct {
	DB *gorm.DB
}

// NewCapsuleStore returns a CapsuleStore backed by GORM.
func NewCapsuleStore(db *gorm.DB) CapsuleStore {
	return &capsuleStore{DB: db}
}

// AddCapsule creates a new capsule.
func (s *capsuleStore) AddCapsule(userID, title, content, topic string, tags []string, isPrivate bool) (*models.Capsule, error) {
	capsule := models.Capsule{
		ID:     utils.GenerateUUID(),
		UserID: userID,
		CapsuleInput: models.CapsuleInput{
			Title:     title,
			Content:   content,
			Topic:     topic,
			Tags:      models.Tags(tags),
			IsPrivate: isPrivate,
		},
	}
	if err := s.DB.Create(&capsule).Error; err != nil {
		return nil, err
	}
	return &capsule, nil
}

// GetCapsulesByUser returns all capsules owned by a specific user.
func (s *capsuleStore) GetCapsulesByUser(userID string) ([]models.Capsule, error) {
	var capsules []models.Capsule
	if err := s.DB.Where("user_id = ?", userID).Find(&capsules).Error; err != nil {
		return nil, err
	}
	return capsules, nil
}

// FindByID returns a capsule by its ID.
func (s *capsuleStore) FindByID(id string) (*models.Capsule, error) {
	var capsule models.Capsule
	err := s.DB.First(&capsule, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("capsule not found")
		}
		return nil, err
	}
	return &capsule, nil
}

// UpdateCapsule updates title, content, topic, and tags.
func (s *capsuleStore) UpdateCapsule(id, userID string, updated models.Capsule) (*models.Capsule, error) {
	result := s.DB.Model(&models.Capsule{}).Where("id = ? AND user_id = ?", id, userID).Updates(map[string]interface{}{
		"title":      updated.Title,
		"content":    updated.Content,
		"topic":      updated.Topic,
		"tags":       updated.Tags,
		"is_private": updated.IsPrivate,
	})
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("capsule not found or unauthorized")
	}
	var capsule models.Capsule
	s.DB.First(&capsule, "id = ?", id)
	return &capsule, nil
}

// DeleteCapsule removes a capsule by ID (only owner).
func (s *capsuleStore) DeleteCapsule(id, userID string) error {
	result := s.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Capsule{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("capsule not found or unauthorized")
	}
	return nil
}

// SearchCapsules performs a simple case-insensitive keyword search.
func (s *capsuleStore) SearchCapsules(userID, query string) ([]models.Capsule, error) {
	query = strings.ToLower(query)
	pattern := "%" + query + "%"
	var capsules []models.Capsule
	err := s.DB.Where("user_id = ? AND (LOWER(title) LIKE ? OR LOWER(content) LIKE ? OR tags::text ILIKE ?)",
		userID, pattern, pattern, pattern).Find(&capsules).Error
	if err != nil {
		return nil, err
	}
	return capsules, nil
}
