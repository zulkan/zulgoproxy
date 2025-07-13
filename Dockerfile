# Node.js build stage for frontend
FROM node:18-alpine AS frontend-builder

WORKDIR /app/ui

# Copy package files
COPY ui/package*.json ./

# Install dependencies
RUN npm ci

# Copy UI source code
COPY ui/ ./

# Build the frontend
RUN npm run build

# Go build stage
FROM golang:1.18-alpine AS backend-builder

# Set the working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Copy built frontend from previous stage
COPY --from=frontend-builder /app/ui/build ./ui/build

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o zulgoproxy ./app

# Final stage
FROM alpine:latest

# Install CA certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Set the working directory
WORKDIR /root/

# Copy the binary from backend builder
COPY --from=backend-builder /app/zulgoproxy .

# Expose the ports the app runs on
EXPOSE 8181 8182

# Command to run the executable
CMD ["./zulgoproxy"]
