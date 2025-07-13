# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

ZulgoProxy is a modern HTTP/HTTPS proxy server implemented in Go with comprehensive authentication, authorization, logging, and management capabilities. The application consists of a proxy server and a REST API with web UI.

## Architecture

### Core Components
- **app/main.go**: Main application entry point with dual server setup (proxy + API)
- **config/**: Configuration management with YAML support and environment variables
- **database/**: PostgreSQL database integration with GORM and connection pooling
- **models/**: Database models for users, sessions, and proxy logs
- **auth/**: JWT-based authentication and password hashing utilities
- **middleware/**: Authentication, authorization, logging, and rate limiting middleware
- **handlers/**: REST API handlers for auth, users, logs, admin, and health endpoints
- **ui/ui.go**: Embedded web UI serving static files

### Database Schema
- **users**: User accounts with role-based access (admin/user)
- **sessions**: JWT refresh token management
- **proxy_logs**: Detailed request/response logging with metrics

### Server Architecture
- **Proxy Server** (port 8181): HTTP/HTTPS proxy with IP filtering and authentication
- **API Server** (port 8182): REST API with JWT authentication and admin endpoints
- **Database**: PostgreSQL with connection pooling and automatic migrations

## Build Commands

### Prerequisites
```bash
# Install PostgreSQL and create database
createdb zulgoproxy

# Install dependencies
go mod download
```

### Local Development
```bash
# Build the application
go build -o app/app ./app

# Run with default config
./app/app

# Run with custom config
CONFIG_PATH=/path/to/config.yaml ./app/app

# Run directly with go
go run ./app/main.go
```

### Database Setup
```bash
# The application automatically:
# - Runs database migrations on startup
# - Creates default admin user (username: admin, password: admin)
# - Sets up connection pooling (25 max connections, 5 idle)
```

### Docker Build
```bash
# Build Docker image
docker build -t zulgoproxy .

# Run with PostgreSQL
docker run -p 8181:8181 -p 8182:8182 \
  -e DB_HOST=host.docker.internal \
  -e DB_USER=postgres \
  -e DB_PASSWORD=postgres \
  -e DB_NAME=zulgoproxy \
  zulgoproxy
```

### Configuration
- Copy config.yaml.example to config.yaml and modify settings
- Environment variables override config file values
- Required: DB_HOST, DB_USER, DB_PASSWORD, DB_NAME
- Optional: JWT_SECRET (auto-generated if not set)

## API Endpoints

### Authentication
- POST `/api/auth/login` - User login
- POST `/api/auth/refresh` - Refresh access token
- POST `/api/auth/logout` - User logout
- GET `/api/auth/me` - Get current user info

### User Management (Admin only)
- GET `/api/users` - List users with pagination
- GET `/api/users/:id` - Get specific user
- POST `/api/users` - Create new user
- PUT `/api/users/:id` - Update user
- DELETE `/api/users/:id` - Delete user

### Logging & Analytics (Admin only)
- GET `/api/logs` - Get proxy logs with filtering
- GET `/api/logs/stats` - Get usage statistics

### Admin Dashboard (Admin only)
- GET `/api/admin/dashboard` - Dashboard statistics
- GET `/api/admin/system` - System information
- DELETE `/api/admin/logs/purge?days=30` - Purge old logs

### Health Checks
- GET `/health` - Overall health status
- GET `/health/readiness` - Readiness probe
- GET `/health/liveness` - Liveness probe

## Development Notes

### Server Configuration
- Proxy server: port 8181 (configurable)
- API server: port 8182 (proxy port + 1)
- Rate limiting: 100 requests/minute per user/IP
- JWT tokens: 24h access, 7d refresh (configurable)

### Security Features
- JWT-based authentication for API
- Role-based access control (admin/user)
- IP whitelist support with CIDR notation
- Rate limiting per user/IP
- Password hashing with bcrypt
- Request/response logging with user tracking

### Database Features
- Automatic migrations on startup
- Connection pooling (25 max, 5 idle, 5min lifetime)
- Soft deletes for user accounts
- Detailed proxy request logging
- Session management for refresh tokens