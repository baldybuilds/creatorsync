# CreatorSync Analytics - Railway Setup Guide

## Quick Setup for Railway

This guide covers the simplified setup process for deploying CreatorSync on Railway with automatic database configuration.

## 1. Environment Variables Setup

### Railway Database (Automatic)
Railway automatically provides a `DATABASE_URL` environment variable when you add a PostgreSQL service. The application will automatically detect and use this.

**No manual database configuration needed!** ✨

### Required Environment Variables in Railway

Add these environment variables in your Railway service settings:

```bash
# Twitch API Configuration
TWITCH_CLIENT_ID=your_twitch_client_id
TWITCH_CLIENT_SECRET=your_twitch_client_secret

# Clerk Authentication
CLERK_SECRET_KEY=your_clerk_secret_key

# Optional: Admin Access Key
ADMIN_KEY=your_admin_key_for_system_stats
```

### Frontend Environment Variables
For the frontend service, add:

```bash
NEXT_PUBLIC_CLERK_PUBLISHABLE_KEY=your_clerk_publishable_key
CLERK_SECRET_KEY=your_clerk_secret_key
NEXT_PUBLIC_API_URL=https://your-backend-domain.railway.app
```

## 2. Database Setup

### Option 1: Railway PostgreSQL Service (Recommended)
1. Add a PostgreSQL service to your Railway project
2. Railway automatically provides the `DATABASE_URL`
3. The application will auto-configure the connection
4. Run migrations (see below)

### Option 2: External Database
If using an external database, set individual variables:
```bash
POSTGRES_DB_HOST=your_host
POSTGRES_DB_PORT=5432
POSTGRES_DB_DATABASE=your_database_name
POSTGRES_DB_USERNAME=your_username
POSTGRES_DB_PASSWORD=your_password
POSTGRES_DB_SCHEMA=public
```

## 3. Running the Setup Script

### Local Development
```bash
# Clone and setup
git clone your-repo
cd creatorsync

# Set environment variables
export DATABASE_URL="postgresql://username:password@host:port/database"
export TWITCH_CLIENT_ID="your_client_id"
export TWITCH_CLIENT_SECRET="your_client_secret"
export CLERK_SECRET_KEY="your_clerk_secret"

# Run setup script
chmod +x scripts/setup-analytics.sh
./scripts/setup-analytics.sh
```

### Railway Deployment
The setup script will automatically detect Railway's `DATABASE_URL` and configure the connection appropriately.

## 4. Database Migrations

### Automatic Migration (during setup)
The setup script will attempt to run migrations automatically if database connection is successful.

### Manual Migration
If needed, run migrations manually:
```bash
cd backend
go run cmd/migrate/main.go
```

### Railway Deployment Migration
Add a migration build command to your Railway service:
```bash
# In Railway Build Command
go run cmd/migrate/main.go && go build -o main cmd/api/main.go
```

## 5. Twitch Application Setup

1. Go to https://dev.twitch.tv/console
2. Create a new application
3. Set OAuth redirect URL: `https://your-frontend-domain.railway.app/auth/callback`
4. Copy Client ID and Client Secret to Railway environment variables

### Required OAuth Scopes
```
channel:read:subscriptions
user:read:follows  
channel:read:analytics
```

## 6. Clerk Authentication Setup

1. Create account at https://clerk.com
2. Create new application
3. Configure OAuth providers (add Twitch)
4. Copy keys to Railway environment variables
5. Set up webhooks endpoint: `https://your-backend-domain.railway.app/webhooks/clerk`

## 7. Testing the Setup

### Health Check
```bash
curl https://your-backend-domain.railway.app/health
curl https://your-backend-domain.railway.app/api/analytics/health
```

### Test Analytics (with auth token)
```bash
export TOKEN="your_auth_token"
curl -H "Authorization: Bearer $TOKEN" \
     https://your-backend-domain.railway.app/api/analytics/overview
```

## 8. Production Deployment Checklist

- [ ] PostgreSQL service added to Railway project
- [ ] All environment variables configured
- [ ] Twitch application created and configured
- [ ] Clerk application created and configured
- [ ] Database migrations completed
- [ ] Health checks passing
- [ ] Analytics endpoints responding
- [ ] Frontend connecting to backend successfully

## 9. Monitoring and Maintenance

### System Stats (Admin)
```bash
curl -H "X-Admin-Key: $ADMIN_KEY" \
     https://your-backend-domain.railway.app/api/analytics/system
```

### Job Status Monitoring
```bash
curl -H "Authorization: Bearer $TOKEN" \
     https://your-backend-domain.railway.app/api/analytics/jobs
```

### Manual Data Collection
```bash
curl -X POST -H "Authorization: Bearer $TOKEN" \
     https://your-backend-domain.railway.app/api/analytics/trigger
```

## 10. Common Issues and Solutions

### Database Connection Issues
- Verify `DATABASE_URL` is present in Railway environment
- Check PostgreSQL service is running
- Ensure connection pool settings are appropriate

### Authentication Issues  
- Verify Clerk keys are correct
- Check Twitch OAuth configuration
- Ensure webhook endpoints are accessible

### API Rate Limiting
- Monitor Twitch API usage (800 req/min limit)
- Check job scheduler timing
- Adjust batch sizes if needed

### Build Issues
- Ensure all Go dependencies are available
- Check for missing environment variables during build
- Verify migration script can run

## 11. Local Development

For local development with Railway database:

```bash
# Install Railway CLI
npm install -g @railway/cli

# Login and link project
railway login
railway link

# Pull environment variables
railway variables

# Run locally with Railway database
railway run go run cmd/api/main.go
```

## 12. Scaling Considerations

### Database
- Monitor connection pool usage
- Consider read replicas for heavy analytics queries
- Index optimization for large datasets

### Background Jobs
- Monitor job completion rates
- Adjust scheduling for user growth
- Consider distributed job processing

### API Performance
- Implement caching for dashboard data
- Rate limiting for API endpoints
- CDN for static assets

---

For detailed architecture information, see [ANALYTICS_ARCHITECTURE.md](./ANALYTICS_ARCHITECTURE.md) 