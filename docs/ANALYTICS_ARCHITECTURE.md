# CreatorSync Analytics Architecture

## Overview

CreatorSync's analytics system is a comprehensive platform for collecting, processing, and visualizing streamer metrics from Twitch. The system provides real-time insights into follower growth, viewership trends, content performance, and engagement metrics.

## Architecture Components

### 1. Database Schema (`backend/migrations/001_create_analytics_tables.sql`)

**Core Tables:**
- `channel_analytics` - Daily snapshots of channel metrics (followers, views, subscribers)
- `stream_sessions` - Individual stream performance data
- `video_analytics` - VOD, clip, and highlight performance
- `game_analytics` - Performance by game/category
- `analytics_jobs` - Background job tracking and status

**Key Features:**
- Comprehensive indexing for fast queries
- Foreign key relationships for data integrity
- Optimized views for common analytics queries
- Support for cross-platform analytics (future expansion)

### 2. Twitch API Integration (`backend/internal/twitch/`)

**Client Features:**
- Complete Twitch Helix API integration
- OAuth token-based authentication
- Rate limiting and error handling
- Comprehensive data models for all Twitch endpoints

**Data Collection:**
- Channel information and metrics
- Follower counts and subscriber data
- Video analytics (VODs, clips, highlights)
- Live stream information
- User profile data

### 3. Data Collection System (`backend/internal/analytics/`)

**Collector Components:**
- `collector.go` - Handles data collection from Twitch API
- `scheduler.go` - Background job scheduler for automated collection
- `models.go` - Data models and structures
- `repository.go` - Database operations and queries

**Collection Types:**
- Daily channel metrics (followers, views, subscribers)
- Video performance data (views, engagement)
- Stream session analytics
- Game/category performance analysis

### 4. Analytics API (`backend/internal/analytics/handlers.go`)

**Endpoints:**
```
GET  /api/analytics/overview        - Dashboard overview metrics
GET  /api/analytics/charts          - Chart data for visualizations
GET  /api/analytics/detailed        - Comprehensive analytics
GET  /api/analytics/growth          - Growth analysis by period
GET  /api/analytics/content         - Content performance analysis
POST /api/analytics/trigger         - Manual data collection
POST /api/analytics/refresh         - Refresh channel data
GET  /api/analytics/jobs            - Job status and history
GET  /api/analytics/health          - Service health check
```

**Admin Endpoints:**
```
POST /api/analytics/daily           - Trigger daily collection (all users)
GET  /api/analytics/system          - System-wide statistics
```

### 5. Background Job System

**Scheduler Features:**
- Automated daily collection at 2 AM UTC
- Batch processing with rate limiting
- Error handling and retry logic
- Job status tracking and monitoring

**Data Collection Process:**
1. Retrieve users with Twitch integration
2. Process users in batches (10 users per batch)
3. Add jitter to respect API rate limits
4. Track job status and errors
5. Update analytics tables with new data

### 6. Real-time Dashboard (`frontend/src/app/dashboard/analytics/`)

**Visualization Components:**
- Interactive charts using Recharts library
- Real-time metric cards with trend indicators
- Responsive design with dark/light mode support
- Time range selectors and data filtering

**Chart Types:**
- Area charts for follower growth
- Line charts for viewership trends
- Bar charts for stream frequency
- Pie charts for top games/categories
- Performance metrics for videos

**Features:**
- Manual data refresh capabilities
- Data collection triggering
- Export functionality
- Mobile-responsive design

## Setup Instructions

### Prerequisites

1. **Backend Dependencies:**
   ```bash
   cd backend
   go mod download
   ```

2. **Frontend Dependencies:**
   ```bash
   cd frontend
   npm install
   ```

3. **Database Setup:**
   - PostgreSQL 13+ running
   - Environment variables configured
   - Migrations applied

### Environment Configuration

**Backend (.env):**
```env
# Database
POSTGRES_DB_HOST=localhost
POSTGRES_DB_PORT=5432
POSTGRES_DB_DATABASE=creatorsync
POSTGRES_DB_USERNAME=your_username
POSTGRES_DB_PASSWORD=your_password
POSTGRES_DB_SCHEMA=public

# Twitch API
TWITCH_CLIENT_ID=your_twitch_client_id
TWITCH_CLIENT_SECRET=your_twitch_client_secret

# Clerk Authentication
CLERK_SECRET_KEY=your_clerk_secret_key
```

