package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func MiddlewareCorrelationID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Generate or retrieve the Correlation ID
		correlationID := c.GetHeader("X-Correlation-ID")
		if correlationID == "" {
			// If no correlation ID is provided, generate one
			correlationID = uuid.New().String()
		}
		c.Set("correlation_id", correlationID)

		// Log the correlation ID and request details
		startTime := time.Now()
		logEntry := logrus.WithFields(logrus.Fields{
			"correlation_id": correlationID,
			"application":    "gin-tutorial",        // Your application name
			"method":         c.Request.Method,      // HTTP method
			"url":            c.Request.URL.Path,    // Request URL
			"client_ip":      c.ClientIP(),          // Client IP
			"user_agent":     c.Request.UserAgent(), // User-Agent header
		})
		logEntry.Info("Received request")

		// Pass correlation ID to the response header for downstream systems
		c.Writer.Header().Set("X-Correlation-ID", correlationID)

		// Process the request
		c.Next()

		// Log response details after the request is handled
		endTime := time.Now()
		latency := endTime.Sub(startTime)
		statusCode := c.Writer.Status()

		logEntry.WithFields(logrus.Fields{
			"status_code":  statusCode,
			"latency":      latency.Seconds(),
			"response_len": c.Writer.Size(),
		}).Info("Completed request")
	}
}
