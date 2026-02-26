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
