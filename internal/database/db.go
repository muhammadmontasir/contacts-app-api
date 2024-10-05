package database

import (
	"fmt"
	"log"

	"github.com/muhammadmontasir/contact-app-api/configs"
	"github.com/muhammadmontasir/contact-app-api/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(config *configs.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		config.DBHost,
		config.DBUser,
		config.DBPassword,
		config.DBName,
		config.DBPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Auto-migrate models
	db.AutoMigrate(&models.User{}, &models.Contact{})

	log.Println("Database auto-migration completed successfully")

	return db, nil
}
