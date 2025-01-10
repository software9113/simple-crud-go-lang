package controllers

import (
	"gin-tutorial/models"
	"gin-tutorial/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UserController defines the interface for the user controller
type UserController interface {
	RegisterUser(c *gin.Context)
	Login(c *gin.Context)
	GetProfile(c *gin.Context)
}

// userControllerImpl is the concrete implementation of UserController
type userControllerImpl struct {
	userService services.UserService
}

// NewUserController creates a new UserController instance
func NewUserController(userService services.UserService) UserController {
	return &userControllerImpl{
		userService: userService,
	}
}

// @Summary Register a new user
// @Description Create a new user with a username, email, and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body models.RegisterUserRequest true "User registration details"
// @Success 200 {object} models.MessageResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /register [post]
func (uc *userControllerImpl) RegisterUser(c *gin.Context) {
	var input models.RegisterUserRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := uc.userService.RegisterUser(input.Username, input.Email, input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User registered successfully",
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

// @Summary Login a user
// @Description Authenticate a user and return a JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Param credentials body models.LoginRequest true "User login credentials"
// @Success 200 {object} models.TokenResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /login [post]
func (uc *userControllerImpl) Login(c *gin.Context) {
	var input models.LoginRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := uc.userService.LoginUser(input.Email, input.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	token, err := uc.userService.GenerateJWT(user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// @Summary Get user profile
// @Description Retrieve the currently authenticated user's profile
// @Tags User
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.ProfileResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /profile [get]
func (uc *userControllerImpl) GetProfile(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user from context"})
		return
	}

	c.JSON(http.StatusOK, user)
}
