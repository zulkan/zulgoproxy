package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zulkan/zulgoproxy/database"
	"github.com/zulkan/zulgoproxy/models"
)

type LogHandler struct{}

func NewLogHandler() *LogHandler {
	return &LogHandler{}
}

func (h *LogHandler) GetLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset := (page - 1) * limit
	
	// Filters
	userID := c.Query("user_id")
	method := c.Query("method")
	host := c.Query("host")
	fromDate := c.Query("from_date")
	toDate := c.Query("to_date")
	
	query := database.GetDB().Model(&models.ProxyLog{}).Preload("User")
	
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	if method != "" {
		query = query.Where("method = ?", method)
	}
	if host != "" {
		query = query.Where("host ILIKE ?", "%"+host+"%")
	}
	if fromDate != "" {
		if from, err := time.Parse("2006-01-02", fromDate); err == nil {
			query = query.Where("timestamp >= ?", from)
		}
	}
	if toDate != "" {
		if to, err := time.Parse("2006-01-02", toDate); err == nil {
			query = query.Where("timestamp <= ?", to.Add(24*time.Hour))
		}
	}
	
	var total int64
	query.Count(&total)
	
	var logs []models.ProxyLog
	if err := query.Order("timestamp DESC").Offset(offset).Limit(limit).Find(&logs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch logs"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"logs":  logs,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

func (h *LogHandler) GetLogStats(c *gin.Context) {
	fromDate := c.DefaultQuery("from_date", time.Now().AddDate(0, 0, -7).Format("2006-01-02"))
	toDate := c.DefaultQuery("to_date", time.Now().Format("2006-01-02"))
	
	from, _ := time.Parse("2006-01-02", fromDate)
	to, _ := time.Parse("2006-01-02", toDate)
	to = to.Add(24 * time.Hour)
	
	// Total requests
	var totalRequests int64
	database.GetDB().Model(&models.ProxyLog{}).Where("timestamp BETWEEN ? AND ?", from, to).Count(&totalRequests)
	
	// Requests by method
	var methodStats []struct {
		Method string `json:"method"`
		Count  int64  `json:"count"`
	}
	database.GetDB().Model(&models.ProxyLog{}).
		Select("method, COUNT(*) as count").
		Where("timestamp BETWEEN ? AND ?", from, to).
		Group("method").
		Find(&methodStats)
	
	// Requests by status code
	var statusStats []struct {
		StatusCode int   `json:"status_code"`
		Count      int64 `json:"count"`
	}
	database.GetDB().Model(&models.ProxyLog{}).
		Select("status_code, COUNT(*) as count").
		Where("timestamp BETWEEN ? AND ?", from, to).
		Group("status_code").
		Find(&statusStats)
	
	// Top hosts
	var hostStats []struct {
		Host  string `json:"host"`
		Count int64  `json:"count"`
	}
	database.GetDB().Model(&models.ProxyLog{}).
		Select("host, COUNT(*) as count").
		Where("timestamp BETWEEN ? AND ?", from, to).
		Group("host").
		Order("count DESC").
		Limit(10).
		Find(&hostStats)
	
	// Average response time
	var avgResponseTime float64
	database.GetDB().Model(&models.ProxyLog{}).
		Select("AVG(duration)").
		Where("timestamp BETWEEN ? AND ?", from, to).
		Scan(&avgResponseTime)
	
	c.JSON(http.StatusOK, gin.H{
		"total_requests":      totalRequests,
		"method_stats":        methodStats,
		"status_stats":        statusStats,
		"host_stats":          hostStats,
		"avg_response_time":   avgResponseTime,
		"from_date":           fromDate,
		"to_date":             toDate,
	})
}