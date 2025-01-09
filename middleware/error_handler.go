package middleware

import (
	"net/http"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// ErrorHandler is a middleware to handle errors and log their origin
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Capture the error location (file, line, function)
				var file string
				var line int
				pc, filePath, lineNo, ok := runtime.Caller(2)
				if ok {
					file = filePath
					line = lineNo
				}

				// Get the function name
				functionName := runtime.FuncForPC(pc).Name()

				// Log the error with its origin
				logrus.WithFields(logrus.Fields{
					"error":    err,
					"file":     file,
					"line":     line,
					"function": functionName,
				}).Error("An error occurred")

				// Respond with a generic error message
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":    "Internal Server Error",
					"file":     file,
					"line":     line,
					"function": functionName,
				})
			}
		}()

		// Continue to the next middleware or route handler
		c.Next()
	}
}
