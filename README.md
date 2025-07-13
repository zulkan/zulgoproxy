# zulgoproxy

HTTP/HTTPS proxy server with authentication and IP filtering capabilities.

## Features

- **HTTP/HTTPS Proxy Server** - High-performance proxy on port 8181
- **Modern Web UI** - React-based admin interface built with Vite
- **JWT Authentication** - Secure token-based authentication system
- **Role-Based Access Control** - Admin and user roles with granular permissions  
- **PostgreSQL Integration** - Robust database backend with connection pooling
- **Real-time Monitoring** - Health checks, metrics, and system monitoring
- **Advanced Logging** - Structured logging with file/line information
- **Rate Limiting** - Per-user/IP request throttling
- **IP Filtering** - CIDR-based access control and whitelisting
- **REST API** - Comprehensive management API
- **Docker Support** - Containerized deployment ready

## Features Status

### Authentication & Security
- [x] JWT-based authentication
- [x] Multiple user accounts with role-based access
- [x] IP whitelist configuration via config file
- [ ] HTTPS certificate management

### Monitoring & Logging
- [x] **Advanced Request/Response Logging** - Detailed logs with timestamps and user tracking
- [x] **Interactive Analytics Dashboard** - Real-time charts and statistics
- [x] **Comprehensive Metrics** - Connection stats, response times, and performance data
- [x] **Health Check Endpoints** - Kubernetes-ready health, readiness, and liveness probes
- [x] **Structured Logging** - File/line information with configurable log levels
- [x] **Log Management** - Automated purging and database optimization

### Configuration & Management  
- [x] **PostgreSQL Database Integration** - Robust data persistence with GORM
- [x] **YAML Configuration** - Flexible config files with environment variable overrides
- [x] **Environment Variable Support** - Docker-friendly configuration
- [x] **Database Connection Pooling** - Optimized connection management (25 max, 5 idle)
- [x] **Comprehensive Admin API** - Full management capabilities via REST API
- [x] **Automatic Database Migrations** - Schema updates handled automatically
- [ ] Runtime configuration reload

### Performance & Reliability
- [x] **Database Connection Pooling** - Optimized PostgreSQL connections
- [x] **Rate Limiting** - Configurable per-user/IP request throttling (100 req/min default)
- [x] **Graceful Shutdown** - Clean shutdown handling with signal management
- [x] **JWT Token Management** - Efficient authentication with refresh tokens
- [x] **Error Handling** - Comprehensive error logging and recovery
- [x] **Performance Monitoring** - Response time tracking and metrics
- [ ] Load balancing support

### UI Enhancements
- [x] **Modern React Frontend** - Built with Vite for optimal performance
- [x] **User Management Interface** - Complete CRUD operations for users
- [x] **Real-time Dashboard** - Live statistics and analytics
- [x] **Traffic Visualization Charts** - Interactive charts using Recharts
- [x] **System Health Monitoring** - Visual health status and metrics
- [x] **Proxy Logs Viewer** - Advanced filtering and search capabilities
- [x] **Responsive Design** - Mobile-friendly interface with Tailwind CSS
- [x] **Role-based UI** - Different views for admin and regular users

## Technology Stack

### Backend
- **Go 1.18+** - High-performance backend with goroutines
- **Gin Framework** - Fast HTTP web framework  
- **GORM** - Go ORM with PostgreSQL driver
- **JWT (golang-jwt/jwt/v5)** - Secure authentication tokens
- **GoProxy** - HTTP/HTTPS proxy implementation
- **bcrypt** - Password hashing and verification

### Frontend  
- **React 18** - Modern UI library with hooks
- **Vite** - Fast build tool with HMR
- **Tailwind CSS** - Utility-first CSS framework
- **React Router v6** - Client-side routing
- **Recharts** - Responsive chart library
- **Lucide React** - Modern icon library
- **Axios** - HTTP client with interceptors

### Database & Infrastructure
- **PostgreSQL 12+** - Robust relational database
- **Docker** - Containerization support
- **YAML** - Configuration management
- **Custom Logger** - Structured logging with file/line info

## API Documentation

### Authentication Endpoints
- `POST /api/auth/login` - User authentication
- `POST /api/auth/refresh` - Refresh access token
- `POST /api/auth/logout` - User logout
- `GET /api/auth/me` - Get current user information

### User Management (Admin Only)
- `GET /api/users` - List all users with pagination
- `GET /api/users/:id` - Get specific user details
- `POST /api/users` - Create new user account
- `PUT /api/users/:id` - Update user information
- `DELETE /api/users/:id` - Delete user account
- `POST /api/change-password` - Change user password

### Logging & Analytics (Admin Only)
- `GET /api/logs` - Get proxy logs with filtering options
- `GET /api/logs/stats` - Get traffic statistics and analytics

### Admin Dashboard (Admin Only)
- `GET /api/admin/dashboard` - Get dashboard statistics
- `GET /api/admin/system` - Get system information
- `DELETE /api/admin/logs/purge` - Purge old log entries

### Health Monitoring
- `GET /health` - Overall application health status
- `GET /health/readiness` - Kubernetes readiness probe
- `GET /health/liveness` - Kubernetes liveness probe

## Setup Instructions

### Prerequisites
- Go 1.18+
- PostgreSQL 12+

### Quick Start

#### Backend Setup
```bash
# 1. Create database
createdb zulgoproxy

# 2. Copy and configure
cp config.yaml.example config.yaml
# Edit config.yaml with your database credentials

# 3. Install Go dependencies
go mod download

# 4. Build and run
go build -o app/app ./app
./app/app

# Or run directly
go run ./app/main.go
```

#### Frontend Development (Optional)
```bash
# Navigate to UI directory
cd ui

# Install Node.js dependencies  
npm install

# Start development server with HMR
npm run dev

# Build for production
npm run build
```

#### Docker Deployment
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

### Default Credentials
- **Username:** admin
- **Password:** admin

### Access Points
- **Proxy Server:** http://localhost:8181 (configurable)
- **API Server:** http://localhost:8182 (proxy port + 1)  
- **Web UI:** http://localhost:8182/ (served by API server)
- **Health Checks:** http://localhost:8182/health

### Configuration
- **Database:** Configure in `config.yaml` or via environment variables
- **Logging:** Set `log_level: debug` for detailed file/line logging
- **JWT Secret:** Change `jwt_secret` in production
- **Rate Limiting:** Default 100 requests/minute per user/IP