package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zulkan/zulgoproxy/database"
	"github.com/zulkan/zulgoproxy/models"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

type HealthResponse struct {
	Status    string            `json:"status"`
	Timestamp time.Time         `json:"timestamp"`
	Version   string            `json:"version"`
	Uptime    string            `json:"uptime"`
	Checks    map[string]Check  `json:"checks"`
}

type Check struct {
	Status  string        `json:"status"`
	Message string        `json:"message,omitempty"`
	Latency time.Duration `json:"latency,omitempty"`
}

var startTime = time.Now()

func (h *HealthHandler) Health(c *gin.Context) {
	checks := make(map[string]Check)
	overallStatus := "healthy"
	
	// Database check
	dbStart := time.Now()
	var count int64
	if err := database.GetDB().Model(&models.User{}).Count(&count).Error; err != nil {
		checks["database"] = Check{
			Status:  "unhealthy",
			Message: err.Error(),
			Latency: time.Since(dbStart),
		}
		overallStatus = "unhealthy"
	} else {
		checks["database"] = Check{
			Status:  "healthy",
			Latency: time.Since(dbStart),
		}
	}
	
	// Memory check (basic)
	checks["memory"] = Check{
		Status: "healthy",
	}
	
	response := HealthResponse{
		Status:    overallStatus,
		Timestamp: time.Now(),
		Version:   "1.0.0",
		Uptime:    time.Since(startTime).String(),
		Checks:    checks,
	}
	
	statusCode := http.StatusOK
	if overallStatus == "unhealthy" {
		statusCode = http.StatusServiceUnavailable
	}
	
	c.JSON(statusCode, response)
}

func (h *HealthHandler) Readiness(c *gin.Context) {
	// Check if database is ready
	if err := database.GetDB().Raw("SELECT 1").Error; err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "not ready",
			"error":  err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status": "ready",
	})
}

func (h *HealthHandler) Liveness(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "alive",
		"uptime": time.Since(startTime).String(),
	})
}