#!/bin/bash

echo "🚀 CreatorSync Production Deployment Fix"
echo "========================================"

# Make this script executable
chmod +x deploy-fix.sh

# Build the backend
echo "📦 Building backend..."
cd backend
go build -o creatorsync ./cmd/main.go

if [ $? -eq 0 ]; then
    echo "✅ Backend build successful"
else
    echo "❌ Backend build failed"
    exit 1
fi

echo ""
echo "🔧 Updated features:"
echo "- Manual JWT parsing (more reliable than Clerk SDK)"
echo "- Better error logging and debugging"
echo "- Fallback to Clerk SDK if manual parsing fails"
echo "- Comprehensive token validation"
echo ""

echo "📋 Deployment Checklist:"
echo "1. ✅ Backend compiled successfully"
echo "2. ⏳ Deploy the new backend binary"
echo "3. ⏳ Ensure environment variables are correct:"
echo "   - CLERK_SECRET_KEY=sk_live_... (NOT sk_test_)"
echo "   - APP_ENV=production"
echo "4. ⏳ Check backend logs for:"
echo "   - 'CLERK_SECRET_KEY length: XX, starts with sk_: true'"
echo "   - 'Token received - Length: XXX, JWT structure: true'"
echo "   - 'Successfully authenticated user: user_XXXXX'"
echo ""

echo "🔍 Debug Commands:"
echo "Backend logs: Check for the debug messages above"
echo "Test endpoint: curl -H \"Authorization: Bearer YOUR_TOKEN\" https://api.creatorsync.app/api/twitch/videos"
echo ""

echo "⚡ This fix should resolve the 401 Unauthorized error!"
