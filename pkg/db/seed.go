package db

import (
	"errors"
	"log/slog"

	"knowledge-capsule/app/models"
	"knowledge-capsule/pkg/utils"

	"gorm.io/gorm"
)

// SeedSuperAdmin creates or updates the superadmin user when env vars are set.
// Call after Open(). If email, password, name are all non-empty, ensures user exists with role=superadmin.
func SeedSuperAdmin(db *gorm.DB, email, password, name string) error {
	if email == "" || password == "" {
		return nil
	}
	if name == "" {
		name = "Super Admin"
	}

	hash, err := utils.HashPassword(password)
	if err != nil {
		return err
	}

	var existing models.User
	err = db.Where("email = ?", email).First(&existing).Error
	if err == nil {
		// Update existing user to superadmin
		if err := db.Model(&existing).Updates(map[string]interface{}{
			"name":          name,
			"password_hash": hash,
			"role":          models.RoleSuperAdmin,
		}).Error; err != nil {
			return err
		}
		slog.Info("Superadmin updated", "email", email)
		return nil
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// Create new superadmin
	user := models.User{
		ID:           utils.GenerateUUID(),
		Name:         name,
		Email:        email,
		PasswordHash: hash,
		Role:         models.RoleSuperAdmin,
	}
	if err := db.Create(&user).Error; err != nil {
		return err
	}
	slog.Info("Superadmin seeded", "email", email)
	return nil
}
