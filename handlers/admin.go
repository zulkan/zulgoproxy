package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zulkan/zulgoproxy/database"
	"github.com/zulkan/zulgoproxy/models"
)

type AdminHandler struct{}

func NewAdminHandler() *AdminHandler {
	return &AdminHandler{}
}

type DashboardStats struct {
	TotalUsers    int64 `json:"total_users"`
	ActiveUsers   int64 `json:"active_users"`
	TotalRequests int64 `json:"total_requests"`
	TodayRequests int64 `json:"today_requests"`
}

func (h *AdminHandler) GetDashboard(c *gin.Context) {
	var stats DashboardStats
	
	// Total users
	database.GetDB().Model(&models.User{}).Count(&stats.TotalUsers)
	
	// Active users
	database.GetDB().Model(&models.User{}).Where("is_active = ?", true).Count(&stats.ActiveUsers)
	
	// Total requests
	database.GetDB().Model(&models.ProxyLog{}).Count(&stats.TotalRequests)
	
	// Today's requests
	database.GetDB().Model(&models.ProxyLog{}).
		Where("DATE(timestamp) = CURRENT_DATE").
		Count(&stats.TodayRequests)
	
	c.JSON(http.StatusOK, gin.H{
		"stats": stats,
	})
}

func (h *AdminHandler) GetSystemInfo(c *gin.Context) {
	// Get database info
	var dbVersion string
	database.GetDB().Raw("SELECT version()").Scan(&dbVersion)
	
	var dbSize string
	database.GetDB().Raw("SELECT pg_size_pretty(pg_database_size(current_database()))").Scan(&dbSize)
	
	// Get table sizes
	var tableSizes []struct {
		TableName string `json:"table_name"`
		Size      string `json:"size"`
	}
	database.GetDB().Raw(`
		SELECT 
			schemaname||'.'||tablename as table_name,
			pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) as size
		FROM pg_tables 
		WHERE schemaname = 'public'
		ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC
	`).Scan(&tableSizes)
	
	c.JSON(http.StatusOK, gin.H{
		"database": gin.H{
			"version":     dbVersion,
			"size":        dbSize,
			"table_sizes": tableSizes,
		},
	})
}

func (h *AdminHandler) PurgeOldLogs(c *gin.Context) {
	days := c.DefaultQuery("days", "30")
	
	var deletedCount int64
	result := database.GetDB().
		Where("timestamp < NOW() - INTERVAL ? day", days).
		Delete(&models.ProxyLog{})
	
	deletedCount = result.RowsAffected
	
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to purge logs"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message":       "Logs purged successfully",
		"deleted_count": deletedCount,
		"days":          days,
	})
}