package router

import (
	"github.com/gin-gonic/gin"
	"grafana-matrix-alerts/actions"
)

func WebhookEndpoints() *gin.Engine {
	// Set Gin to release mode
	gin.SetMode(gin.ReleaseMode)

	// Create a new Gin router
	r := gin.New()

	// Set default middlewares
	r.Use(zerologMiddleware())
	r.Use(gin.Recovery())

	// Define routes
	r.POST("/api/v1/unified/:room", actions.HandleWebhook)

	return r
}
