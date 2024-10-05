package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/muhammadmontasir/contact-app-api/configs"
	"github.com/muhammadmontasir/contact-app-api/internal/database"
	"github.com/muhammadmontasir/contact-app-api/internal/handlers"
	"github.com/muhammadmontasir/contact-app-api/internal/middleware"
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

	// Set the database connection for the handlers
	handlers.SetDB(db)

	r := mux.NewRouter()
	api := r.PathPrefix("/api/v1").Subrouter()

	// User routes
	api.HandleFunc("/users", handlers.CreateUser).Methods("POST")
	api.HandleFunc("/users/activate", handlers.ActivateUser).Methods("POST")

	// Auth routes
	api.HandleFunc("/token/auth", handlers.AuthenticateUser).Methods("POST")

	// Contact routes (protected by JWT middleware)
	contacts := api.PathPrefix("/contacts").Subrouter()
	contacts.Use(middleware.JWTAuth)
	contacts.HandleFunc("", handlers.GetAllContacts).Methods("GET")
	contacts.HandleFunc("", handlers.CreateContact).Methods("POST")
	contacts.HandleFunc("/{id}", handlers.GetContactDetails).Methods("GET")
	contacts.HandleFunc("/{id}", handlers.UpdateContact).Methods("PATCH")
	contacts.HandleFunc("/{id}", handlers.DeleteContact).Methods("DELETE")

	log.Printf("Server starting on port %s...", config.ServerPort)
	err = http.ListenAndServe(":"+config.ServerPort, handlers.LoggingMiddleware(r))
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
