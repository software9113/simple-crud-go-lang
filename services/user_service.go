package services

import (
	"errors"
	"gin-tutorial/models"
	"gin-tutorial/repository"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var jwtSecretKey = []byte("your_jwt_secret_key")

// UserService defines the interface for the user service
type UserService interface {
	RegisterUser(username, email, password string) (*models.User, error)
	LoginUser(email, password string) (*models.User, error)
	GenerateJWT(username string) (string, error)
	GetProfile(userID uint) (*models.User, error)
}

// userServiceImpl is the concrete implementation of UserService
type userServiceImpl struct {
	userRepo repository.UserRepository
}

// NewUserService creates a new UserService instance
func NewUserService(userRepo repository.UserRepository) UserService {
	return &userServiceImpl{
		userRepo: userRepo,
	}
}

// RegisterUser handles user registration
func (us *userServiceImpl) RegisterUser(username, email, password string) (*models.User, error) {
	// Check if user exists
	existingUser, _ := us.userRepo.FindByEmail(email)
	if existingUser != nil {
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
	if err := us.userRepo.Create(&user); err != nil {
		return nil, errors.New("failed to create user")
	}

	return &user, nil
}

// LoginUser authenticates the user
func (us *userServiceImpl) LoginUser(email, password string) (*models.User, error) {
	// Find user by email
	user, err := us.userRepo.FindByEmail(email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Check password
	if !user.CheckPassword(password) {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

// GenerateJWT generates a JWT token for the user
func (us *userServiceImpl) GenerateJWT(username string) (string, error) {
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

// GetProfile retrieves a user's profile by ID
func (us *userServiceImpl) GetProfile(userID uint) (*models.User, error) {
	return us.userRepo.FindByID(userID)
}
