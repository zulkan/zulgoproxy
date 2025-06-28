# Build stage
FROM golang:1.18-alpine AS builder

# Set the working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o zulgoproxy ./app

# Final stage
FROM alpine:latest

# Install CA certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Set the working directory
WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/zulgoproxy .

# Expose the port the app runs on
EXPOSE 8181

# Command to run the executable
CMD ["./zulgoproxy"]
