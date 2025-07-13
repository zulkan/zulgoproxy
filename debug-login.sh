#!/bin/bash

echo "Testing ZulgoProxy login..."

# Test if API server is running
echo "1. Checking if API server is accessible..."
curl -s -o /dev/null -w "API Health Check: %{http_code}\n" http://localhost:8182/health

# Test login endpoint
echo "2. Testing login with admin/admin..."
curl -X POST \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin"}' \
  -w "\nHTTP Status: %{http_code}\n" \
  http://localhost:8182/api/auth/login

echo ""
echo "If you see 'Connection refused', the server is not running."
echo "If you see 401, there may be a password issue."
echo "If you see 200, login should work."