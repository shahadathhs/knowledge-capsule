package store

import (
	"errors"

	"knowledge-capsule/app/models"
	"knowledge-capsule/pkg/utils"

	"gorm.io/gorm"
)

// userStore implements user storage with GORM.
type userStore struct {
	DB *gorm.DB
}

// NewUserStore returns a UserStore backed by GORM.
func NewUserStore(db *gorm.DB) UserStore {
	return &userStore{DB: db}
}

// AddUser adds a new user if email is not taken.
func (s *userStore) AddUser(name, email, password string) (*models.User, error) {
	var existing models.User
	if err := s.DB.Where("email = ?", email).First(&existing).Error; err == nil {
		return nil, errors.New("email already exists")
	}

	passwordHash, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := models.User{
		ID:           utils.GenerateUUID(),
		Name:         name,
		Email:        email,
		PasswordHash: passwordHash,
		Role:         models.RoleUser,
	}

	if err := s.DB.Create(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByEmail finds a user by email.
func (s *userStore) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := s.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

// FindByID finds a user by ID.
func (s *userStore) FindByID(id string) (*models.User, error) {
	var user models.User
	err := s.DB.First(&user, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

// ListUsers returns paginated users with optional search and role filter.
func (s *userStore) ListUsers(q, role string, page, limit int) ([]models.User, int, error) {
	var users []models.User
	query := s.DB.Model(&models.User{})

	if q != "" {
		pattern := "%" + q + "%"
		query = query.Where("name ILIKE ? OR email ILIKE ?", pattern, pattern)
	}
	if role != "" {
		query = query.Where("role = ?", role)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if offset < 0 {
		offset = 0
	}
	if err := query.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, 0, err
	}
	return users, int(total), nil
}

// UpdateProfile updates name and/or avatar_url for a user.
func (s *userStore) UpdateProfile(userID string, name, avatarURL string) error {
	updates := make(map[string]interface{})
	if name != "" {
		updates["name"] = name
	}
	if avatarURL != "" {
		updates["avatar_url"] = avatarURL
	}
	if len(updates) == 0 {
		return nil
	}
	result := s.DB.Model(&models.User{}).Where("id = ?", userID).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}
	return nil
}

// UpdateUserRole sets the role for a user (admin/superadmin only).
func (s *userStore) UpdateUserRole(id, role string) error {
	valid := map[string]bool{models.RoleUser: true, models.RoleAdmin: true, models.RoleSuperAdmin: true}
	if !valid[role] {
		return errors.New("invalid role")
	}
	result := s.DB.Model(&models.User{}).Where("id = ?", id).Update("role", role)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}
	return nil
}

// SearchUsers searches users by name or email.
func (s *userStore) SearchUsers(query string, limit int) ([]models.User, error) {
	if limit <= 0 {
		limit = 20
	}
	pattern := "%" + query + "%"
	var users []models.User
	err := s.DB.Where("name ILIKE ? OR email ILIKE ?", pattern, pattern).Limit(limit).Find(&users).Error
	return users, err
}

// ListAdmins returns users with role admin or superadmin.
func (s *userStore) ListAdmins(page, limit int) ([]models.User, int, error) {
	var users []models.User
	query := s.DB.Model(&models.User{}).Where("role IN ?", []string{models.RoleAdmin, models.RoleSuperAdmin})

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if offset < 0 {
		offset = 0
	}
	if err := query.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, 0, err
	}
	return users, int(total), nil
}
