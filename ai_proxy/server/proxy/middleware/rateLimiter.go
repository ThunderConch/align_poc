package middleware

import (
	"ai-proxy/server/proxy/models"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var rateLimiters = make(map[string]*RateLimiter)
var mutex = &sync.Mutex{}

type RateLimiter struct {
	requests int
	bytes    int64
	resetAt  time.Time
}

func RateLimiterMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("API-Key")
		if apiKey == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "API key is required"})
			return
		}

		apiKeyInfo, err := models.GetAPIKeyInfo(apiKey)
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": "Internal Server Error"})
			return
		}

		mutex.Lock()
		if limiter, exists := rateLimiters[apiKey]; exists {
			if time.Now().After(limiter.resetAt) {
				limiter.requests = 0
				limiter.bytes = 0
				limiter.resetAt = time.Now().Add(time.Minute)
			}
			limiter.requests++
			limiter.bytes += c.Request.ContentLength
			if limiter.requests > apiKeyInfo.RPM || limiter.bytes > apiKeyInfo.BPM {
				mutex.Unlock()
				c.AbortWithStatus(429)
				return
			}
		} else {
			rateLimiters[apiKey] = &RateLimiter{requests: 1, bytes: c.Request.ContentLength, resetAt: time.Now().Add(time.Minute)}
		}
		mutex.Unlock()
		c.Next()
	}
}
