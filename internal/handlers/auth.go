package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/muhammadmontasir/contact-app-api/internal/models"
	"github.com/muhammadmontasir/contact-app-api/pkg/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var db *gorm.DB

// SetDB sets the database connection for the handlers
func SetDB(database *gorm.DB) {
	db = database
}

func AuthenticateUser(w http.ResponseWriter, r *http.Request) {
	var credentials models.Credentials
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var user models.User
	result := db.Where("email = ?", credentials.Email).First(&user)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		} else {
			http.Error(w, "Database error", http.StatusInternalServerError)
		}
		return
	}

	// Check if the user is active
	if !user.Active {
		http.Error(w, "User account is not active", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password))
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Generate token
	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

// Add this middleware
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received request: ", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
		log.Println("Finished processing request: ", r.Method, r.URL.Path)
	})
}
