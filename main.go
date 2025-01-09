package main

import (
	"gin-tutorial/config"
	"gin-tutorial/controllers"
	"gin-tutorial/database"
	"gin-tutorial/docs"
	"gin-tutorial/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Gin Tutorial API
// @version 1.0
// @description This is a sample Gin application for learning purposes
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	cfg := config.LoadConfig()

	// Configure logging
	config.ConfigureLogger()

	// Connect to database
	database.ConnectDatabase(cfg.DatabaseURL)

	// Run migrations
	database.RunMigrations(database.DB)

	r := gin.Default()

	// Middleware
	r.Use(middleware.CORS())
	r.Use(middleware.LoggerAndErrorHandlerMiddleware())

	// Swagger
	docs.SwaggerInfo.BasePath = "/"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Routes
	r.POST("/register", controllers.RegisterUser)
	r.POST("/login", controllers.Login)

	authorized := r.Group("/").Use(middleware.AuthMiddleware(cfg.JWTSecret))
	authorized.GET("/profile", controllers.GetProfile)

	r.Run(":" + cfg.Port)
}
