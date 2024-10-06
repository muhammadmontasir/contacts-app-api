package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/muhammadmontasir/contact-app-api/internal/models"
	"github.com/muhammadmontasir/contact-app-api/pkg/utils"
	"gorm.io/gorm"
)

// Helper function to extract user ID from Authorization header
func getUserIDFromRequest(r *http.Request) (uint64, error) {
	authHeader := r.Header.Get("Authorization")
	token, err := utils.ExtractBearerToken(authHeader)
	if err != nil {
		return 0, err
	}

	userID, err := utils.GetUserIDFromToken(token)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func GetAllContacts(w http.ResponseWriter, r *http.Request) {
	// Get user ID from the token
	userID, err := getUserIDFromRequest(r)
	if err != nil {
		http.Error(w, "Failed to get user ID from token: "+err.Error(), http.StatusUnauthorized)
		return
	}

	// Implement pagination
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page == 0 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if pageSize == 0 {
		pageSize = 10
	}

	var contacts []models.Contact
	var total int64

	offset := (page - 1) * pageSize

	result := db.Model(&models.Contact{}).Where("user_id = ?", userID).Count(&total)
	if result.Error != nil {
		http.Error(w, "Failed to count contacts", http.StatusInternalServerError)
		return
	}

	result = db.Where("user_id = ?", userID).Offset(offset).Limit(pageSize).Find(&contacts)
	if result.Error != nil {
		http.Error(w, "Failed to retrieve contacts", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"contacts": contacts,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func CreateContact(w http.ResponseWriter, r *http.Request) {
	var contact models.Contact
	err := json.NewDecoder(r.Body).Decode(&contact)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get user ID from the token
	userID, err := getUserIDFromRequest(r)
	if err != nil {
		http.Error(w, "Failed to get user ID from token: "+err.Error(), http.StatusUnauthorized)
		return
	}

	contact.UserID = userID

	result := db.Create(&contact)
	if result.Error != nil {
		http.Error(w, "Failed to create contact", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(contact)
}

func GetContactDetails(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid contact ID", http.StatusBadRequest)
		return
	}

	userID, err := getUserIDFromRequest(r)
	if err != nil {
		http.Error(w, "Failed to get user ID from token: "+err.Error(), http.StatusUnauthorized)
		return
	}

	var contact models.Contact
	result := db.Where("id = ? AND user_id = ?", id, userID).First(&contact)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			http.Error(w, "Contact not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to retrieve contact", http.StatusInternalServerError)
		}
		return
	}

	json.NewEncoder(w).Encode(contact)
}

func UpdateContact(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid contact ID", http.StatusBadRequest)
		return
	}

	var updatedContact models.Contact
	err = json.NewDecoder(r.Body).Decode(&updatedContact)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID, err := getUserIDFromRequest(r)
	if err != nil {
		http.Error(w, "Failed to get user ID from token: "+err.Error(), http.StatusUnauthorized)
		return
	}

	var contact models.Contact
	result := db.Where("id = ? AND user_id = ?", id, userID).First(&contact)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			http.Error(w, "Contact not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to retrieve contact", http.StatusInternalServerError)
		}
		return
	}

	contact.Name = updatedContact.Name
	contact.Email = updatedContact.Email
	contact.Phone = updatedContact.Phone

	result = db.Save(&contact)
	if result.Error != nil {
		http.Error(w, "Failed to update contact", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(contact)
}

func DeleteContact(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid contact ID", http.StatusBadRequest)
		return
	}

	userID, err := getUserIDFromRequest(r)
	if err != nil {
		http.Error(w, "Failed to get user ID from token: "+err.Error(), http.StatusUnauthorized)
		return
	}

	result := db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Contact{})
	if result.Error != nil {
		http.Error(w, "Failed to delete contact", http.StatusInternalServerError)
		return
	}

	if result.RowsAffected == 0 {
		http.Error(w, "Contact not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
