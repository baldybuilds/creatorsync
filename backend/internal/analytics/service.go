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
	GetDashboardOverview(ctx context.Context, userID string) (*DashboardOverview, error)
	GetAnalyticsChartData(ctx context.Context, userID string, days int) (*AnalyticsChartData, error)
	GetDetailedAnalytics(ctx context.Context, userID string) (*DetailedAnalytics, error)
	GetEnhancedAnalytics(ctx context.Context, userID string, days int) (*EnhancedAnalytics, error)

	TriggerDataCollection(ctx context.Context, userID string) error
	RefreshChannelData(ctx context.Context, userID string) error

	GetGrowthAnalysis(ctx context.Context, userID string, period string) (*GrowthAnalysis, error)
	GetContentPerformance(ctx context.Context, userID string) (*ContentPerformance, error)
	GetAnalyticsJobs(ctx context.Context, userID string, limit int) ([]AnalyticsJob, error)
	GetSystemStats(ctx context.Context) (*SystemStats, error)

	CheckUserAnalyticsData(ctx context.Context, userID string) (bool, *time.Time, error)
	CheckTwitchConnection(ctx context.Context, userID string) (bool, error)
	IsHealthy() bool
}

type service struct {
	dbService    database.DatabaseInterface
	standardDB   database.Service
	twitchClient *twitch.Client
	repository   Repository
}

func NewService(dbService database.DatabaseInterface, twitchClient *twitch.Client) Service {
	standardDB := dbService.GetStandardDB()

	return &service{
		dbService:    dbService,
		standardDB:   standardDB,
		twitchClient: twitchClient,
		repository:   NewRepository(standardDB),
	}
}

