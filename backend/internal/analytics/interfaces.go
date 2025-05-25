package analytics

import (
	"context"
	"time"
)

// Platform represents a social media platform
type Platform string

const (
	PlatformTwitch  Platform = "twitch"
	PlatformYouTube Platform = "youtube"
	PlatformTikTok  Platform = "tiktok"
)

// PlatformMetrics represents metrics that can be collected from any platform
type PlatformMetrics struct {
	UserID          string                 `json:"user_id"`
	Platform        Platform               `json:"platform"`
	Date            time.Time              `json:"date"`
	FollowersCount  int                    `json:"followers_count"`
	TotalViews      int                    `json:"total_views"`
	SubscriberCount int                    `json:"subscriber_count"`
	VideoCount      int                    `json:"video_count"`
	MetricsData     map[string]interface{} `json:"metrics_data"` // Platform-specific data
}

// VideoMetrics represents video-specific metrics from any platform
type VideoMetrics struct {
	UserID      string                 `json:"user_id"`
	Platform    Platform               `json:"platform"`
	VideoID     string                 `json:"video_id"`
	Title       string                 `json:"title"`
	ViewCount   int                    `json:"view_count"`
	Duration    string                 `json:"duration"`
	PublishedAt time.Time              `json:"published_at"`
	VideoData   map[string]interface{} `json:"video_data"` // Platform-specific data
}

// PlatformCollector defines the interface that each platform collector must implement
type PlatformCollector interface {
	// GetPlatform returns the platform this collector handles
	GetPlatform() Platform

	// IsConnected checks if the user has connected their account for this platform
	IsConnected(ctx context.Context, userID string) (bool, error)

	// CollectChannelMetrics collects high-level channel/account metrics
	CollectChannelMetrics(ctx context.Context, userID string) (*PlatformMetrics, error)

	// CollectVideoMetrics collects metrics for individual videos/content
	CollectVideoMetrics(ctx context.Context, userID string, limit int) ([]VideoMetrics, error)

	// ValidateConnection validates that the platform connection is working
	ValidateConnection(ctx context.Context, userID string) error
}

// UniversalAnalyticsCollector orchestrates collection across all platforms
type UniversalAnalyticsCollector interface {
	// RegisterPlatform adds a platform collector
	RegisterPlatform(collector PlatformCollector)

	// GetConnectedPlatforms returns all platforms the user has connected
	GetConnectedPlatforms(ctx context.Context, userID string) ([]Platform, error)

	// CollectUserData collects data from all connected platforms for a user
	CollectUserData(ctx context.Context, userID string) error

	// CollectPlatformData collects data from a specific platform for a user
	CollectPlatformData(ctx context.Context, userID string, platform Platform) error

	// ScheduleCollection schedules regular data collection for a user
	ScheduleCollection(ctx context.Context, userID string, interval time.Duration) error
}

// AnalyticsAggregator combines data from multiple platforms
type AnalyticsAggregator interface {
	// GetCombinedMetrics returns aggregated metrics across all platforms
	GetCombinedMetrics(ctx context.Context, userID string, dateRange DateRange) (*CombinedMetrics, error)

	// GetPlatformComparison compares performance across platforms
	GetPlatformComparison(ctx context.Context, userID string, dateRange DateRange) (*PlatformComparison, error)

	// GetGrowthAnalysis analyzes growth trends across platforms
	GetGrowthAnalysis(ctx context.Context, userID string, dateRange DateRange) (*MultiPlatformGrowthAnalysis, error)
}

// Supporting types for aggregated analytics
type DateRange struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

type CombinedMetrics struct {
	TotalFollowers    int                           `json:"total_followers"`
	TotalViews        int                           `json:"total_views"`
	TotalSubscribers  int                           `json:"total_subscribers"`
	TotalVideos       int                           `json:"total_videos"`
	PlatformBreakdown map[Platform]*PlatformMetrics `json:"platform_breakdown"`
	DateRange         DateRange                     `json:"date_range"`
}

type PlatformComparison struct {
	Platforms    []Platform                    `json:"platforms"`
	Metrics      map[Platform]*PlatformMetrics `json:"metrics"`
	BestPlatform Platform                      `json:"best_platform"`
	Insights     []string                      `json:"insights"`
}

type MultiPlatformGrowthAnalysis struct {
	Platform    Platform  `json:"platform"`
	GrowthRate  float64   `json:"growth_rate"`
	Trend       string    `json:"trend"` // "up", "down", "stable"
	Predictions []string  `json:"predictions"`
	DateRange   DateRange `json:"date_range"`
}
