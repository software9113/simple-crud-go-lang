package models

import "github.com/golang-jwt/jwt/v4"

// Claims defines custom claims for JWT
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}
