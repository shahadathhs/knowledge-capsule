package db

import (
	"log/slog"

	"knowledge-capsule/app/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Open connects to PostgreSQL and runs auto-migrations.
func Open(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(
		&models.User{},
		&models.Topic{},
		&models.Capsule{},
		&models.Message{},
	); err != nil {
		return nil, err
	}

	slog.Info("Database connected and migrated")
	return db, nil
}
