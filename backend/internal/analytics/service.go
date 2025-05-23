package analytics

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/baldybuilds/creatorsync/internal/database"
	"github.com/baldybuilds/creatorsync/internal/twitch"
)

type Service interface {
	// Data retrieval for dashboard
	GetDashboardOverview(ctx context.Context, userID string) (*DashboardOverview, error)
	GetAnalyticsChartData(ctx context.Context, userID string, days int) (*AnalyticsChartData, error)
	GetDetailedAnalytics(ctx context.Context, userID string) (*DetailedAnalytics, error)
	GetEnhancedAnalytics(ctx context.Context, userID string, days int) (*EnhancedAnalytics, error)

	// Manual data collection triggers
	TriggerDataCollection(ctx context.Context, userID string) error
	RefreshChannelData(ctx context.Context, userID string) error

	// Data analysis
	GetGrowthAnalysis(ctx context.Context, userID string, period string) (*GrowthAnalysis, error)
	GetContentPerformance(ctx context.Context, userID string) (*ContentPerformance, error)

	// Job management
	GetAnalyticsJobs(ctx context.Context, userID string, limit int) ([]AnalyticsJob, error)

	// System stats (admin only)
	GetSystemStats(ctx context.Context) (*SystemStats, error)

	// Data freshness check
	CheckUserAnalyticsData(ctx context.Context, userID string) (hasData bool, lastUpdate *time.Time, err error)
}

type service struct {
	repo      Repository
	collector DataCollector
	db        database.Service
}

func NewService(db database.Service, twitchClient *twitch.Client) Service {
	repo := NewRepository(db.GetDB())
	collector := NewDataCollector(repo, twitchClient)

	return &service{
		repo:      repo,
		collector: collector,
		db:        db,
	}
}

// GetDashboardOverview returns summary metrics for the main dashboard
func (s *service) GetDashboardOverview(ctx context.Context, userID string) (*DashboardOverview, error) {
	overview, err := s.repo.GetDashboardOverview(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get dashboard overview: %w", err)
	}

	// If no data exists, return default overview
	if overview.CurrentFollowers == 0 && overview.TotalViews == 0 {
		return &DashboardOverview{
			CurrentFollowers:   0,
			CurrentSubscribers: 0,
			TotalViews:         0,
			AverageViewers:     0,
		}, nil
	}

	return overview, nil
}

// GetAnalyticsChartData returns chart data for analytics visualization
func (s *service) GetAnalyticsChartData(ctx context.Context, userID string, days int) (*AnalyticsChartData, error) {
	chartData, err := s.repo.GetAnalyticsChartData(ctx, userID, days)
	if err != nil {
		return nil, fmt.Errorf("failed to get chart data: %w", err)
	}

	// Generate mock data if no real data exists yet
	if len(chartData.FollowerGrowth) == 0 {
		chartData = s.generateMockChartData(days)
	}

	return chartData, nil
}

// GetDetailedAnalytics returns comprehensive analytics for the analytics page
func (s *service) GetDetailedAnalytics(ctx context.Context, userID string) (*DetailedAnalytics, error) {
	analytics, err := s.repo.GetDetailedAnalytics(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get detailed analytics: %w", err)
	}

	// Generate recent activity
	analytics.RecentActivity = s.generateRecentActivity(userID)

	return analytics, nil
}

// GetEnhancedAnalytics returns video-based analytics for the new dashboard design
func (s *service) GetEnhancedAnalytics(ctx context.Context, userID string, days int) (*EnhancedAnalytics, error) {
	analytics, err := s.repo.GetEnhancedAnalytics(ctx, userID, days)
	if err != nil {
		return nil, fmt.Errorf("failed to get enhanced analytics: %w", err)
	}

	// If no video data exists, return default structure with zero values
	if analytics.Overview.VideoCount == 0 {
		return &EnhancedAnalytics{
			Overview: VideoBasedOverview{
				TotalViews:           0,
				VideoCount:           0,
				AverageViewsPerVideo: 0,
				TotalWatchTimeHours:  0,
				CurrentFollowers:     analytics.Overview.CurrentFollowers,
				CurrentSubscribers:   analytics.Overview.CurrentSubscribers,
				FollowerChange:       0,
				SubscriberChange:     0,
			},
			Performance: PerformanceData{
				ViewsOverTime:       []ChartDataPoint{},
				ContentDistribution: []ContentTypeData{},
			},
			TopVideos:    []VideoAnalytics{},
			RecentVideos: []VideoAnalytics{},
		}, nil
	}

	return analytics, nil
}

