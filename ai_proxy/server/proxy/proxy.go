package proxy

import (
	"ai-proxy/server/proxy/config"
	"ai-proxy/server/proxy/handlers"
	"ai-proxy/server/proxy/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	config.InitLogger()
	defer config.Logger.Sync()

	router := gin.Default()

	// Apply rate limiting middleware
	router.Use(middleware.RateLimiterMiddleware())

	// Define routes
	router.POST("/proxy", handlers.ProxyRequest)

	router.Run(":8080")
}
