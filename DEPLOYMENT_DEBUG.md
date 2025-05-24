# CreatorSync Deployment Debug Guide

This guide helps diagnose and fix data fetching issues between local, staging, and production environments.

## Current Issue Summary

- **Local Development**: ✅ Working perfectly
- **Staging**: ⚠️ Partial data loading (Followers & Subscribers only)
- **Production**: ❌ No data loading (500 Internal Server Error)

## Quick Diagnosis Steps

### 1. Check Backend Logs

Run the debug tool to verify environment configuration:

```bash
cd backend
go run cmd/debug/main.go
```

This will check:
- Environment variables
- Database connectivity
- Clerk authentication setup
- Required database tables

### 2. Check Frontend Console

Open browser dev tools on staging/production and look for:
- API endpoint URLs being called
- Authentication token presence
- Detailed error messages
- Debug information from failed requests

### 3. Test API Endpoints Directly

Test the health endpoint:
```bash
# Staging
curl https://api-dev.creatorsync.app/health

# Production  
curl https://api.creatorsync.app/health
```

Test the analytics debug endpoint (requires auth token):
```bash
# Get your auth token from browser dev tools, then:
curl -H "Authorization: Bearer YOUR_TOKEN" \
     https://api-dev.creatorsync.app/api/analytics/debug/data-status
```

## Common Issues & Solutions

### Issue 1: Database Connection Problems

**Symptoms**: 500 errors, "database connection failed" in logs

**Solutions**:
1. Verify `DATABASE_URL` environment variable is set correctly
2. Check database server is accessible from deployment environment
3. Verify SSL settings (`sslmode=require` for production databases)

### Issue 2: Missing Environment Variables

**Symptoms**: "NOT CONFIGURED" messages in startup logs

**Solutions**:
1. Verify all required environment variables are set:
   - `APP_ENV` (staging/production)
   - `DATABASE_URL` or individual DB variables
   - `CLERK_SECRET_KEY`
   - `TWITCH_CLIENT_ID` and `TWITCH_CLIENT_SECRET`
   - `PORT`

### Issue 3: Cross-Environment User Issues

**Symptoms**: Users exist in Clerk but not in database

**Solutions**:
1. The updated code now automatically creates user records
2. Check logs for "Created basic user record" messages
3. Verify user sync endpoint is working: `/api/user/sync`

### Issue 4: CORS Issues

**Symptoms**: CORS errors in browser console

**Solutions**:
1. Verify `APP_ENV` is set correctly (staging/production)
2. Check CORS configuration in `backend/internal/server/routes.go`
3. Ensure frontend domains match CORS allowed origins

### Issue 5: Authentication Problems

**Symptoms**: "User not authenticated" errors

**Solutions**:
1. Verify Clerk publishable key matches environment
2. Check Clerk secret key is correct for the environment
3. Ensure JWT tokens are being passed correctly

## Environment-Specific Checks

### Staging Environment

**Expected Configuration**:
- `APP_ENV=staging`
- Frontend: `dev.creatorsync.app`
- Backend: `api-dev.creatorsync.app`
- Database: Production database (shared)
- Clerk: Production Clerk environment

**Verification Steps**:
```bash
# Check if staging backend is responding
curl https://api-dev.creatorsync.app/health

# Check environment configuration
curl https://api-dev.creatorsync.app/ 
# Should return: {"message": "The Server is running!!"}
```

### Production Environment

**Expected Configuration**:
- `APP_ENV=production`
- Frontend: `creatorsync.app`
- Backend: `api.creatorsync.app`
- Database: Production database
- Clerk: Production Clerk environment

**Verification Steps**:
```bash
# Check if production backend is responding
curl https://api.creatorsync.app/health

# Check environment configuration
curl https://api.creatorsync.app/
# Should return: {"message": "The Server is running!!"}
```

## Debugging Commands

### Backend Debug Tool
```bash
cd backend
go run cmd/debug/main.go
```

### Check Database Tables
```bash
cd backend
go run cmd/migrate/main.go
```

### View Backend Logs
Check your deployment platform (Railway, Vercel, etc.) for backend logs showing:
- Server startup messages
- Environment configuration
- Database connection status
- API request logs

### Frontend Debug Information
Open browser dev tools and check console for:
- API base URL being used
- Environment detection
- Authentication token presence
- Detailed error responses

## Quick Fixes

### 1. Force User Sync
Visit the analytics page and check browser console. The updated code will automatically:
- Detect missing users
- Create basic user records
- Log the process

### 2. Trigger Manual Data Collection
Use the "Collect Data" button on the analytics page to manually trigger data collection.

### 3. Check Debug Endpoint
Visit `/api/analytics/debug/data-status` (requires authentication) to see:
- User existence status
- Database health
- Analytics data availability
- Environment information

## Next Steps

1. **Deploy the updated code** with improved error handling and logging
2. **Check the logs** after deployment to see detailed error information
3. **Test the debug endpoint** to verify system status
4. **Monitor the console** on frontend for detailed debugging information

The updated code includes:
- ✅ Automatic user creation for cross-environment compatibility
- ✅ Better error handling and logging
- ✅ Database connection health checks
- ✅ Detailed debug information
- ✅ Environment configuration logging

This should resolve the data fetching issues between environments. 