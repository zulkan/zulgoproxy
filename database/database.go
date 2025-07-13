package database

import (
	"fmt"
	"time"

	"github.com/zulkan/zulgoproxy/auth"
	"github.com/zulkan/zulgoproxy/config"
	"github.com/zulkan/zulgoproxy/logger"
	"github.com/zulkan/zulgoproxy/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDatabase(cfg *config.Config) error {
	var err error
	
	gormConfig := &gorm.Config{}
	
	if cfg.Server.LogLevel == "debug" {
		gormConfig.Logger = gormLogger.Default.LogMode(gormLogger.Info)
	}
	
	DB, err = gorm.Open(postgres.Open(cfg.Database.DSN()), gormConfig)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	
	// Configure connection pool
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}
	
	// Set maximum number of open connections
	sqlDB.SetMaxOpenConns(25)
	
	// Set maximum number of idle connections
	sqlDB.SetMaxIdleConns(5)
	
	// Set maximum lifetime of connections
	sqlDB.SetConnMaxLifetime(5 * time.Minute)
	
	// Auto-migrate the schema
	err = DB.AutoMigrate(&models.User{}, &models.Session{}, &models.ProxyLog{})
	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}
	
	// Create default admin user if it doesn't exist
	if err := createDefaultAdmin(); err != nil {
		logger.Warn("Failed to create default admin: %v", err)
	}
	
	logger.Info("Database connection established and migrations completed")
	return nil
}

func createDefaultAdmin() error {
	var count int64
	if err := DB.Model(&models.User{}).Where("role = ?", models.RoleAdmin).Count(&count).Error; err != nil {
		return err
	}
	
	logger.Debug("Found %d admin users in database", count)
	
	if count == 0 {
		logger.Info("No admin users found, creating default admin user")
		hashedPassword, err := auth.HashPassword("admin")
		if err != nil {
			return err
		}
		
		adminUser := models.User{
			Username: "admin",
			Password: hashedPassword,
			Email:    "admin@zulgoproxy.local",
			Role:     models.RoleAdmin,
			IsActive: true,
		}
		
		if err := DB.Create(&adminUser).Error; err != nil {
			logger.Error("Failed to create admin user: %v", err)
			return err
		}
		
		logger.Info("Default admin user created (username: admin, password: admin)")
	} else {
		logger.Debug("Admin user already exists, skipping creation")
	}
	
	return nil
}

func GetDB() *gorm.DB {
	return DB
}