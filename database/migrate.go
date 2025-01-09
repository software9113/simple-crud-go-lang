package database

import (
	"log"

	"gin-tutorial/models"

	"gorm.io/gorm"
)

// RunMigrations runs all database migrations and seed data
func RunMigrations(db *gorm.DB) {
	log.Println("Running migrations...")

	// Automatically migrate User model
	err := db.AutoMigrate(
		&models.User{}, // Add your models here
		// Add additional models here, e.g., &models.Product{}, &models.Order{}
	)
	if err != nil {
		log.Fatalf("Error running migrations: %v", err)
	}

	log.Println("Database migration completed successfully")

	// Seed data
	seedData(db)
}

// seedData seeds the database with initial data
func seedData(db *gorm.DB) {
	log.Println("Seeding data...")

	// Create an admin user with a hashed password
	admin := models.User{
		Username: "admin",
		Email:    "admin@example.com",
		Password: "password123", // Plain text password
	}

	// Hash the password
	if err := admin.HashPassword(); err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}

	// Save the admin user if not exists
	db.FirstOrCreate(&admin, models.User{Email: "admin@example.com"})

	log.Println("Seeding data completed")
}
