package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zulkan/zulgoproxy/auth"
	"github.com/zulkan/zulgoproxy/config"
	"github.com/zulkan/zulgoproxy/database"
	"github.com/zulkan/zulgoproxy/logger"
	"github.com/zulkan/zulgoproxy/models"
)

type AuthHandler struct {
	cfg *config.Config
}

func NewAuthHandler(cfg *config.Config) *AuthHandler {
	return &AuthHandler{cfg: cfg}
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	User  *models.User     `json:"user"`
	Token *auth.TokenPair  `json:"token"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("Login: Failed to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	logger.Info("Login attempt for username: %s", req.Username)
	
	var user models.User
	if err := database.GetDB().Where("username = ? AND is_active = ?", req.Username, true).First(&user).Error; err != nil {
		logger.Warn("Login: User not found or inactive: %s, error: %v", req.Username, err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}
	
	logger.Debug("Login: User found: %s, checking password", user.Username)
	
	if !auth.CheckPassword(req.Password, user.Password) {
		logger.Warn("Login: Password check failed for user: %s", req.Username)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}
	
	logger.Info("Login: Password check passed for user: %s", req.Username)
	
	tokens, err := auth.GenerateTokenPair(&user, h.cfg.Auth.JWTSecret, h.cfg.Auth.TokenExpiry, h.cfg.Auth.RefreshExpiry)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}
	
	// Store refresh token in database
	session := models.Session{
		UserID:       user.ID,
		RefreshToken: tokens.RefreshToken,
		ExpiresAt:    time.Unix(tokens.ExpiresAt, 0).Add(time.Duration(h.cfg.Auth.RefreshExpiry-h.cfg.Auth.TokenExpiry) * time.Hour),
	}
	
	if err := database.GetDB().Create(&session).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
		return
	}
	
	user.Password = "" // Don't return password
	c.JSON(http.StatusOK, LoginResponse{
		User:  &user,
		Token: tokens,
	})
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Validate refresh token
	var session models.Session
	if err := database.GetDB().Where("refresh_token = ? AND expires_at > ?", req.RefreshToken, time.Now()).First(&session).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}
	
	// Get user
	var user models.User
	if err := database.GetDB().First(&user, session.UserID).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}
	
	if !user.IsActive {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User account is inactive"})
		return
	}
	
	// Generate new access token
	newAccessToken, err := auth.RefreshAccessToken(req.RefreshToken, h.cfg.Auth.JWTSecret, h.cfg.Auth.TokenExpiry)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to refresh token"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"access_token": newAccessToken,
		"expires_at":   time.Now().Add(time.Duration(h.cfg.Auth.TokenExpiry) * time.Hour).Unix(),
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	// Get refresh token from request
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Delete session
	database.GetDB().Where("refresh_token = ?", req.RefreshToken).Delete(&models.Session{})
	
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func (h *AuthHandler) Me(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"user": user})
}