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

// GetCapsulesByUser returns capsules owned by a user with optional filters.
func (s *capsuleStore) GetCapsulesByUser(userID string, filters *models.CapsuleFilters) ([]models.Capsule, error) {
	query := s.DB.Where("user_id = ?", userID)

	if filters != nil {
		if filters.Topic != "" {
			query = query.Where("topic ILIKE ?", "%"+filters.Topic+"%")
		}
		if len(filters.Tags) > 0 {
			for _, tag := range filters.Tags {
				query = query.Where("tags::text ILIKE ?", "%"+tag+"%")
			}
		}
		if filters.Q != "" {
			pattern := "%" + strings.ToLower(filters.Q) + "%"
			query = query.Where("LOWER(title) LIKE ? OR LOWER(content) LIKE ? OR tags::text ILIKE ?",
				pattern, pattern, pattern)
		}
		if filters.IsPrivate != nil {
			query = query.Where("is_private = ?", *filters.IsPrivate)
		}
	}

	var capsules []models.Capsule
	if err := query.Find(&capsules).Error; err != nil {
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

// SearchAllCapsules searches all capsules (admin only, no user filter).
func (s *capsuleStore) SearchAllCapsules(query string, limit int) ([]models.Capsule, error) {
	if limit <= 0 {
		limit = 20
	}
	pattern := "%" + strings.ToLower(query) + "%"
	var capsules []models.Capsule
	err := s.DB.Where("LOWER(title) LIKE ? OR LOWER(content) LIKE ? OR tags::text ILIKE ?",
		pattern, pattern, pattern).Limit(limit).Find(&capsules).Error
	return capsules, err
}
