package analytics

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/baldybuilds/creatorsync/internal/database"
	"github.com/baldybuilds/creatorsync/internal/twitch"
)

type Repository interface {
	GetDashboardOverview(ctx context.Context, userID string) (*DashboardOverview, error)
	GetAnalyticsChartData(ctx context.Context, userID string, days int) (*AnalyticsChartData, error)
	GetDetailedAnalytics(ctx context.Context, userID string) (*DetailedAnalytics, error)
	GetEnhancedAnalytics(ctx context.Context, userID string, days int) (*EnhancedAnalytics, error)

	SaveVideos(ctx context.Context, conn *database.RequestConnection, userID string, videos []twitch.VideoInfo) (int, error)
	GetVideos(ctx context.Context, conn *database.RequestConnection, userID string, limit int) ([]VideoAnalytics, error)
	GetVideosInDateRange(ctx context.Context, conn *database.RequestConnection, userID string, days int, limit int) ([]VideoAnalytics, error)
	GetVideoCount(ctx context.Context, conn *database.RequestConnection, userID string) (int, error)
}

type repository struct {
	db database.Service
}

func NewRepository(db database.Service) Repository {
	return &repository{db: db}
}

func (r *repository) GetDashboardOverview(ctx context.Context, userID string) (*DashboardOverview, error) {
	overview := &DashboardOverview{}

	db := r.db.GetDB()

	query := `
		SELECT 
			COALESCE(MAX(followers_count), 0) as current_followers,
			COALESCE(MAX(subscriber_count), 0) as current_subscribers,
			COALESCE(MAX(total_views), 0) as total_views
		FROM channel_analytics 
		WHERE user_id = $1 
		AND date >= CURRENT_DATE - INTERVAL '30 days'
	`

	err := db.QueryRowContext(ctx, query, userID).Scan(
		&overview.CurrentFollowers,
		&overview.CurrentSubscribers,
		&overview.TotalViews,
	)

	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to get dashboard overview: %w", err)
	}

	return overview, nil
}

func (r *repository) GetAnalyticsChartData(ctx context.Context, userID string, days int) (*AnalyticsChartData, error) {
	chartData := &AnalyticsChartData{
		FollowerGrowth:   []ChartDataPoint{},
		ViewershipTrends: []ChartDataPoint{},
		StreamFrequency:  []ChartDataPoint{},
		TopGames:         []ChartDataPoint{},
		VideoPerformance: []ChartDataPoint{},
	}

	db := r.db.GetDB()

	query := `
		SELECT date, followers_count 
		FROM channel_analytics 
		WHERE user_id = $1 
		AND date >= CURRENT_DATE - ($2 || ' days')::INTERVAL
		ORDER BY date ASC
	`

	rows, err := db.QueryContext(ctx, query, userID, days)
	if err != nil {
		return chartData, nil
	}
	defer rows.Close()

	for rows.Next() {
		var date time.Time
		var count int

		if err := rows.Scan(&date, &count); err != nil {
			continue
		}

		chartData.FollowerGrowth = append(chartData.FollowerGrowth, ChartDataPoint{
			Date:  date.Format("2006-01-02"),
			Value: float64(count),
		})
	}

	return chartData, nil
}

func (r *repository) GetDetailedAnalytics(ctx context.Context, userID string) (*DetailedAnalytics, error) {
	overview, err := r.GetDashboardOverview(ctx, userID)
	if err != nil {
		return nil, err
	}

	chartData, err := r.GetAnalyticsChartData(ctx, userID, 30)
	if err != nil {
		return nil, err
	}

	return &DetailedAnalytics{
		Overview:       *overview,
		Charts:         *chartData,
		TopStreams:     []StreamSession{},
		TopVideos:      []VideoAnalytics{},
		TopGames:       []GameAnalytics{},
		RecentActivity: []ActivityItem{},
	}, nil
}

func (r *repository) GetEnhancedAnalytics(ctx context.Context, userID string, days int) (*EnhancedAnalytics, error) {
	db := r.db.GetDB()

	var videoCount int
	var totalViews int
	var totalDuration int

	query := `
		SELECT 
			COUNT(*) as video_count,
			COALESCE(SUM(view_count), 0) as total_views,
			COALESCE(SUM(duration_seconds), 0) as total_duration
		FROM video_analytics 
		WHERE user_id = $1
		AND published_at >= CURRENT_DATE - ($2 || ' days')::INTERVAL
	`

	err := db.QueryRowContext(ctx, query, userID, days).Scan(
		&videoCount, &totalViews, &totalDuration,
	)

	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to get enhanced analytics: %w", err)
	}

	avgViews := float64(0)
	if videoCount > 0 {
		avgViews = float64(totalViews) / float64(videoCount)
	}

	analytics := &EnhancedAnalytics{
		Overview: VideoBasedOverview{
			TotalViews:           totalViews,
			VideoCount:           videoCount,
			AverageViewsPerVideo: avgViews,
			TotalWatchTimeHours:  float64(totalDuration) / 3600.0,
		},
		Performance: PerformanceData{
			ViewsOverTime:       []ChartDataPoint{},
			ContentDistribution: []ContentTypeData{},
		},
		TopVideos:    []VideoAnalytics{},
		RecentVideos: []VideoAnalytics{},
	}

	videos, err := r.getVideosForAnalytics(ctx, userID, 20)
	if err == nil {
		if len(videos) > 10 {
			analytics.TopVideos = videos[:10]
			analytics.RecentVideos = videos[:10]
		} else {
			analytics.TopVideos = videos
			analytics.RecentVideos = videos
		}
	}

	return analytics, nil
}

