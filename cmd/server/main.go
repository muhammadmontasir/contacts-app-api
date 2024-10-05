package main

import (
	"log"

	"github.com/muhammadmontasir/contact-app-api/configs"
	"github.com/muhammadmontasir/contact-app-api/internal/database"
	"github.com/muhammadmontasir/contact-app-api/pkg/utils"
)

func main() {
	// Load configuration
	config, err := configs.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	utils.SetConfig(config)
	log.Println("Configuration has been set in utils package")

	// Initialize database connection
	db, err := database.InitDB(config)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
}
