#!/bin/bash

echo "🔍 CreatorSync Production Debugging"
echo "===================================="

API_BASE="https://api.creatorsync.app"

echo ""
echo "1. Testing basic health check..."
curl -s "${API_BASE}/health" | jq . 2>/dev/null || echo "❌ Health check failed or returned non-JSON"

echo ""
echo "2. Testing analytics health check..."
curl -s "${API_BASE}/api/analytics/health" | jq . 2>/dev/null || echo "❌ Analytics health check failed or returned non-JSON"

echo ""
echo "3. Testing root endpoint..."
curl -s "${API_BASE}/" | jq . 2>/dev/null || echo "❌ Root endpoint failed or returned non-JSON"

echo ""
echo "4. Testing with authentication (you'll need to provide a token)..."
echo "To test authenticated endpoints, run:"
echo ""
echo "export TOKEN=\"your_clerk_jwt_token_here\""
echo "curl -H \"Authorization: Bearer \$TOKEN\" ${API_BASE}/api/analytics/debug/data-status"
echo "curl -H \"Authorization: Bearer \$TOKEN\" ${API_BASE}/api/analytics/enhanced?days=30"

echo ""
echo "📋 Common 500 error causes:"
echo "1. Missing environment variables (CLERK_SECRET_KEY, DATABASE_URL, etc.)"
echo "2. Database connection issues"
echo "3. Database tables don't exist (need migrations)"
echo "4. Twitch API credentials missing/invalid"
echo "5. Authentication middleware failing"

echo ""
echo "🔧 To check Railway logs:"
echo "1. Go to Railway dashboard"
echo "2. Select your backend service"
echo "3. Click on 'Logs' tab"
echo "4. Look for error messages when accessing /api/analytics/enhanced"

echo ""
echo "💡 Debug environment variables in Railway:"
echo "Required variables:"
echo "- DATABASE_URL (auto-provided by Railway PostgreSQL)"
echo "- CLERK_SECRET_KEY (should start with sk_live_ for production)"
echo "- TWITCH_CLIENT_ID"
echo "- TWITCH_CLIENT_SECRET"
echo "- APP_ENV=production (recommended)" 