// TriggerDataCollection manually triggers data collection for a user
func (s *service) TriggerDataCollection(ctx context.Context, userID string) error {
	log.Printf("Manually triggering data collection for user %s", userID)

	go func() {
		// Run in background to avoid blocking the API response
		bgCtx := context.Background()
		if err := s.collector.CollectAllUserData(bgCtx, userID); err != nil {
			log.Printf("Background data collection failed for user %s: %v", userID, err)
		}
	}()

	return nil
}

// RefreshChannelData specifically refreshes channel metrics
func (s *service) RefreshChannelData(ctx context.Context, userID string) error {
	return s.collector.CollectDailyChannelData(ctx, userID)
}

// GetGrowthAnalysis provides growth trend analysis
func (s *service) GetGrowthAnalysis(ctx context.Context, userID string, period string) (*GrowthAnalysis, error) {
	// Get historical data based on period
	var days int
	switch period {
	case "week":
		days = 7
	case "month":
		days = 30
	case "quarter":
		days = 90
	case "year":
		days = 365
	default:
		days = 30
	}

	analytics, err := s.repo.GetChannelAnalytics(ctx, userID, days)
	if err != nil {
		return nil, fmt.Errorf("failed to get growth analysis: %w", err)
	}

	// Calculate growth metrics
	growth := &GrowthAnalysis{
		Period:  period,
		Metrics: make(map[string]GrowthMetric),
	}

	if len(analytics) >= 2 {
		latest := analytics[0]
		oldest := analytics[len(analytics)-1]

		// Calculate follower growth
		followerGrowth := latest.FollowersCount - oldest.FollowersCount
		followerPercent := 0.0
		if oldest.FollowersCount > 0 {
			followerPercent = float64(followerGrowth) / float64(oldest.FollowersCount) * 100
		}

		growth.Metrics["followers"] = GrowthMetric{
			Current:       latest.FollowersCount,
			Previous:      oldest.FollowersCount,
			Change:        followerGrowth,
			PercentChange: followerPercent,
			Trend:         getTrend(followerPercent),
		}

		// Calculate view growth
		viewGrowth := latest.TotalViews - oldest.TotalViews
		viewPercent := 0.0
		if oldest.TotalViews > 0 {
			viewPercent = float64(viewGrowth) / float64(oldest.TotalViews) * 100
		}

		growth.Metrics["views"] = GrowthMetric{
			Current:       latest.TotalViews,
			Previous:      oldest.TotalViews,
			Change:        viewGrowth,
			PercentChange: viewPercent,
			Trend:         getTrend(viewPercent),
		}
	}

	return growth, nil
}

// GetContentPerformance analyzes video and stream performance
func (s *service) GetContentPerformance(ctx context.Context, userID string) (*ContentPerformance, error) {
	// Get top videos
	videos, err := s.repo.GetVideoAnalytics(ctx, userID, 10)
	if err != nil {
		return nil, fmt.Errorf("failed to get video analytics: %w", err)
	}

	// Get top games
	games, err := s.repo.GetTopGames(ctx, userID, 5)
	if err != nil {
		return nil, fmt.Errorf("failed to get game analytics: %w", err)
	}

	performance := &ContentPerformance{
		TopVideos: videos,
		TopGames:  games,
		Insights:  s.generateContentInsights(videos, games),
	}

	return performance, nil
}

// GetAnalyticsJobs returns the status of analytics jobs for a user
func (s *service) GetAnalyticsJobs(ctx context.Context, userID string, limit int) ([]AnalyticsJob, error) {
	jobs, err := s.repo.GetAnalyticsJobs(ctx, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get analytics jobs: %w", err)
	}
	return jobs, nil
}

