#!/bin/bash

echo "=== ZulgoProxy Detailed Debug ==="
echo ""

# Test if API server is running
echo "1. Testing API server health..."
health_response=$(curl -s http://localhost:8182/health)
if [ $? -eq 0 ]; then
    echo "✓ API server is responding"
    echo "Health response: $health_response"
else
    echo "✗ API server is not responding"
    exit 1
fi

echo ""
echo "2. Testing database connection via health check..."
readiness_response=$(curl -s -w "%{http_code}" http://localhost:8182/health/readiness)
echo "Readiness check: $readiness_response"

echo ""
echo "3. Testing login with detailed response..."
login_response=$(curl -X POST \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin"}' \
  -s -w "\nHTTP_CODE:%{http_code}" \
  http://localhost:8182/api/auth/login)

echo "Login response:"
echo "$login_response"

echo ""
echo "4. Testing with different credentials to verify error handling..."
bad_login_response=$(curl -X POST \
  -H "Content-Type: application/json" \
  -d '{"username":"wronguser","password":"wrongpass"}' \
  -s -w "\nHTTP_CODE:%{http_code}" \
  http://localhost:8182/api/auth/login)

echo "Bad login response:"
echo "$bad_login_response"

echo ""
echo "5. Testing login endpoint exists..."
curl -X OPTIONS \
  -H "Content-Type: application/json" \
  -s -w "OPTIONS Status: %{http_code}\n" \
  http://localhost:8182/api/auth/login > /dev/null

echo ""
echo "=== Debug Complete ==="