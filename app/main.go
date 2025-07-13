package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/elazarl/goproxy"
	"github.com/elazarl/goproxy/ext/auth"
	"github.com/gin-gonic/gin"
	"github.com/zulkan/zulgoproxy/config"
	"github.com/zulkan/zulgoproxy/database"
	"github.com/zulkan/zulgoproxy/handlers"
	"github.com/zulkan/zulgoproxy/logger"
	"github.com/zulkan/zulgoproxy/middleware"
	"github.com/zulkan/zulgoproxy/ui"
)

var cfg *config.Config

func main() {
	// Load configuration
	var err error
	cfg, err = config.LoadConfig("")
	if err != nil {
		logger.Fatal("Failed to load config: %v", err)
	}

	// Set log level based on config
	logger.SetLevelFromString(cfg.Server.LogLevel)
	logger.Info("ZulgoProxy starting up...")
	logger.Debug("Debug logging enabled")

	// Initialize database
	if err := database.InitDatabase(cfg); err != nil {
		logger.Fatal("Failed to initialize database: %v", err)
	}

	// Setup graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		logger.Info("Shutting down gracefully...")
		os.Exit(0)
	}()

	// Start proxy server
	go startProxyServer()

	// Start API server
	startAPIServer()
}

func startProxyServer() {
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = cfg.Server.LogLevel == "debug"

	proxy.OnRequest().DoFunc(filterIP)
	proxy.OnRequest().HandleConnect(getHandleConnect())

	logger.Info("Proxy server starting on port %d", cfg.Server.Port)
	logger.Fatal("Proxy server error: %v", http.ListenAndServe(fmt.Sprintf(":%d", cfg.Server.Port), proxy))
}

func startAPIServer() {
	if cfg.Server.LogLevel != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// Add rate limiting (100 requests per minute)
	rateLimiter := middleware.NewRateLimiter(100, time.Minute)
	router.Use(middleware.RateLimitMiddleware(rateLimiter))

	// Add logging middleware
	router.Use(middleware.ProxyLoggingMiddleware())

	// Health check endpoints (no auth required)
	healthHandler := handlers.NewHealthHandler()
	router.GET("/health", healthHandler.Health)
	router.GET("/health/readiness", healthHandler.Readiness)
	router.GET("/health/liveness", healthHandler.Liveness)

	// Auth endpoints
	authHandler := handlers.NewAuthHandler(cfg)
	auth := router.Group("/api/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", authHandler.RefreshToken)
		auth.POST("/logout", authHandler.Logout)
		auth.GET("/me", middleware.AuthMiddleware(cfg), authHandler.Me)
	}

	// Protected API routes
	api := router.Group("/api")
	api.Use(middleware.AuthMiddleware(cfg))
	{
		// User management (admin only)
		userHandler := handlers.NewUserHandler()
		users := api.Group("/users")
		users.Use(middleware.AdminMiddleware())
		{
			users.GET("", userHandler.GetUsers)
			users.GET("/:id", userHandler.GetUser)
			users.POST("", userHandler.CreateUser)
			users.PUT("/:id", userHandler.UpdateUser)
			users.DELETE("/:id", userHandler.DeleteUser)
		}

		// Change password (for authenticated users)
		api.POST("/change-password", userHandler.ChangePassword)

		// Logs (admin only)
		logHandler := handlers.NewLogHandler()
		logs := api.Group("/logs")
		logs.Use(middleware.AdminMiddleware())
		{
			logs.GET("", logHandler.GetLogs)
			logs.GET("/stats", logHandler.GetLogStats)
		}

		// Admin endpoints
		adminHandler := handlers.NewAdminHandler()
		admin := api.Group("/admin")
		admin.Use(middleware.AdminMiddleware())
		{
			admin.GET("/dashboard", adminHandler.GetDashboard)
			admin.GET("/system", adminHandler.GetSystemInfo)
			admin.DELETE("/logs/purge", adminHandler.PurgeOldLogs)
		}
	}

	// Serve UI
	ui.AddRoutes(router)

	apiPort := cfg.Server.Port + 1
	logger.Info("API server starting on port %d", apiPort)
	logger.Fatal("API server error: %v", http.ListenAndServe(fmt.Sprintf(":%d", apiPort), router))
}

func filterIP(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	logger.Debug("Request from %s to %s", req.RemoteAddr, req.URL.String())
	return req, nil
}

func isIPAllowed(remoteAddr string) bool {
	clientIP, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		return false
	}

	for _, allowedRange := range cfg.Server.AllowedIPs {
		if strings.Contains(allowedRange, "/") {
			// CIDR notation
			_, subnet, err := net.ParseCIDR(allowedRange)
			if err != nil {
				continue
			}
			if subnet.Contains(net.ParseIP(clientIP)) {
				return true
			}
		} else {
			// Direct IP comparison
			if clientIP == allowedRange {
				return true
			}
		}
	}

	return false
}

func zulUserPass(user, passwd string) bool {
	// TODO: Replace with database authentication
	if user == "zulkan" && passwd == "zulkan" {
		return true
	}
	return false
}

func getHandleConnect() goproxy.HttpsHandler {
	return goproxy.FuncHttpsHandler(func(host string, ctx *goproxy.ProxyCtx) (*goproxy.ConnectAction, string) {
		logger.Debug("CONNECT request to %s from %s", host, ctx.Req.RemoteAddr)

		if !isIPAllowed(ctx.Req.RemoteAddr) {
			return auth.BasicConnect("ZulgoProxy", zulUserPass).HandleConnect(host, ctx)
		}
		return goproxy.OkConnect, host
	})
}