// GetSystemStats returns system-wide analytics statistics
func (s *service) GetSystemStats(ctx context.Context) (*SystemStats, error) {
	stats, err := s.repo.GetSystemStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get system stats: %w", err)
	}
	return stats, nil
}

// CheckUserAnalyticsData checks if a user has analytics data and when it was last updated
func (s *service) CheckUserAnalyticsData(ctx context.Context, userID string) (bool, *time.Time, error) {
	return s.repo.CheckUserAnalyticsData(ctx, userID)
}

// Helper function to generate mock chart data when no real data exists
func (s *service) generateMockChartData(days int) *AnalyticsChartData {
	chartData := &AnalyticsChartData{}

	// Generate mock follower growth data
	baseFollowers := 1000
	for i := days; i >= 0; i-- {
		date := time.Now().AddDate(0, 0, -i).Format("2006-01-02")
		// Simulate growth with some randomness
		growth := baseFollowers + (days-i)*5 + (i%3)*2
		chartData.FollowerGrowth = append(chartData.FollowerGrowth, ChartDataPoint{
			Date:  date,
			Value: float64(growth),
		})
	}

	// Generate mock viewership trends
	for i := days; i >= 0; i-- {
		date := time.Now().AddDate(0, 0, -i).Format("2006-01-02")
		// Simulate viewership with some variance
		viewers := 50 + (i%7)*20 + (i%3)*10
		chartData.ViewershipTrends = append(chartData.ViewershipTrends, ChartDataPoint{
			Date:  date,
			Value: float64(viewers),
		})
	}

	return chartData
}

// Helper function to generate recent activity
func (s *service) generateRecentActivity(userID string) []ActivityItem {
	return []ActivityItem{
		{
			Type:        "stream",
			Title:       "New Stream Session",
			Description: "Just finished a 3-hour gaming session",
			Value:       "156 viewers",
			Timestamp:   time.Now().Add(-2 * time.Hour),
			Icon:        "video",
		},
		{
			Type:        "milestone",
			Title:       "Follower Milestone",
			Description: "Reached 1,200 followers!",
			Timestamp:   time.Now().Add(-1 * time.Hour),
			Icon:        "users",
		},
		{
			Type:        "video",
			Title:       "New Clip Created",
			Description: "Epic win moment got clipped",
			Value:       "89 views",
			Timestamp:   time.Now().Add(-30 * time.Minute),
			Icon:        "video",
		},
	}
}

// Helper function to generate content insights
func (s *service) generateContentInsights(videos []VideoAnalytics, games []GameAnalytics) []string {
	insights := []string{}

	if len(videos) > 0 {
		totalViews := 0
		for _, video := range videos {
			totalViews += video.ViewCount
		}
		avgViews := totalViews / len(videos)
		insights = append(insights, fmt.Sprintf("Your videos average %d views", avgViews))
	}

	if len(games) > 0 {
		topGame := games[0]
		insights = append(insights, fmt.Sprintf("%s is your most streamed game with %.1f hours", topGame.GameName, topGame.TotalHoursStreamed))
	}

	if len(insights) == 0 {
		insights = append(insights, "Start streaming to see performance insights!")
	}

	return insights
}

// Helper function to determine trend direction
func getTrend(percent float64) string {
	if percent > 5 {
		return "up"
	} else if percent < -5 {
		return "down"
	}
	return "stable"
}

// Additional types for analytics responses
type GrowthAnalysis struct {
	Period  string                  `json:"period"`
	Metrics map[string]GrowthMetric `json:"metrics"`
}

type GrowthMetric struct {
	Current       int     `json:"current"`
	Previous      int     `json:"previous"`
	Change        int     `json:"change"`
	PercentChange float64 `json:"percent_change"`
	Trend         string  `json:"trend"` // "up", "down", "stable"
}

type ContentPerformance struct {
	TopVideos []VideoAnalytics `json:"top_videos"`
	TopGames  []GameAnalytics  `json:"top_games"`
	Insights  []string         `json:"insights"`
}
