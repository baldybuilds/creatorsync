# Environment Variables for Production

## Backend (.env or environment variables)
```bash
# Clerk Configuration
CLERK_SECRET_KEY=sk_live_xxxxxxxxxx  # Your Clerk secret key
APP_ENV=production

# Database
DATABASE_URL=your_database_url

# Other variables as needed
```

## Frontend (.env.local or environment variables)
```bash
# Clerk Configuration
NEXT_PUBLIC_CLERK_PUBLISHABLE_KEY=pk_live_xxxxxxxxxx  # Your Clerk publishable key
CLERK_SECRET_KEY=sk_live_xxxxxxxxxx  # Same as backend

# API Configuration
NEXT_PUBLIC_API_URL=https://api.creatorsync.app
```

## Quick Check Commands

### Verify Backend Environment
```bash
cd backend
echo $CLERK_SECRET_KEY
# Should output: sk_live_...
```

### Verify Frontend Environment
```bash
cd frontend
echo $NEXT_PUBLIC_CLERK_PUBLISHABLE_KEY
# Should output: pk_live_...
```

### Run Environment Check Script
```bash
chmod +x check-env.sh
./check-env.sh
```

## Common Issues

1. **Wrong Key Format**: Make sure you're using `sk_live_` for secret key and `pk_live_` for publishable key
2. **Mismatched Keys**: Frontend and backend must use keys from the same Clerk application
3. **Environment Mismatch**: Don't mix test and live keys
4. **Missing Variables**: Ensure all required environment variables are set in your deployment platform

## Clerk Dashboard Checklist

1. Go to your Clerk dashboard
2. Navigate to "Developers" → "API Keys"
3. Copy the **Secret Key** (starts with `sk_live_`) for backend
4. Copy the **Publishable Key** (starts with `pk_live_`) for frontend
5. Make sure you're using keys from the **production** instance, not development

## Twitch OAuth Configuration

1. In Clerk dashboard, go to "User & Authentication" → "Social Connections"
2. Configure Twitch OAuth with proper scopes:
   - `user:read:email`
   - `clips:edit` (if needed)
   - `channel:read:videos` (for VOD access)
3. Ensure the redirect URLs match your production domain
