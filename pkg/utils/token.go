package utils

import (
	"errors"
	"fmt"
	"log"
	"strings"
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

func ExtractBearerToken(authHeader string) (string, error) {
	if authHeader == "" {
		return "", errors.New("Authorization header missing")
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", errors.New("Authorization header format must be Bearer {token}")
	}

	return parts[1], nil
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

	returnedToken, err := token.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		log.Println("Error signing token:", err)
		return "", err
	}

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

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.JWTSecret), nil
	})

	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok := claims["user_id"].(float64)
		if !ok {
			return 0, errors.New("user ID not found in token")
		}
		return uint64(userID), nil
	}

	return 0, errors.New("invalid token")
}
