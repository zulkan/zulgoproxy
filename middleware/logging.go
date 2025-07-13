package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zulkan/zulgoproxy/database"
	"github.com/zulkan/zulgoproxy/models"
)

func ProxyLoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		
		// Process request
		c.Next()
		
		// Log after processing
		duration := time.Since(start)
		
		// Get user ID if authenticated
		var userID *uint
		if id, exists := c.Get("user_id"); exists {
			if uid, ok := id.(uint); ok {
				userID = &uid
			}
		}
		
		// Create log entry
		logEntry := models.ProxyLog{
			UserID:       userID,
			RemoteAddr:   c.ClientIP(),
			Method:       c.Request.Method,
			URL:          c.Request.URL.String(),
			Host:         c.Request.Host,
			UserAgent:    c.Request.UserAgent(),
			StatusCode:   c.Writer.Status(),
			ResponseSize: int64(c.Writer.Size()),
			Duration:     duration.Milliseconds(),
			Timestamp:    start,
		}
		
		// Save to database (async to avoid blocking)
		go func() {
			database.GetDB().Create(&logEntry)
		}()
	}
}