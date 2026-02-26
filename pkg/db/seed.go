package db

import (
	"errors"

	"knowledge-capsule/app/models"
	"knowledge-capsule/pkg/logger"
	"knowledge-capsule/pkg/utils"

	"gorm.io/gorm"
)

// SeedSuperAdmin creates or updates the superadmin user when env vars are set.
// Call after Open(). If email, password, name are all non-empty, ensures user exists with role=superadmin.
func SeedSuperAdmin(db *gorm.DB, email, password, name string) error {
	if email == "" || password == "" {
		logger.Info(logger.EventSeed, logger.Attr("action", "skipped"), logger.Attr("reason", "email or password not set"))
		return nil
	}
	if name == "" {
		name = "Super Admin"
	}

	hash, err := utils.HashPassword(password)
	if err != nil {
		logger.Error(logger.EventSeed, err, logger.Attr("action", "hash_password"), logger.Attr("email", email))
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
			logger.Error(logger.EventSeed, err, logger.Attr("action", "update_superadmin"), logger.Attr("email", email))
			return err
		}
		logger.Info(logger.EventSeed, logger.Attr("action", "superadmin_updated"), logger.Attr("email", email), logger.Attr("user_id", existing.ID))
		return nil
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Error(logger.EventSeed, err, logger.Attr("action", "find_superadmin"), logger.Attr("email", email))
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
		logger.Error(logger.EventSeed, err, logger.Attr("action", "create_superadmin"), logger.Attr("email", email))
		return err
	}
	logger.Info(logger.EventSeed, logger.Attr("action", "superadmin_created"), logger.Attr("email", email), logger.Attr("user_id", user.ID))
	return nil
}
