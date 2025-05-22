package analytics

import (
	"time"
)

// User represents a creator user in the system
type User struct {
	ID              string    `json:"id" db:"id"`
	ClerkUserID     string    `json:"clerk_user_id" db:"clerk_user_id"`
	TwitchUserID    string    `json:"twitch_user_id" db:"twitch_user_id"`
	Username        string    `json:"username" db:"username"`
	DisplayName     string    `json:"display_name" db:"display_name"`
	Email           string    `json:"email" db:"email"`
	ProfileImageURL string    `json:"profile_image_url" db:"profile_image_url"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

// ChannelAnalytics represents daily channel metrics
type ChannelAnalytics struct {
	ID              int       `json:"id" db:"id"`
	UserID          string    `json:"user_id" db:"user_id"`
	Date            time.Time `json:"date" db:"date"`
	FollowersCount  int       `json:"followers_count" db:"followers_count"`
	FollowingCount  int       `json:"following_count" db:"following_count"`
	TotalViews      int       `json:"total_views" db:"total_views"`
	SubscriberCount int       `json:"subscriber_count" db:"subscriber_count"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
}

// StreamSession represents individual stream performance
type StreamSession struct {
	ID                int        `json:"id" db:"id"`
	UserID            string     `json:"user_id" db:"user_id"`
	StreamID          string     `json:"stream_id" db:"stream_id"`
	Title             string     `json:"title" db:"title"`
	GameName          string     `json:"game_name" db:"game_name"`
	GameID            string     `json:"game_id" db:"game_id"`
	StartedAt         *time.Time `json:"started_at" db:"started_at"`
	EndedAt           *time.Time `json:"ended_at" db:"ended_at"`
	DurationMinutes   int        `json:"duration_minutes" db:"duration_minutes"`
	PeakViewers       int        `json:"peak_viewers" db:"peak_viewers"`
	AverageViewers    int        `json:"average_viewers" db:"average_viewers"`
	TotalChatters     int        `json:"total_chatters" db:"total_chatters"`
	FollowersGained   int        `json:"followers_gained" db:"followers_gained"`
	SubscribersGained int        `json:"subscribers_gained" db:"subscribers_gained"`
	CreatedAt         time.Time  `json:"created_at" db:"created_at"`
}

// VideoAnalytics represents video performance metrics
type VideoAnalytics struct {
	ID           int        `json:"id" db:"id"`
	UserID       string     `json:"user_id" db:"user_id"`
	VideoID      string     `json:"video_id" db:"video_id"`
	Title        string     `json:"title" db:"title"`
	VideoType    string     `json:"video_type" db:"video_type"`
	Duration     int        `json:"duration_seconds" db:"duration_seconds"`
	ViewCount    int        `json:"view_count" db:"view_count"`
	LikeCount    int        `json:"like_count" db:"like_count"`
	CommentCount int        `json:"comment_count" db:"comment_count"`
	ThumbnailURL string     `json:"thumbnail_url" db:"thumbnail_url"`
	PublishedAt  *time.Time `json:"published_at" db:"published_at"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
}

// VideoDailyStats represents daily video performance tracking
type VideoDailyStats struct {
	ID               int       `json:"id" db:"id"`
	VideoID          string    `json:"video_id" db:"video_id"`
	Date             time.Time `json:"date" db:"date"`
	ViewCount        int       `json:"view_count" db:"view_count"`
	LikeCount        int       `json:"like_count" db:"like_count"`
	CommentCount     int       `json:"comment_count" db:"comment_count"`
	WatchTimeMinutes int       `json:"watch_time_minutes" db:"watch_time_minutes"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
}

// GameAnalytics represents performance by game/category
type GameAnalytics struct {
	ID                   int        `json:"id" db:"id"`
	UserID               string     `json:"user_id" db:"user_id"`
	GameID               string     `json:"game_id" db:"game_id"`
	GameName             string     `json:"game_name" db:"game_name"`
	TotalStreams         int        `json:"total_streams" db:"total_streams"`
	TotalHoursStreamed   float64    `json:"total_hours_streamed" db:"total_hours_streamed"`
	AverageViewers       float64    `json:"average_viewers" db:"average_viewers"`
	PeakViewers          int        `json:"peak_viewers" db:"peak_viewers"`
	TotalFollowersGained int        `json:"total_followers_gained" db:"total_followers_gained"`
	LastStreamedAt       *time.Time `json:"last_streamed_at" db:"last_streamed_at"`
	CreatedAt            time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at" db:"updated_at"`
}

