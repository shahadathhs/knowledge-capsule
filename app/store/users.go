package store

import (
	"errors"
	"time"

	"knowledge-capsule/app/models"
	"knowledge-capsule/pkg/utils"
)

type UserStore struct {
	FileStore[models.User]
}

// AddUser adds a new user if email is not taken.
func (s *UserStore) AddUser(name, email, password string) (*models.User, error) {
	users, err := s.Load()
	if err != nil {
		return nil, err
	}

	// Check duplicate email
	for _, u := range users {
		if u.Email == email {
			return nil, errors.New("email already exists")
		}
	}

	passwordHash, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	newUser := models.User{
		ID:           utils.GenerateUUID(),
		Name:         name,
		Email:        email,
		PasswordHash: passwordHash,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	users = append(users, newUser)
	if err := s.Save(users); err != nil {
		return nil, err
	}
	return &newUser, nil
}

// FindByEmail finds a user by email.
func (s *UserStore) FindByEmail(email string) (*models.User, error) {
	users, err := s.Load()
	if err != nil {
		return nil, err
	}

	for _, u := range users {
		if u.Email == email {
			return &u, nil
		}
	}

	return nil, errors.New("user not found")
}

// FindByID finds a user by ID.
func (s *UserStore) FindByID(id string) (*models.User, error) {
	users, err := s.Load()
	if err != nil {
		return nil, err
	}

	for _, u := range users {
		if u.ID == id {
			return &u, nil
		}
	}
	return nil, errors.New("user not found")
}
