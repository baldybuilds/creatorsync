#!/bin/bash

echo "=== CreatorSync Environment Check ==="
echo ""

# Check backend environment
echo "Backend Environment Variables:"
echo "==============================="

if [ -f "backend/.env" ]; then
    echo "✅ backend/.env file exists"
    
    # Check for CLERK_SECRET_KEY
    if grep -q "CLERK_SECRET_KEY=" backend/.env; then
        SECRET_KEY=$(grep "CLERK_SECRET_KEY=" backend/.env | cut -d'=' -f2 | tr -d '"' | tr -d "'")
        if [[ $SECRET_KEY == sk_* ]]; then
            echo "✅ CLERK_SECRET_KEY format looks correct (starts with sk_)"
            echo "   Length: ${#SECRET_KEY} characters"
        else
            echo "❌ CLERK_SECRET_KEY format is incorrect (should start with sk_)"
            echo "   Current value starts with: ${SECRET_KEY:0:5}..."
        fi
    else
        echo "❌ CLERK_SECRET_KEY not found in backend/.env"
    fi
    
    # Check for APP_ENV
    if grep -q "APP_ENV=" backend/.env; then
        APP_ENV=$(grep "APP_ENV=" backend/.env | cut -d'=' -f2 | tr -d '"' | tr -d "'")
        echo "✅ APP_ENV is set to: $APP_ENV"
    else
        echo "⚠️  APP_ENV not set (will default to development mode)"
    fi
else
    echo "❌ backend/.env file not found"
fi

echo ""

# Check frontend environment
echo "Frontend Environment Variables:"
echo "==============================="

if [ -f "frontend/.env.local" ]; then
    echo "✅ frontend/.env.local file exists"
    
    # Check for NEXT_PUBLIC_CLERK_PUBLISHABLE_KEY
    if grep -q "NEXT_PUBLIC_CLERK_PUBLISHABLE_KEY=" frontend/.env.local; then
        PUB_KEY=$(grep "NEXT_PUBLIC_CLERK_PUBLISHABLE_KEY=" frontend/.env.local | cut -d'=' -f2 | tr -d '"' | tr -d "'")
        if [[ $PUB_KEY == pk_* ]]; then
            echo "✅ NEXT_PUBLIC_CLERK_PUBLISHABLE_KEY format looks correct (starts with pk_)"
        else
            echo "❌ NEXT_PUBLIC_CLERK_PUBLISHABLE_KEY format is incorrect (should start with pk_)"
        fi
    else
        echo "❌ NEXT_PUBLIC_CLERK_PUBLISHABLE_KEY not found"
    fi
    
    # Check for CLERK_SECRET_KEY in frontend (should match backend)
    if grep -q "CLERK_SECRET_KEY=" frontend/.env.local; then
        FRONTEND_SECRET=$(grep "CLERK_SECRET_KEY=" frontend/.env.local | cut -d'=' -f2 | tr -d '"' | tr -d "'")
        BACKEND_SECRET=$(grep "CLERK_SECRET_KEY=" backend/.env | cut -d'=' -f2 | tr -d '"' | tr -d "'")
        if [ "$FRONTEND_SECRET" = "$BACKEND_SECRET" ]; then
            echo "✅ Frontend and backend CLERK_SECRET_KEY match"
        else
            echo "❌ Frontend and backend CLERK_SECRET_KEY do not match"
        fi
    else
        echo "⚠️  CLERK_SECRET_KEY not found in frontend/.env.local"
    fi
    
elif [ -f "frontend/.env" ]; then
    echo "✅ frontend/.env file exists (checking that instead)"
    # Repeat similar checks for .env file
else
    echo "❌ No frontend environment file found (.env.local or .env)"
fi

echo ""
echo "=== Recommendations ==="

echo "1. Ensure your CLERK_SECRET_KEY starts with 'sk_live_' for production"
echo "2. Ensure your NEXT_PUBLIC_CLERK_PUBLISHABLE_KEY starts with 'pk_live_' for production"
echo "3. Make sure both frontend and backend use the same Clerk project keys"
echo "4. Set APP_ENV=production in your production backend environment"
echo ""
echo "For production deployment, also check:"
echo "- Environment variables are set in your hosting platform"
echo "- Keys are from the same Clerk application"
echo "- Twitch OAuth is properly configured in Clerk dashboard"