func (s *service) GetDashboardOverview(ctx context.Context, userID string) (*DashboardOverview, error) {
	overview, err := s.repository.GetDashboardOverview(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get dashboard overview: %w", err)
	}

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

func (s *service) GetAnalyticsChartData(ctx context.Context, userID string, days int) (*AnalyticsChartData, error) {
	chartData, err := s.repository.GetAnalyticsChartData(ctx, userID, days)
	if err != nil {
		return nil, fmt.Errorf("failed to get chart data: %w", err)
	}

	if len(chartData.FollowerGrowth) == 0 {
		chartData = s.generateMockChartData(days)
	}

	return chartData, nil
}

func (s *service) GetDetailedAnalytics(ctx context.Context, userID string) (*DetailedAnalytics, error) {
	analytics, err := s.repository.GetDetailedAnalytics(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get detailed analytics: %w", err)
	}

	analytics.RecentActivity = s.generateRecentActivity(userID)
	return analytics, nil
}

func (s *service) GetEnhancedAnalytics(ctx context.Context, userID string, days int) (*EnhancedAnalytics, error) {
	log.Printf("📊 GetEnhancedAnalytics: Starting for user %s (days: %d)", userID, days)

	conn, err := s.dbService.GetConnection(ctx)
	if err != nil {
		log.Printf("❌ Failed to get database connection, falling back to legacy: %v", err)
		return s.fallbackToLegacy(ctx, userID, days)
	}
	defer conn.Close()

	videoCount, err := s.repository.GetVideoCount(ctx, conn, userID)
	if err != nil {
		log.Printf("❌ Failed to get video count: %v", err)
		return s.fallbackToLegacy(ctx, userID, days)
	}

	if videoCount == 0 {
		log.Printf("⚠️ No video data found for user %s", userID)
		return s.getEmptyAnalytics(), nil
	}

	// Get all videos for overview metrics (total views, etc.)
	allVideos, err := s.repository.GetVideos(ctx, conn, userID, 50)
	if err != nil {
		log.Printf("❌ Failed to get videos: %v", err)
		return s.fallbackToLegacy(ctx, userID, days)
	}

	// Get time-filtered videos for chart data
	filteredVideos, err := s.repository.GetVideosInDateRange(ctx, conn, userID, days, 50)
	if err != nil {
		log.Printf("❌ Failed to get filtered videos: %v", err)
		filteredVideos = allVideos // Fallback to all videos
	}

	// Get channel data for followers/subscribers
	overview, err := s.repository.GetDashboardOverview(ctx, userID)
	if err != nil {
		log.Printf("⚠️ Failed to get channel data, using defaults: %v", err)
		overview = &DashboardOverview{CurrentFollowers: 0, CurrentSubscribers: 0}
	}

	analytics := s.buildAnalyticsFromVideos(allVideos, filteredVideos, videoCount, overview)
	log.Printf("✅ Enhanced analytics completed for user %s", userID)

	return analytics, nil
}

func (s *service) TriggerDataCollection(ctx context.Context, userID string) error {
	log.Printf("🚀 Triggering data collection for user %s", userID)

	conn, err := s.dbService.GetConnection(ctx)
	if err != nil {
		log.Printf("❌ Failed to get database connection: %v", err)
		return fmt.Errorf("database connection failed: %w", err)
	}
	defer conn.Close()

	return s.collectUserData(ctx, conn, userID)
}

func (s *service) RefreshChannelData(ctx context.Context, userID string) error {
	return s.TriggerDataCollection(ctx, userID)
}

func (s *service) GetGrowthAnalysis(ctx context.Context, userID string, period string) (*GrowthAnalysis, error) {
	return &GrowthAnalysis{
		Period:  period,
		Metrics: make(map[string]GrowthMetric),
	}, nil
}

func (s *service) GetContentPerformance(ctx context.Context, userID string) (*ContentPerformance, error) {
	return &ContentPerformance{
		TopVideos: []VideoAnalytics{},
		TopGames:  []GameAnalytics{},
		Insights:  []string{},
	}, nil
}

func (s *service) GetAnalyticsJobs(ctx context.Context, userID string, limit int) ([]AnalyticsJob, error) {
	return []AnalyticsJob{}, nil
}

func (s *service) GetSystemStats(ctx context.Context) (*SystemStats, error) {
	return &SystemStats{
		TotalUsers:            0,
		ActiveUsers:           0,
		TotalJobs:             0,
		SuccessfulJobs:        0,
		FailedJobs:            0,
		SuccessRate:           0,
		AverageCollectionTime: "0s",
		LastCollectionRun:     time.Now(),
	}, nil
}

func (s *service) CheckUserAnalyticsData(ctx context.Context, userID string) (bool, *time.Time, error) {
	conn, err := s.dbService.GetConnection(ctx)
	if err != nil {
		return false, nil, err
	}
	defer conn.Close()

	count, err := s.repository.GetVideoCount(ctx, conn, userID)
	if err != nil {
		return false, nil, err
	}

	hasData := count > 0
	var lastUpdate *time.Time

	if hasData {
		var maxTime time.Time
		row := conn.QueryRow(ctx, "SELECT MAX(updated_at) FROM video_analytics WHERE user_id = $1", userID)
		if err := row.Scan(&maxTime); err == nil {
			lastUpdate = &maxTime
		}
	}

	return hasData, lastUpdate, nil
}

func (s *service) CheckTwitchConnection(ctx context.Context, userID string) (bool, error) {
	// Use the same token helper approach as other parts of the system for consistency
	tokenHelper, err := twitch.NewTwitchTokenHelper(s.standardDB)
	if err != nil {
		log.Printf("❌ Failed to create token helper for connection check: %v", err)
		return false, err
	}

	// Try to get a valid token - this will handle refresh if needed
	_, err = tokenHelper.GetValidTokenForUser(ctx, userID)
	if err != nil {
		if err.Error() == "twitch not connected" {
			return false, nil // Not connected, but no error
		}
		log.Printf("❌ Token validation failed for user %s: %v", userID, err)
		return false, err
	}

	return true, nil
}

func (s *service) IsHealthy() bool {
	return s.dbService.IsHealthy()
}

func (s *service) collectUserData(ctx context.Context, conn *database.RequestConnection, userID string) error {
	// Use a context with timeout to prevent hanging operations
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()

	// Collect both video data and channel data
	err := s.collectVideoData(ctxWithTimeout, conn, userID)
	if err != nil {
		log.Printf("❌ Video data collection failed for user %s: %v", userID, err)
		return err
	}

	err = s.collectChannelData(ctxWithTimeout, conn, userID)
	if err != nil {
		log.Printf("❌ Channel data collection failed for user %s: %v", userID, err)
		return err
	}

	log.Printf("✅ Data collection completed for user %s", userID)
	return nil
}

func (s *service) collectVideoData(ctx context.Context, conn *database.RequestConnection, userID string) error {
	oauthConfig, err := twitch.NewOAuthConfig()
	if err != nil {
		return fmt.Errorf("failed to create OAuth config: %w", err)
	}

	token, err := oauthConfig.GetStoredTokens(ctx, s.standardDB, userID)
	if err != nil {
		return fmt.Errorf("failed to get Twitch token: %w", err)
	}

	userInfo, err := s.twitchClient.GetUserInfo(token.AccessToken)
	if err != nil {
		return fmt.Errorf("failed to get user info: %w", err)
	}

	videos, err := s.fetchAllVideos(ctx, token.AccessToken, userInfo.ID)
	if err != nil {
		return fmt.Errorf("failed to fetch videos: %w", err)
	}

	savedCount, err := s.repository.SaveVideos(ctx, conn, userID, videos)
	if err != nil {
		return fmt.Errorf("failed to save videos: %w", err)
	}

	log.Printf("🎉 Saved %d videos for user %s", savedCount, userID)
	return nil
}

func (s *service) collectChannelData(ctx context.Context, conn *database.RequestConnection, userID string) error {
	log.Printf("📊 Collecting channel data for user %s", userID)

	oauthConfig, err := twitch.NewOAuthConfig()
	if err != nil {
		return fmt.Errorf("failed to create OAuth config: %w", err)
	}

	token, err := oauthConfig.GetStoredTokens(ctx, s.standardDB, userID)
	if err != nil {
		return fmt.Errorf("failed to get Twitch token: %w", err)
	}

	userInfo, err := s.twitchClient.GetUserInfo(token.AccessToken)
	if err != nil {
		return fmt.Errorf("failed to get user info: %w", err)
	}

	followerCount, err := s.twitchClient.GetFollowerCount(token.AccessToken)
	if err != nil {
		log.Printf("⚠️ Failed to get follower count: %v", err)
		followerCount = 0
	}

	subscriberCount, err := s.twitchClient.GetSubscriberCount(token.AccessToken)
	if err != nil {
		log.Printf("⚠️ Failed to get subscriber count: %v", err)
		subscriberCount = 0
	}

	query := `
		INSERT INTO channel_analytics (
			user_id, date, followers_count, following_count, 
			total_views, subscriber_count
		) VALUES ($1, CURRENT_DATE, $2, $3, $4, $5)
		ON CONFLICT (user_id, date) 
		DO UPDATE SET 
			followers_count = EXCLUDED.followers_count,
			following_count = EXCLUDED.following_count,
			total_views = EXCLUDED.total_views,
			subscriber_count = EXCLUDED.subscriber_count
	`

	_, err = conn.Exec(ctx, query, userID, followerCount, 0, userInfo.ViewCount, subscriberCount)
	if err != nil {
		return fmt.Errorf("failed to save channel data: %w", err)
	}

	log.Printf("📈 Saved channel data for user %s (followers: %d, subscribers: %d, views: %d)",
		userID, followerCount, subscriberCount, userInfo.ViewCount)
	return nil
}

func (s *service) fetchAllVideos(ctx context.Context, token, twitchUserID string) ([]twitch.VideoInfo, error) {
	var allVideos []twitch.VideoInfo
	limit := 100
	maxVideos := 500

	// Fetch regular videos (VODs, highlights, uploads)
	for len(allVideos) < maxVideos {
		videos, _, err := s.twitchClient.GetUserVideos(ctx, token, twitchUserID, limit)
		if err != nil {
			return allVideos, err
		}

		if len(videos) == 0 || len(videos) < limit {
			allVideos = append(allVideos, videos...)
			break
		}

		allVideos = append(allVideos, videos...)
	}

	// Fetch clips and convert them to VideoInfo format
	clips, err := s.twitchClient.GetClips(ctx, token, twitchUserID, 50)
	if err != nil {
		log.Printf("⚠️ Failed to fetch clips for user %s: %v (continuing without clips)", twitchUserID, err)
	} else {
		log.Printf("📎 Fetched %d clips for Twitch user %s", len(clips), twitchUserID)

		// Convert clips to VideoInfo format
		for _, clip := range clips {
			clipAsVideo := twitch.VideoInfo{
				ID:           clip.ID,
				UserID:       clip.BroadcasterID,
				UserName:     clip.BroadcasterName,
				Title:        clip.Title,
				Description:  "", // Clips don't have descriptions
				CreatedAt:    clip.CreatedAt,
				PublishedAt:  clip.CreatedAt, // Use created_at as published_at for clips
				URL:          clip.URL,
				ThumbnailURL: clip.ThumbnailURL,
				ViewCount:    clip.ViewCount,
				Language:     clip.Language,
				Type:         "clip",
				Duration:     fmt.Sprintf("%.0fs", clip.Duration), // Convert float seconds to string format
			}
			allVideos = append(allVideos, clipAsVideo)
		}
	}

	log.Printf("📹 Total content fetched: %d videos + clips for user %s", len(allVideos), twitchUserID)
	return allVideos, nil
}

func (s *service) fallbackToLegacy(ctx context.Context, userID string, days int) (*EnhancedAnalytics, error) {
	log.Printf("⚠️ Falling back to legacy analytics for user %s", userID)
	return s.repository.GetEnhancedAnalytics(ctx, userID, days)
}

func (s *service) getEmptyAnalytics() *EnhancedAnalytics {
	return &EnhancedAnalytics{
		Overview: VideoBasedOverview{
			TotalViews:           0,
			VideoCount:           0,
			AverageViewsPerVideo: 0,
			TotalWatchTimeHours:  0,
			CurrentFollowers:     0,
			CurrentSubscribers:   0,
			FollowerChange:       0,
			SubscriberChange:     0,
		},
		Performance: PerformanceData{
			ViewsOverTime:       []ChartDataPoint{},
			ContentDistribution: []ContentTypeData{},
		},
		TopVideos:    []VideoAnalytics{},
		RecentVideos: []VideoAnalytics{},
	}
}

func (s *service) buildAnalyticsFromVideos(allVideos []VideoAnalytics, filteredVideos []VideoAnalytics, videoCount int, overview *DashboardOverview) *EnhancedAnalytics {
	// Calculate overview metrics from ALL videos (all-time data)
	totalViews := 0
	totalWatchTimeSeconds := 0

	for _, video := range allVideos {
		totalViews += video.ViewCount
		totalWatchTimeSeconds += video.Duration
	}

	avgViews := float64(0)
	if len(allVideos) > 0 {
		avgViews = float64(totalViews) / float64(len(allVideos))
	}

	// Generate chart data from FILTERED videos (respects time range)
	viewsOverTime := s.generateViewsOverTime(filteredVideos)
	contentDistribution := s.generateContentDistribution(filteredVideos)

	return &EnhancedAnalytics{
		Overview: VideoBasedOverview{
			TotalViews:           totalViews,
			VideoCount:           videoCount,
			AverageViewsPerVideo: avgViews,
			TotalWatchTimeHours:  float64(totalWatchTimeSeconds) / 3600.0,
			CurrentFollowers:     overview.CurrentFollowers,
			CurrentSubscribers:   overview.CurrentSubscribers,
			FollowerChange:       overview.FollowerChange,
			SubscriberChange:     overview.SubscriberChange,
		},
		Performance: PerformanceData{
			ViewsOverTime:       viewsOverTime,
			ContentDistribution: contentDistribution,
		},
		TopVideos:    s.limitVideos(filteredVideos, 10),
		RecentVideos: s.limitVideos(filteredVideos, 10),
	}
}

func (s *service) generateMockChartData(days int) *AnalyticsChartData {
	chartData := &AnalyticsChartData{}

	baseFollowers := 1000
	for i := days; i >= 0; i-- {
		date := time.Now().AddDate(0, 0, -i).Format("2006-01-02")
		growth := baseFollowers + (days-i)*5 + (i%3)*2
		chartData.FollowerGrowth = append(chartData.FollowerGrowth, ChartDataPoint{
			Date:  date,
			Value: float64(growth),
		})
	}

	for i := days; i >= 0; i-- {
		date := time.Now().AddDate(0, 0, -i).Format("2006-01-02")
		viewers := 50 + (i%7)*20 + (i%3)*10
		chartData.ViewershipTrends = append(chartData.ViewershipTrends, ChartDataPoint{
			Date:  date,
			Value: float64(viewers),
		})
	}

	return chartData
}

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
			Value:       "2.5k views",
			Timestamp:   time.Now().Add(-30 * time.Minute),
			Icon:        "film",
		},
	}
}

