package config

import (
    "log"
    "os"

    "github.com/joho/godotenv"
)

type Config struct {
    Port        string
    DatabaseURL string
    JWTSecret   string
}

func LoadConfig() *Config {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    return &Config{
        Port:        getEnv("PORT", "8080"),
        DatabaseURL: getEnv("DATABASE_URL", ""),
        JWTSecret:   getEnv("JWT_SECRET", ""),
    }
}

func getEnv(key, defaultValue string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value
    }
    return defaultValue
}
