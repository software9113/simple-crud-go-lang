package middleware

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// responseBodyWriter is used to capture the response body
type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w responseBodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b) // Write the response body to the buffer
	return w.ResponseWriter.Write(b)
}

// LoggerAndErrorHandlerMiddleware logs requests, responses, and errors, and handles errors
func LoggerAndErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()

		// Read the request body
		bodyBytes, _ := ioutil.ReadAll(c.Request.Body)
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes)) // Reset request body for further use

		// Log request data
		logrus.WithFields(logrus.Fields{
			"method": c.Request.Method,
			"path":   c.Request.URL.Path,
			"query":  c.Request.URL.RawQuery,
			"body":   string(bodyBytes),
			"ip":     c.ClientIP(),
		}).Info("Incoming request")

		// Replace the ResponseWriter with a custom writer to capture the response body
		writer := &responseBodyWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = writer

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
				functionName := runtime.FuncForPC(pc).Name()

				// Log the error
				logrus.WithFields(logrus.Fields{
					"error":    err,
					"file":     file,
					"line":     line,
					"function": functionName,
					"method":   c.Request.Method,
					"path":     c.Request.URL.Path,
				}).Error("Unhandled error occurred")

				// Respond with an error message
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":    "Internal Server Error",
					"details":  err,
					"file":     file,
					"line":     line,
					"function": functionName,
				})
				c.Abort()
				return
			}
		}()

		// Process the request
		c.Next()

		// Log response data
		duration := time.Since(start)
		statusCode := c.Writer.Status()

		logEntry := logrus.WithFields(logrus.Fields{
			"status":     statusCode,
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"query":      c.Request.URL.RawQuery,
			"request":    string(bodyBytes),
			"response":   writer.body.String(),
			"duration":   duration,
			"ip":         c.ClientIP(),
			"user_agent": c.Request.UserAgent(),
		})

		// Log based on status code
		if statusCode >= 400 {
			logEntry.Error("Request completed with error")
		} else {
			logEntry.Info("Request completed successfully")
		}
	}
}
