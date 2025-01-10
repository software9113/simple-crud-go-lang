package repository

import (
	"gin-tutorial/models"

	"gorm.io/gorm"
)

// UserRepository defines the methods for user-related database operations
type UserRepository interface {
	FindByEmail(email string) (*models.User, error)
	FindByID(userID uint) (*models.User, error)
	Create(user *models.User) error
}

// userRepositoryImpl is the concrete implementation of UserRepository
type userRepositoryImpl struct {
	db *gorm.DB
}

// NewUserRepository creates a new instance of UserRepository
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepositoryImpl{db: db}
}

// FindByEmail finds a user by email
func (ur *userRepositoryImpl) FindByEmail(email string) (*models.User, error) {
	var user models.User
	if err := ur.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByID finds a user by ID
func (ur *userRepositoryImpl) FindByID(userID uint) (*models.User, error) {
	var user models.User
	if err := ur.db.First(&user, userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Create saves a new user in the database
func (ur *userRepositoryImpl) Create(user *models.User) error {
	return ur.db.Create(user).Error
}
