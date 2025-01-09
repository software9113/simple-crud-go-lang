package models

// RegisterUserRequest defines the request body for user registration
type RegisterUserRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginRequest defines the request body for user login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// TokenResponse defines the response body containing the JWT token
type TokenResponse struct {
	Token string `json:"token"`
}

// ProfileResponse defines the response body for user profile
type ProfileResponse struct {
	User string `json:"user"`
}

// MessageResponse is a generic message response
type MessageResponse struct {
	Message string `json:"message"`
}

// ErrorResponse is a generic error response
type ErrorResponse struct {
	Error string `json:"error"`
}