func (s *service) limitVideos(videos []VideoAnalytics, limit int) []VideoAnalytics {
	if len(videos) <= limit {
		return videos
	}
	return videos[:limit]
}

func (s *service) generateViewsOverTime(videos []VideoAnalytics) []ChartDataPoint {
	if len(videos) == 0 {
		return []ChartDataPoint{}
	}

	// Group videos by date and sum views
	viewsByDate := make(map[string]int)
	for _, video := range videos {
		if video.PublishedAt != nil {
			dateStr := video.PublishedAt.Format("2006-01-02")
			viewsByDate[dateStr] += video.ViewCount
		}
	}

	// Convert to sorted chart data points
	var chartData []ChartDataPoint
	for date, views := range viewsByDate {
		chartData = append(chartData, ChartDataPoint{
			Date:  date,
			Value: float64(views),
		})
	}

	// Sort by date
	for i := 0; i < len(chartData); i++ {
		for j := i + 1; j < len(chartData); j++ {
			if chartData[i].Date > chartData[j].Date {
				chartData[i], chartData[j] = chartData[j], chartData[i]
			}
		}
	}

	return chartData
}

func (s *service) generateContentDistribution(videos []VideoAnalytics) []ContentTypeData {
	if len(videos) == 0 {
		return []ContentTypeData{}
	}

	// Group content by date and type
	distributionByDate := make(map[string]map[string]int)
	for _, video := range videos {
		if video.PublishedAt != nil {
			dateStr := video.PublishedAt.Format("2006-01-02")
			if distributionByDate[dateStr] == nil {
				distributionByDate[dateStr] = make(map[string]int)
			}

			// Categorize video types
			switch video.VideoType {
			case "archive":
				distributionByDate[dateStr]["broadcasts"]++
			case "highlight", "clip":
				distributionByDate[dateStr]["clips"]++
			case "upload":
				distributionByDate[dateStr]["uploads"]++
			default:
				distributionByDate[dateStr]["uploads"]++ // Default to uploads
			}
		}
	}

	// Convert to sorted chart data
	var chartData []ContentTypeData
	for date, types := range distributionByDate {
		chartData = append(chartData, ContentTypeData{
			Date:       date,
			Broadcasts: types["broadcasts"],
			Clips:      types["clips"],
			Uploads:    types["uploads"],
		})
	}

	// Sort by date
	for i := 0; i < len(chartData); i++ {
		for j := i + 1; j < len(chartData); j++ {
			if chartData[i].Date > chartData[j].Date {
				chartData[i], chartData[j] = chartData[j], chartData[i]
			}
		}
	}

	return chartData
}

type GrowthAnalysis struct {
	Period  string                  `json:"period"`
	Metrics map[string]GrowthMetric `json:"metrics"`
}

type GrowthMetric struct {
	Current       int     `json:"current"`
	Previous      int     `json:"previous"`
	Change        int     `json:"change"`
	PercentChange float64 `json:"percent_change"`
	Trend         string  `json:"trend"`
}

type ContentPerformance struct {
	TopVideos []VideoAnalytics `json:"top_videos"`
	TopGames  []GameAnalytics  `json:"top_games"`
	Insights  []string         `json:"insights"`
}
