package config

import (
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

// ConfigureLogger sets up logrus to write logs to daily rotated files
func ConfigureLogger() {
	// Create the logs directory if it doesn't exist
	logDir := "logs"
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		err := os.MkdirAll(logDir, os.ModePerm)
		if err != nil {
			logrus.Fatal("Failed to create log directory:", err)
		}
	}

	// Configure lumberjack for daily log rotation
	logFile := &lumberjack.Logger{
		Filename:   filepath.Join(logDir, time.Now().Format("2006-01-02")+".log"), // Log filename with current date
		MaxSize:    10,                                                            // Maximum size in MB before rotating
		MaxBackups: 7,                                                             // Maximum number of old log files to keep
		MaxAge:     30,                                                            // Maximum number of days to retain old log files
		Compress:   true,                                                          // Compress old log files
	}

	// Set log output to both file and console
	multiWriter := io.MultiWriter(logFile, os.Stdout)
	logrus.SetOutput(multiWriter)

	// Set log format (plain text or JSON)
	logrus.SetFormatter(&logrus.JSONFormatter{}) // Use JSON format for structured logs

	// Set the log level
	logrus.SetLevel(logrus.DebugLevel) // Set to DebugLevel for development
}