func (r *repository) SaveVideos(ctx context.Context, conn *database.RequestConnection, userID string, videos []twitch.VideoInfo) (int, error) {
	if len(videos) == 0 {
		return 0, nil
	}

	query := `
		INSERT INTO video_analytics (
			user_id, video_id, title, video_type, duration_seconds, view_count,
			thumbnail_url, published_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (video_id) 
		DO UPDATE SET 
			title = EXCLUDED.title,
			view_count = EXCLUDED.view_count,
			updated_at = CURRENT_TIMESTAMP
	`

	savedCount := 0
	var lastError error

	for _, video := range videos {
		// Parse duration from Twitch format (e.g., "1h2m3s") to seconds
		durationSeconds := twitch.ParseDurationToSeconds(video.Duration)

		_, err := conn.Exec(ctx, query,
			userID, video.ID, video.Title, video.Type,
			durationSeconds, video.ViewCount, video.ThumbnailURL, &video.PublishedAt,
		)

		if err != nil {
			log.Printf("⚠️ Failed to save video %s: %v (continuing with others)", video.ID, err)
			lastError = err
			continue // Continue processing other videos instead of failing completely
		}
		savedCount++
	}

	if savedCount > 0 {
		log.Printf("✅ Saved %d/%d videos for user %s", savedCount, len(videos), userID)
	}

	// Only return error if no videos were saved at all
	if savedCount == 0 && lastError != nil {
		return 0, fmt.Errorf("failed to save any videos, last error: %w", lastError)
	}

	return savedCount, nil
}

func (r *repository) GetVideos(ctx context.Context, conn *database.RequestConnection, userID string, limit int) ([]VideoAnalytics, error) {
	query := `
		SELECT video_id, title, video_type, duration_seconds, view_count, 
		       thumbnail_url, published_at, created_at, updated_at
		FROM video_analytics 
		WHERE user_id = $1 
		ORDER BY published_at DESC 
		LIMIT $2
	`

	rows, err := conn.Query(ctx, query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var videos []VideoAnalytics
	for rows.Next() {
		var video VideoAnalytics
		var publishedAt, createdAt, updatedAt time.Time

		err := rows.Scan(
			&video.VideoID, &video.Title, &video.VideoType,
			&video.Duration, &video.ViewCount, &video.ThumbnailURL,
			&publishedAt, &createdAt, &updatedAt,
		)
		if err != nil {
			return nil, err
		}

		video.UserID = userID
		video.PublishedAt = &publishedAt
		video.CreatedAt = createdAt
		video.UpdatedAt = updatedAt
		video.DurationFormatted = twitch.FormatSecondsToHMS(video.Duration)
		videos = append(videos, video)
	}

	return videos, nil
}

func (r *repository) GetVideosInDateRange(ctx context.Context, conn *database.RequestConnection, userID string, days int, limit int) ([]VideoAnalytics, error) {
	query := `
		SELECT video_id, title, video_type, duration_seconds, view_count, 
		       thumbnail_url, published_at, created_at, updated_at
		FROM video_analytics 
		WHERE user_id = $1 
		AND published_at >= CURRENT_DATE - ($2 || ' days')::INTERVAL
		ORDER BY published_at DESC 
		LIMIT $3
	`

	rows, err := conn.Query(ctx, query, userID, days, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var videos []VideoAnalytics
	for rows.Next() {
		var video VideoAnalytics
		var publishedAt, createdAt, updatedAt time.Time

		err := rows.Scan(
			&video.VideoID, &video.Title, &video.VideoType,
			&video.Duration, &video.ViewCount, &video.ThumbnailURL,
			&publishedAt, &createdAt, &updatedAt,
		)
		if err != nil {
			return nil, err
		}

		video.UserID = userID
		video.PublishedAt = &publishedAt
		video.CreatedAt = createdAt
		video.UpdatedAt = updatedAt
		video.DurationFormatted = twitch.FormatSecondsToHMS(video.Duration)
		videos = append(videos, video)
	}

	return videos, nil
}

func (r *repository) GetVideoCount(ctx context.Context, conn *database.RequestConnection, userID string) (int, error) {
	var count int
	row := conn.QueryRow(ctx, "SELECT COUNT(*) FROM video_analytics WHERE user_id = $1", userID)
	err := row.Scan(&count)
	return count, err
}

func (r *repository) getVideosForAnalytics(ctx context.Context, userID string, limit int) ([]VideoAnalytics, error) {
	db := r.db.GetDB()

	query := `
		SELECT video_id, title, video_type, duration_seconds, view_count, 
		       thumbnail_url, published_at, created_at, updated_at
		FROM video_analytics 
		WHERE user_id = $1 
		ORDER BY view_count DESC, published_at DESC
		LIMIT $2
	`

	rows, err := db.QueryContext(ctx, query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var videos []VideoAnalytics
	for rows.Next() {
		var video VideoAnalytics
		var publishedAt, createdAt, updatedAt time.Time

		err := rows.Scan(
			&video.VideoID, &video.Title, &video.VideoType,
			&video.Duration, &video.ViewCount, &video.ThumbnailURL,
			&publishedAt, &createdAt, &updatedAt,
		)
		if err != nil {
			continue
		}

		video.UserID = userID
		video.PublishedAt = &publishedAt
		video.CreatedAt = createdAt
		video.UpdatedAt = updatedAt
		video.DurationFormatted = twitch.FormatSecondsToHMS(video.Duration)
		videos = append(videos, video)
	}

	return videos, nil
}