**Frontend (.env):**
```env
NEXT_PUBLIC_CLERK_PUBLISHABLE_KEY=your_clerk_publishable_key
CLERK_SECRET_KEY=your_clerk_secret_key
NEXT_PUBLIC_API_URL=http://localhost:8080
```

### Running the System

1. **Start the Backend:**
   ```bash
   cd backend
   go run cmd/api/main.go
   ```

2. **Start the Frontend:**
   ```bash
   cd frontend
   npm run dev
   ```

3. **Run Database Migrations:**
   ```bash
   cd backend
   go run cmd/migrate/main.go
   ```

### Twitch Integration Setup

1. **Create Twitch Application:**
   - Go to https://dev.twitch.tv/console
   - Create a new application
   - Set OAuth redirect URL: `http://localhost:3000/auth/callback`
   - Copy Client ID and Client Secret

2. **Configure OAuth Scopes:**
   Required scopes for full functionality:
   - `channel:read:subscriptions` (subscriber count)
   - `user:read:follows` (follower data)
   - `channel:read:analytics` (channel analytics)

## API Usage Examples

### Get Dashboard Overview
```typescript
import { analyticsService } from '@/services/analytics';

const overview = await analyticsService.getDashboardOverview(token);
console.log(`Current followers: ${overview.currentFollowers}`);
```

### Trigger Data Collection
```typescript
const result = await analyticsService.triggerDataCollection(token);
console.log(result.message); // "Data collection triggered successfully"
```

### Get Chart Data
```typescript
const chartData = await analyticsService.getAnalyticsChartData(token, 30);
// Returns follower growth, viewership trends, etc. for last 30 days
```

## System Monitoring

### Job Status Monitoring
```bash
# Check recent analytics jobs
curl -H "Authorization: Bearer $TOKEN" \
     http://localhost:8080/api/analytics/jobs?limit=10
```

### Health Checks
```bash
# Analytics service health
curl http://localhost:8080/api/analytics/health

# Database health  
curl http://localhost:8080/health
```

### System Statistics (Admin)
```bash
# Get system-wide analytics stats
curl -H "X-Admin-Key: $ADMIN_KEY" \
     http://localhost:8080/api/analytics/system
```

## Performance Considerations

### Database Optimization
- Indexes on user_id, date columns for fast queries
- Partitioning by date for large datasets
- Regular VACUUM and ANALYZE operations
- Connection pooling for concurrent requests

### API Rate Limiting
- Twitch API: 800 requests per minute
- Batch processing with delays between requests
- Exponential backoff for failed requests
- Jitter to avoid thundering herd

### Caching Strategy
- Redis for frequently accessed data
- Client-side caching for dashboard metrics
- CDN for static chart images
- Database query result caching

## Security

### Authentication & Authorization
- Clerk-based user authentication
- JWT tokens for API access
- Admin-only endpoints with separate keys
- User isolation in all queries

### Data Protection
- Encrypted database connections
- Secure API key storage
- Input validation and sanitization
- Rate limiting to prevent abuse

## Troubleshooting

### Common Issues

1. **Data Collection Failures:**
   - Check Twitch OAuth token validity
   - Verify API scopes and permissions
   - Review rate limiting constraints

2. **Database Connection Issues:**
   - Verify environment variables
   - Check PostgreSQL service status
   - Review connection pool settings

3. **Frontend API Errors:**
   - Confirm API endpoint URLs
   - Check CORS configuration
   - Verify authentication tokens

### Debug Logs
```bash
# Enable debug logging
export LOG_LEVEL=debug
go run cmd/api/main.go

# Check specific component logs
grep "analytics" /var/log/creatorsync.log
```

## Future Enhancements

### Planned Features
- Multi-platform support (YouTube, TikTok, Instagram)
- Advanced ML-based insights and predictions
- Custom alert system for milestone achievements
- Export functionality for external analysis
- Real-time streaming analytics during live sessions

### Scalability Improvements
- Microservices architecture for better scaling
- Message queue system for background processing
- Horizontal database sharding
- CDN integration for global performance

## Contributing

### Development Workflow
1. Create feature branch from `main`
2. Implement changes with tests
3. Update documentation
4. Submit pull request with description
5. Code review and approval process

### Testing
```bash
# Backend tests
cd backend && go test ./...

# Frontend tests  
cd frontend && npm run test

# Integration tests
make test-integration
```

### Code Standards
- Go: Follow standard Go conventions
- TypeScript: ESLint + Prettier configuration
- Database: Migrations for all schema changes
- Documentation: Update relevant docs with changes

---

For additional support or questions, please refer to the main project documentation or create an issue in the repository. 