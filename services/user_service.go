package services

import (
	"errors"
	"gin-tutorial/database"
	"gin-tutorial/models"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

// Secret key for JWT generation
var jwtSecretKey = []byte("your_jwt_secret_key")

// RegisterUserService handles user registration
func RegisterUserService(username, email, password string) (*models.User, error) {
	// Check if user exists
	var existingUser models.User
	if err := database.DB.Where("email = ?", email).First(&existingUser).Error; err == nil {
		return nil, errors.New("user with this email already exists")
	}

	// Create user object
	user := models.User{
		Username: username,
		Email:    email,
		Password: password,
	}

	// Hash password
	if err := user.HashPassword(); err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Save user to database
	if err := database.DB.Create(&user).Error; err != nil {
		return nil, errors.New("failed to create user")
	}

	return &user, nil
}

// LoginUserService authenticates the user
func LoginUserService(email, password string) (*models.User, error) {
	var user models.User

	// Find user by email
	if err := database.DB.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid credentials")
		}
		return nil, errors.New("database error")
	}

	// Check password
	if !user.CheckPassword(password) {
		return nil, errors.New("invalid credentials")
	}

	return &user, nil
}

// GenerateJWT generates a JWT token for the user
func GenerateJWT(username string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &models.Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecretKey)
}

// GetProfileService retrieves a user's profile by ID
func GetProfileService(userID uint) (*models.User, error) {
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, errors.New("database error")
	}

	return &user, nil
}