// AnalyticsJob represents data collection job status
type AnalyticsJob struct {
	ID           int        `json:"id" db:"id"`
	UserID       string     `json:"user_id" db:"user_id"`
	JobType      string     `json:"job_type" db:"job_type"`
	Status       string     `json:"status" db:"status"`
	StartedAt    *time.Time `json:"started_at" db:"started_at"`
	CompletedAt  *time.Time `json:"completed_at" db:"completed_at"`
	ErrorMessage string     `json:"error_message" db:"error_message"`
	DataDate     *time.Time `json:"data_date" db:"data_date"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
}

// Dashboard Analytics Response Types

// DashboardOverview provides high-level metrics for the dashboard
type DashboardOverview struct {
	CurrentFollowers      int     `json:"current_followers"`
	FollowerChange        int     `json:"follower_change"`
	FollowerChangePercent float64 `json:"follower_change_percent"`
	CurrentSubscribers    int     `json:"current_subscribers"`
	SubscriberChange      int     `json:"subscriber_change"`
	TotalViews            int     `json:"total_views"`
	ViewChange            int     `json:"view_change"`
	AverageViewers        int     `json:"average_viewers"`
	ViewerChange          int     `json:"viewer_change"`
	StreamsLast30Days     int     `json:"streams_last_30_days"`
	HoursStreamedLast30   float64 `json:"hours_streamed_last_30"`
}

// ChartDataPoint represents a data point for charts
type ChartDataPoint struct {
	Date  string  `json:"date"`
	Value float64 `json:"value"`
	Label string  `json:"label,omitempty"`
}

// AnalyticsChartData contains chart data for the analytics dashboard
type AnalyticsChartData struct {
	FollowerGrowth   []ChartDataPoint `json:"follower_growth"`
	ViewershipTrends []ChartDataPoint `json:"viewership_trends"`
	StreamFrequency  []ChartDataPoint `json:"stream_frequency"`
	TopGames         []ChartDataPoint `json:"top_games"`
	VideoPerformance []ChartDataPoint `json:"video_performance"`
}

// DetailedAnalytics provides comprehensive analytics for the analytics page
type DetailedAnalytics struct {
	Overview       DashboardOverview  `json:"overview"`
	Charts         AnalyticsChartData `json:"charts"`
	TopStreams     []StreamSession    `json:"top_streams"`
	TopVideos      []VideoAnalytics   `json:"top_videos"`
	TopGames       []GameAnalytics    `json:"top_games"`
	RecentActivity []ActivityItem     `json:"recent_activity"`
}

// ActivityItem represents recent activity for the dashboard
type ActivityItem struct {
	Type        string    `json:"type"` // 'stream', 'video', 'milestone'
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Value       string    `json:"value,omitempty"`
	Timestamp   time.Time `json:"timestamp"`
	Icon        string    `json:"icon"`
}

// SystemStats represents system-wide analytics statistics
type SystemStats struct {
	TotalUsers            int       `json:"total_users"`
	ActiveUsers           int       `json:"active_users"`
	TotalJobs             int       `json:"total_jobs"`
	SuccessfulJobs        int       `json:"successful_jobs"`
	FailedJobs            int       `json:"failed_jobs"`
	SuccessRate           float64   `json:"success_rate"`
	AverageCollectionTime string    `json:"average_collection_time"`
	LastCollectionRun     time.Time `json:"last_collection_run"`
}

// Enhanced Analytics Response Types for the new dashboard design

// VideoBasedOverview provides metrics calculated from video analytics
type VideoBasedOverview struct {
	TotalViews           int     `json:"totalViews"`
	VideoCount           int     `json:"videoCount"`
	AverageViewsPerVideo float64 `json:"averageViewsPerVideo"`
	TotalWatchTimeHours  float64 `json:"totalWatchTimeHours"`
	CurrentFollowers     int     `json:"currentFollowers"`
	CurrentSubscribers   int     `json:"currentSubscribers"`
	FollowerChange       int     `json:"followerChange"`
	SubscriberChange     int     `json:"subscriberChange"`
}

// PerformanceData represents performance metrics over time
type PerformanceData struct {
	ViewsOverTime       []ChartDataPoint  `json:"viewsOverTime"`
	ContentDistribution []ContentTypeData `json:"contentDistribution"`
}

// ContentTypeData represents content distribution by type and date
type ContentTypeData struct {
	Date       string `json:"date"`
	Broadcasts int    `json:"broadcasts"`
	Clips      int    `json:"clips"`
	Uploads    int    `json:"uploads"`
}

// EnhancedAnalytics provides comprehensive analytics for the new dashboard design
type EnhancedAnalytics struct {
	Overview     VideoBasedOverview `json:"overview"`
	Performance  PerformanceData    `json:"performance"`
	TopVideos    []VideoAnalytics   `json:"topVideos"`
	RecentVideos []VideoAnalytics   `json:"recentVideos"`
}
