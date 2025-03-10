package router

import (
	"github.com/gin-gonic/gin"
	"grafana-matrix-alerts/logger"
	"time"
)

// Custom middleware to log requests
func zerologMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Capture the start time for the request
		start := time.Now()

		// Log the incoming request
		logger.Log.Info().
			Str("method", c.Request.Method).
			Str("path", c.Request.URL.Path).
			Msg("Request received")

		// Process the request
		c.Next()

		// Log the response time
		duration := time.Since(start)
		logger.Log.Info().
			Str("method", c.Request.Method).
			Str("path", c.Request.URL.Path).
			Int("status", c.Writer.Status()).
			Dur("duration", duration).
			Msg("Request processed")
	}
}
