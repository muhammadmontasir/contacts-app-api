package utils

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/muhammadmontasir/contact-app-api/configs"
)

var (
	config      *configs.Config
	configMutex sync.RWMutex
)

// SetConfig initializes the config variable safely
func SetConfig(c *configs.Config) {
	configMutex.Lock()
	defer configMutex.Unlock()
	config = c
}

// getConfig safely retrieves the config variable
func getConfig() *configs.Config {
	configMutex.RLock()
	defer configMutex.RUnlock()
	return config
}

// GenerateToken creates a JWT token for a given user ID
func GenerateToken(userID uint64) (string, error) {
	cfg := getConfig()
	if cfg == nil {
		return "", fmt.Errorf("configuration not set")
	}

	if cfg.JWTSecret == "" {
		return "", fmt.Errorf("JWTSecret is not set in configuration")
	}

	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	log.Println("claims:", claims, "token:", token)

	log.Println("config.JWTSecret:", cfg.JWTSecret)
	returnedToken, err := token.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		log.Println("Error signing token:", err)
		return "", err
	}
	log.Println("returnedToken:", returnedToken, "err:", err)
	return returnedToken, nil
}

// VerifyToken validates a JWT token string
func VerifyToken(tokenString string) (bool, error) {
	cfg := getConfig()
	if cfg == nil {
		return false, fmt.Errorf("configuration not set")
	}

	if cfg.JWTSecret == "" {
		return false, fmt.Errorf("JWTSecret is not set in configuration")
	}

	// Add logging and check for empty token
	if tokenString == "" {
		log.Println("Error: Empty token string")
		return false, fmt.Errorf("empty token string")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.JWTSecret), nil
	})

	if err != nil {
		log.Printf("Error parsing token: %v", err)
		return false, fmt.Errorf("error parsing token: %w", err)
	}

	return token.Valid, nil
}

// GetUserIDFromToken extracts the user ID from a JWT token string
func GetUserIDFromToken(tokenString string) (uint64, error) {
	cfg := getConfig()
	if cfg == nil {
		return 0, fmt.Errorf("configuration not set")
	}

	if cfg.JWTSecret == "" {
		return 0, fmt.Errorf("JWTSecret is not set in configuration")
	}

	// Add logging and check for empty token
	if tokenString == "" {
		log.Println("Error: Empty token string")
		return 0, fmt.Errorf("empty token string")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.JWTSecret), nil
	})

	if err != nil {
		log.Printf("Error parsing token: %v", err)
		return 0, fmt.Errorf("error parsing token: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok := claims["user_id"].(float64)
		if !ok {
			log.Println("Error: User ID not found in token")
			return 0, fmt.Errorf("user ID not found in token")
		}
		return uint64(userID), nil
	}

	log.Println("Error: Invalid token")
	return 0, fmt.Errorf("invalid token")
}
