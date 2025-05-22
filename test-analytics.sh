#!/bin/bash

# Test analytics endpoints
echo "Testing analytics system..."

API_BASE="http://localhost:8080"

echo "1. Testing health endpoint..."
curl -s "$API_BASE/api/analytics/health" | jq .

echo -e "\n2. Testing overview endpoint (requires auth)..."
echo "Please set TOKEN environment variable with your auth token"
if [ ! -z "$TOKEN" ]; then
    curl -s -H "Authorization: Bearer $TOKEN" "$API_BASE/api/analytics/overview" | jq .
else
    echo "Skipping - no TOKEN set"
fi

echo -e "\nAnalytics test completed"
