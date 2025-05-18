# API Design Philosophy

## 🌐 REST-First Approach

The API follows RESTful principles with clear resource-based endpoints. Each endpoint has a single responsibility, making the API predictable and easy to consume from the frontend.

## 🔑 Authentication Strategy

### Clerk JWT Integration
Every protected endpoint expects a Clerk JWT token in the Authorization header. The backend validates these tokens with Clerk's service, extracting user context without maintaining session state. This approach eliminates the complexity of OAuth token management while ensuring security.

### User Context Isolation
Once authenticated, users can only access their own resources. The API automatically filters all queries by the authenticated user ID, preventing data leaks and simplifying authorization logic.

## 📋 Key API Endpoints

### Authentication & User Management
- `GET /api/health` - Health check endpoint
- `POST /api/auth/sync` - Sync user data from Clerk to database
- `GET /api/user/profile` - Get current user profile and settings

### Twitch Content Integration
- `GET /api/twitch/clips` - Get user's Twitch clips with filtering/pagination
- `POST /api/twitch/clips/import` - Import clips from Twitch to local database
- `GET /api/twitch/broadcasts` - Get user's recent streams and VODs
- `GET /api/twitch/stream/live` - Check if user is currently live (no stream key exposure)
- `GET /api/twitch/followers` - Get follower count and recent followers
- `GET /api/twitch/subscribers` - Get subscriber data (for partners/affiliates)
- `GET /api/twitch/analytics` - Get stream analytics and insights

### Clip Management
- `GET /api/clips` - Get user's imported clips from local database
- `GET /api/clips/:id` - Get specific clip details and metadata
- `DELETE /api/clips/:id` - Remove clip from local database

### Video Rendering Pipeline
- `POST /api/render` - Create new render job with editing parameters
- `GET /api/render/:job_id` - Get render job status and progress
- `GET /api/render/jobs` - List user's render jobs with filtering
- `DELETE /api/render/:job_id` - Cancel pending render job

### File Delivery
- `GET /api/download/:job_id` - Download completed render as video file
- `GET /api/download/:job_id/thumbnail` - Get thumbnail of rendered video

### Webhook Handlers
- `POST /api/webhooks/remotion` - Receive render completion from Lambda
- `POST /api/webhooks/twitch` - Handle Twitch EventSub notifications

## 🔄 Async Processing Model

### Job-Based Workflow
Video rendering creates a job record immediately, then processes asynchronously via Remotion Lambda. This prevents request timeouts and allows users to leave the page while rendering continues.

### Status Polling
The frontend polls job status endpoints regularly. This simple approach avoids WebSocket complexity while providing real-time updates on render progress.

### Webhook Integration
Remotion Lambda sends completion webhooks to update job status. This ensures accurate state tracking without requiring complex queue management.

## 🛡️ Security & Performance

### Rate Limiting
API endpoints are rate-limited per user to prevent abuse, especially for resource-intensive operations like clip importing and render job creation.

### Input Validation
All request parameters are validated and sanitized. Video processing parameters have reasonable bounds to prevent resource exhaustion.

### Twitch Token Security
User Twitch access tokens are never stored permanently. They're used for API calls and discarded, with Clerk managing the refresh token lifecycle.

## 🎯 Why This Design

**Comprehensive Twitch Integration**: Full access to streamer data (clips, analytics, live status) provides rich content for the dashboard.

**Secure Token Handling**: Never exposing or storing sensitive data like stream keys while still providing live stream status.

**Scalable Analytics**: Aggregating Twitch analytics provides value-added insights without complex data processing.

**Flexible Content Import**: Supporting both clips and broadcast segments gives creators more content options for repurposing.