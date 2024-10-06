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

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	user.Password = string(hashedPassword)

	result := db.Create(&user)
	if result.Error != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	activationToken, err := utils.GenerateToken(user.ID)
	if err != nil {
		http.Error(w, "Failed to generate activation token", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"activation_token": activationToken})
}

func ActivateUser(w http.ResponseWriter, r *http.Request) {
	var activation struct {
		Token string `json:"activation_token"`
	}
	err := json.NewDecoder(r.Body).Decode(&activation)

	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Check if token is empty
	if activation.Token == "" {
		http.Error(w, "Activation token is required", http.StatusBadRequest)
		return
	}

	// Verify token
	userID, err := utils.GetUserIDFromToken(activation.Token)
	if err != nil {
		http.Error(w, "Invalid or expired activation token", http.StatusBadRequest)
		return
	}

	// Activate user in database
	var user models.User
	result := db.First(&user, userID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to retrieve user", http.StatusInternalServerError)
		}
		return
	}

	log.Println("user", user)
	if user.Active {
		http.Error(w, "User is already activated", http.StatusBadRequest)
		return
	}

	user.Active = true
	result = db.Save(&user)
	if result.Error != nil {
		http.Error(w, "Failed to activate user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User activated successfully"})
}
