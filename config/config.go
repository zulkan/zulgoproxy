package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Database DatabaseConfig `yaml:"database"`
	Server   ServerConfig   `yaml:"server"`
	Auth     AuthConfig     `yaml:"auth"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode"`
}

type ServerConfig struct {
	Port         int      `yaml:"port"`
	AllowedIPs   []string `yaml:"allowed_ips"`
	LogLevel     string   `yaml:"log_level"`
	EnableHTTPS  bool     `yaml:"enable_https"`
	CertFile     string   `yaml:"cert_file"`
	KeyFile      string   `yaml:"key_file"`
}

type AuthConfig struct {
	JWTSecret     string `yaml:"jwt_secret"`
	TokenExpiry   int    `yaml:"token_expiry"` // in hours
	RefreshExpiry int    `yaml:"refresh_expiry"` // in hours
}

func LoadConfig(configPath string) (*Config, error) {
	config := &Config{}
	
	// Set defaults
	config.Server.Port = 8181
	config.Server.LogLevel = "info"
	config.Auth.TokenExpiry = 24
	config.Auth.RefreshExpiry = 168 // 7 days
	config.Database.SSLMode = "disable"
	
	if configPath == "" {
		configPath = "config.yaml"
	}
	
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return config, nil // Return defaults if config file doesn't exist
	}
	
	file, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}
	
	err = yaml.Unmarshal(file, config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}
	
	// Override with environment variables
	if dbHost := os.Getenv("DB_HOST"); dbHost != "" {
		config.Database.Host = dbHost
	}
	if dbUser := os.Getenv("DB_USER"); dbUser != "" {
		config.Database.User = dbUser
	}
	if dbPassword := os.Getenv("DB_PASSWORD"); dbPassword != "" {
		config.Database.Password = dbPassword
	}
	if dbName := os.Getenv("DB_NAME"); dbName != "" {
		config.Database.DBName = dbName
	}
	if jwtSecret := os.Getenv("JWT_SECRET"); jwtSecret != "" {
		config.Auth.JWTSecret = jwtSecret
	}
	
	return config, nil
}

func (d *DatabaseConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		d.Host, d.Port, d.User, d.Password, d.DBName, d.SSLMode)
}