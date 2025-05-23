package analytics

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
)

type Repository interface {
	// User Management
	CreateOrUpdateUser(ctx context.Context, user *User) error
	GetUserByClerkID(ctx context.Context, clerkUserID string) (*User, error)

	// Channel Analytics
	SaveChannelAnalytics(ctx context.Context, analytics *ChannelAnalytics) error
	GetChannelAnalytics(ctx context.Context, userID string, days int) ([]ChannelAnalytics, error)
	GetLatestChannelAnalytics(ctx context.Context, userID string) (*ChannelAnalytics, error)

	// Stream Sessions
	SaveStreamSession(ctx context.Context, session *StreamSession) error
	GetStreamSessions(ctx context.Context, userID string, limit int) ([]StreamSession, error)
	GetStreamSessionsByDateRange(ctx context.Context, userID string, start, end time.Time) ([]StreamSession, error)

	// Video Analytics
	SaveVideoAnalytics(ctx context.Context, video *VideoAnalytics) error
	GetVideoAnalytics(ctx context.Context, userID string, limit int) ([]VideoAnalytics, error)
	UpdateVideoAnalytics(ctx context.Context, videoID string, views, likes, comments int) error

	// Game Analytics
	SaveGameAnalytics(ctx context.Context, game *GameAnalytics) error
	GetTopGames(ctx context.Context, userID string, limit int) ([]GameAnalytics, error)

	// Dashboard Data
	GetDashboardOverview(ctx context.Context, userID string) (*DashboardOverview, error)
	GetAnalyticsChartData(ctx context.Context, userID string, days int) (*AnalyticsChartData, error)
	GetDetailedAnalytics(ctx context.Context, userID string) (*DetailedAnalytics, error)
	GetEnhancedAnalytics(ctx context.Context, userID string, days int) (*EnhancedAnalytics, error)

	// Jobs
	CreateAnalyticsJob(ctx context.Context, job *AnalyticsJob) error
	UpdateAnalyticsJob(ctx context.Context, jobID int, status string, errorMsg *string) error
	GetAnalyticsJobs(ctx context.Context, userID string, limit int) ([]AnalyticsJob, error)

	// System Stats
	GetSystemStats(ctx context.Context) (*SystemStats, error)

	// Data freshness check
	CheckUserAnalyticsData(ctx context.Context, userID string) (hasData bool, lastUpdate *time.Time, err error)
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{
		db: sqlx.NewDb(db, "postgres"),
	}
}

// User Management Methods

func (r *repository) CreateOrUpdateUser(ctx context.Context, user *User) error {
	query := `
		INSERT INTO users (id, clerk_user_id, twitch_user_id, username, display_name, email, profile_image_url)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (clerk_user_id) 
		DO UPDATE SET 
			twitch_user_id = EXCLUDED.twitch_user_id,
			username = EXCLUDED.username,
			display_name = EXCLUDED.display_name,
			email = EXCLUDED.email,
			profile_image_url = EXCLUDED.profile_image_url,
			updated_at = NOW()
		RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRowContext(ctx, query,
		user.ID, user.ClerkUserID, user.TwitchUserID, user.Username,
		user.DisplayName, user.Email, user.ProfileImageURL).Scan(
		&user.ID, &user.CreatedAt, &user.UpdatedAt)
	return err
}

func (r *repository) GetUserByClerkID(ctx context.Context, clerkUserID string) (*User, error) {
	query := `
		SELECT id, clerk_user_id, twitch_user_id, username, display_name, email, profile_image_url, created_at, updated_at
		FROM users 
		WHERE clerk_user_id = $1
	`

	var user User
	err := r.db.GetContext(ctx, &user, query, clerkUserID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &user, err
}

// Channel Analytics Methods

func (r *repository) SaveChannelAnalytics(ctx context.Context, analytics *ChannelAnalytics) error {
	query := `
		INSERT INTO channel_analytics (user_id, date, followers_count, following_count, total_views, subscriber_count)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (user_id, date) 
		DO UPDATE SET 
			followers_count = EXCLUDED.followers_count,
			following_count = EXCLUDED.following_count,
			total_views = EXCLUDED.total_views,
			subscriber_count = EXCLUDED.subscriber_count
	`
	_, err := r.db.ExecContext(ctx, query,
		analytics.UserID, analytics.Date, analytics.FollowersCount,
		analytics.FollowingCount, analytics.TotalViews, analytics.SubscriberCount)
	return err
}

func (r *repository) GetChannelAnalytics(ctx context.Context, userID string, days int) ([]ChannelAnalytics, error) {
	query := `
		SELECT id, user_id, date, followers_count, following_count, total_views, subscriber_count, created_at
		FROM channel_analytics 
		WHERE user_id = $1 AND date >= CURRENT_DATE - INTERVAL '%d days'
		ORDER BY date DESC
	`

	var analytics []ChannelAnalytics
	err := r.db.SelectContext(ctx, &analytics, fmt.Sprintf(query, days), userID)
	return analytics, err
}

func (r *repository) GetLatestChannelAnalytics(ctx context.Context, userID string) (*ChannelAnalytics, error) {
	query := `
		SELECT id, user_id, date, followers_count, following_count, total_views, subscriber_count, created_at
		FROM channel_analytics 
		WHERE user_id = $1 
		ORDER BY date DESC 
		LIMIT 1
	`

	var analytics ChannelAnalytics
	err := r.db.GetContext(ctx, &analytics, query, userID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &analytics, err
}

// Stream Sessions Methods

func (r *repository) SaveStreamSession(ctx context.Context, session *StreamSession) error {
	query := `
		INSERT INTO stream_sessions (
			user_id, stream_id, title, game_name, game_id, started_at, ended_at,
			duration_minutes, peak_viewers, average_viewers, total_chatters,
			followers_gained, subscribers_gained
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		ON CONFLICT (stream_id) 
		DO UPDATE SET 
			ended_at = EXCLUDED.ended_at,
			duration_minutes = EXCLUDED.duration_minutes,
			peak_viewers = EXCLUDED.peak_viewers,
			average_viewers = EXCLUDED.average_viewers,
			total_chatters = EXCLUDED.total_chatters,
			followers_gained = EXCLUDED.followers_gained,
			subscribers_gained = EXCLUDED.subscribers_gained
	`
	_, err := r.db.ExecContext(ctx, query,
		session.UserID, session.StreamID, session.Title, session.GameName, session.GameID,
		session.StartedAt, session.EndedAt, session.DurationMinutes, session.PeakViewers,
		session.AverageViewers, session.TotalChatters, session.FollowersGained, session.SubscribersGained)
	return err
}

func (r *repository) GetStreamSessions(ctx context.Context, userID string, limit int) ([]StreamSession, error) {
	query := `
		SELECT id, user_id, stream_id, title, game_name, game_id, started_at, ended_at,
			   duration_minutes, peak_viewers, average_viewers, total_chatters,
			   followers_gained, subscribers_gained, created_at
		FROM stream_sessions 
		WHERE user_id = $1 
		ORDER BY started_at DESC 
		LIMIT $2
	`

	var sessions []StreamSession
	err := r.db.SelectContext(ctx, &sessions, query, userID, limit)
	return sessions, err
}

func (r *repository) GetStreamSessionsByDateRange(ctx context.Context, userID string, start, end time.Time) ([]StreamSession, error) {
	query := `
		SELECT id, user_id, stream_id, title, game_name, game_id, started_at, ended_at,
			   duration_minutes, peak_viewers, average_viewers, total_chatters,
			   followers_gained, subscribers_gained, created_at
		FROM stream_sessions 
		WHERE user_id = $1 AND started_at >= $2 AND started_at <= $3
		ORDER BY started_at DESC
	`

	var sessions []StreamSession
	err := r.db.SelectContext(ctx, &sessions, query, userID, start, end)
	return sessions, err
}

// Video Analytics Methods

func (r *repository) SaveVideoAnalytics(ctx context.Context, video *VideoAnalytics) error {
	query := `
		INSERT INTO video_analytics (
			user_id, video_id, title, video_type, duration_seconds, view_count,
			like_count, comment_count, thumbnail_url, published_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (video_id) 
		DO UPDATE SET 
			title = EXCLUDED.title,
			view_count = EXCLUDED.view_count,
			like_count = EXCLUDED.like_count,
			comment_count = EXCLUDED.comment_count,
			updated_at = NOW()
	`
	_, err := r.db.ExecContext(ctx, query,
		video.UserID, video.VideoID, video.Title, video.VideoType, video.Duration,
		video.ViewCount, video.LikeCount, video.CommentCount, video.ThumbnailURL, video.PublishedAt)
	return err
}

func (r *repository) GetVideoAnalytics(ctx context.Context, userID string, limit int) ([]VideoAnalytics, error) {
	query := `
		SELECT id, user_id, video_id, title, video_type, duration_seconds, view_count,
			   like_count, comment_count, thumbnail_url, published_at, created_at, updated_at
		FROM video_analytics 
		WHERE user_id = $1 
		ORDER BY published_at DESC 
		LIMIT $2
	`

	var videos []VideoAnalytics
	err := r.db.SelectContext(ctx, &videos, query, userID, limit)
	return videos, err
}

func (r *repository) UpdateVideoAnalytics(ctx context.Context, videoID string, views, likes, comments int) error {
	query := `
		UPDATE video_analytics 
		SET view_count = $2, like_count = $3, comment_count = $4, updated_at = NOW()
		WHERE video_id = $1
	`
	_, err := r.db.ExecContext(ctx, query, videoID, views, likes, comments)
	return err
}

// Game Analytics Methods

func (r *repository) SaveGameAnalytics(ctx context.Context, game *GameAnalytics) error {
	query := `
		INSERT INTO game_analytics (
			user_id, game_id, game_name, total_streams, total_hours_streamed,
			average_viewers, peak_viewers, total_followers_gained, last_streamed_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (user_id, game_id) 
		DO UPDATE SET 
			game_name = EXCLUDED.game_name,
			total_streams = EXCLUDED.total_streams,
			total_hours_streamed = EXCLUDED.total_hours_streamed,
			average_viewers = EXCLUDED.average_viewers,
			peak_viewers = EXCLUDED.peak_viewers,
			total_followers_gained = EXCLUDED.total_followers_gained,
			last_streamed_at = EXCLUDED.last_streamed_at,
			updated_at = NOW()
	`
	_, err := r.db.ExecContext(ctx, query,
		game.UserID, game.GameID, game.GameName, game.TotalStreams, game.TotalHoursStreamed,
		game.AverageViewers, game.PeakViewers, game.TotalFollowersGained, game.LastStreamedAt)
	return err
}

func (r *repository) GetTopGames(ctx context.Context, userID string, limit int) ([]GameAnalytics, error) {
	query := `
		SELECT id, user_id, game_id, game_name, total_streams, total_hours_streamed,
			   average_viewers, peak_viewers, total_followers_gained, last_streamed_at,
			   created_at, updated_at
		FROM game_analytics 
		WHERE user_id = $1 
		ORDER BY total_hours_streamed DESC 
		LIMIT $2
	`

	var games []GameAnalytics
	err := r.db.SelectContext(ctx, &games, query, userID, limit)
	return games, err
}

// Dashboard Methods

func (r *repository) GetDashboardOverview(ctx context.Context, userID string) (*DashboardOverview, error) {
	query := `
SELECT 
COALESCE(current_analytics.followers_count, 0) as current_followers,
COALESCE(current_analytics.followers_count - previous_analytics.followers_count, 0) as follower_change,
COALESCE(current_analytics.subscriber_count, 0) as current_subscribers,
COALESCE(current_analytics.subscriber_count - previous_analytics.subscriber_count, 0) as subscriber_change,
COALESCE(current_analytics.total_views, 0) as total_views,
COALESCE(current_analytics.total_views - previous_analytics.total_views, 0) as view_change,
COALESCE(stream_stats.average_viewers, 0) as average_viewers,
COALESCE(stream_stats.streams_count, 0) as streams_last_30_days,
COALESCE(stream_stats.total_hours, 0) as hours_streamed_last_30
FROM (
SELECT followers_count, subscriber_count, total_views
FROM channel_analytics 
WHERE user_id = $1 
ORDER BY date DESC 
LIMIT 1
) current_analytics
LEFT JOIN (
SELECT followers_count, subscriber_count, total_views
FROM channel_analytics 
WHERE user_id = $1 
ORDER BY date DESC 
LIMIT 1 OFFSET 7
) previous_analytics ON true
LEFT JOIN (
SELECT 
AVG(average_viewers) as average_viewers,
COUNT(*) as streams_count,
SUM(duration_minutes) / 60.0 as total_hours
FROM stream_sessions 
WHERE user_id = $1 
AND started_at >= CURRENT_DATE - INTERVAL '30 days'
) stream_stats ON true
`

	var overview DashboardOverview
	row := r.db.QueryRowContext(ctx, query, userID)

	var avgViewers sql.NullFloat64
	err := row.Scan(
		&overview.CurrentFollowers, &overview.FollowerChange,
		&overview.CurrentSubscribers, &overview.SubscriberChange,
		&overview.TotalViews, &overview.ViewChange,
		&avgViewers, &overview.StreamsLast30Days, &overview.HoursStreamedLast30,
	)

	if err != nil {
		return nil, err
	}

	overview.AverageViewers = int(avgViewers.Float64)

	// Calculate percentage changes
	if overview.CurrentFollowers > 0 {
		overview.FollowerChangePercent = float64(overview.FollowerChange) / float64(overview.CurrentFollowers-overview.FollowerChange) * 100
	}

	return &overview, nil
}

func (r *repository) GetAnalyticsChartData(ctx context.Context, userID string, days int) (*AnalyticsChartData, error) {
	chartData := &AnalyticsChartData{}

	// Follower growth chart
	followerQuery := `
		SELECT date, followers_count 
		FROM channel_analytics 
		WHERE user_id = $1 AND date >= CURRENT_DATE - INTERVAL '%d days'
		ORDER BY date ASC
	`

	rows, err := r.db.QueryContext(ctx, fmt.Sprintf(followerQuery, days), userID)
	if err != nil {
		return nil, err
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

	// Add more chart data queries here...

	return chartData, nil
}

func (r *repository) GetDetailedAnalytics(ctx context.Context, userID string) (*DetailedAnalytics, error) {
	analytics := &DetailedAnalytics{}

	// Get overview
	overview, err := r.GetDashboardOverview(ctx, userID)
	if err != nil {
		return nil, err
	}
	analytics.Overview = *overview

	// Get chart data
	chartData, err := r.GetAnalyticsChartData(ctx, userID, 30)
	if err != nil {
		return nil, err
	}
	analytics.Charts = *chartData

	// Get top streams
	topStreams, err := r.GetStreamSessions(ctx, userID, 5)
	if err != nil {
		return nil, err
	}
	analytics.TopStreams = topStreams

	// Get top videos
	topVideos, err := r.GetVideoAnalytics(ctx, userID, 5)
	if err != nil {
		return nil, err
	}
	analytics.TopVideos = topVideos

	// Get top games
	topGames, err := r.GetTopGames(ctx, userID, 5)
	if err != nil {
		return nil, err
	}
	analytics.TopGames = topGames

	return analytics, nil
}

// GetEnhancedAnalytics provides video-based analytics for the new dashboard design
func (r *repository) GetEnhancedAnalytics(ctx context.Context, userID string, days int) (*EnhancedAnalytics, error) {
	analytics := &EnhancedAnalytics{}

	// Calculate video-based overview metrics
	overview, err := r.getVideoBasedOverview(ctx, userID, days)
	if err != nil {
		return nil, fmt.Errorf("failed to get video-based overview: %w", err)
	}
	analytics.Overview = *overview

	// Get performance data over time
	performance, err := r.getPerformanceData(ctx, userID, days)
	if err != nil {
		return nil, fmt.Errorf("failed to get performance data: %w", err)
	}
	analytics.Performance = *performance

	// Get top videos by view count
	topVideos, err := r.GetVideoAnalytics(ctx, userID, 5)
	if err != nil {
		return nil, fmt.Errorf("failed to get top videos: %w", err)
	}
	analytics.TopVideos = topVideos

	// Get recent videos
	recentVideos, err := r.GetVideoAnalytics(ctx, userID, 10)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent videos: %w", err)
	}
	analytics.RecentVideos = recentVideos

	return analytics, nil
}

// Helper method to calculate video-based overview metrics
func (r *repository) getVideoBasedOverview(ctx context.Context, userID string, days int) (*VideoBasedOverview, error) {
	// Get video metrics - show ALL videos for the user, not filtered by publish date
	// because we want to show total channel metrics, not just recent videos
	videoQuery := `
		SELECT 
			COALESCE(SUM(view_count), 0) as total_views,
			COUNT(*) as video_count,
			COALESCE(AVG(view_count), 0) as avg_views,
			COALESCE(SUM(duration_seconds), 0) / 3600.0 as total_hours
		FROM video_analytics 
		WHERE user_id = $1
	`

	var totalViews, videoCount int
	var avgViews, totalHours float64
	err := r.db.QueryRowContext(ctx, videoQuery, userID).Scan(
		&totalViews, &videoCount, &avgViews, &totalHours)
	if err != nil {
		log.Printf("❌ Error executing video query for user %s: %v", userID, err)
		return nil, err
	}

	log.Printf("📊 Enhanced analytics query for user %s: found %d videos, %d total views, %.2f avg views, %.2f total hours", userID, videoCount, totalViews, avgViews, totalHours)

	// Get channel metrics (followers, subscribers) from latest channel analytics
	channelQuery := `
		SELECT 
			COALESCE(followers_count, 0) as current_followers,
			COALESCE(subscriber_count, 0) as current_subscribers
		FROM channel_analytics 
		WHERE user_id = $1 
		ORDER BY date DESC 
		LIMIT 1
	`

	var currentFollowers, currentSubscribers int
	err = r.db.QueryRowContext(ctx, channelQuery, userID).Scan(
		&currentFollowers, &currentSubscribers)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	// Calculate follower/subscriber changes (simplified for now)
	followerChange := 0   // TODO: Calculate from previous period
	subscriberChange := 0 // TODO: Calculate from previous period

	return &VideoBasedOverview{
		TotalViews:           totalViews,
		VideoCount:           videoCount,
		AverageViewsPerVideo: avgViews,
		TotalWatchTimeHours:  totalHours,
		CurrentFollowers:     currentFollowers,
		CurrentSubscribers:   currentSubscribers,
		FollowerChange:       followerChange,
		SubscriberChange:     subscriberChange,
	}, nil
}

// Helper method to get performance data over time
func (r *repository) getPerformanceData(ctx context.Context, userID string, days int) (*PerformanceData, error) {
	performance := &PerformanceData{}

	// Views over time (aggregate by day)
	viewsQuery := `
		SELECT 
			DATE(published_at) as date,
			SUM(view_count) as daily_views
		FROM video_analytics 
		WHERE user_id = $1 
		AND published_at >= CURRENT_DATE - INTERVAL '%d days'
		GROUP BY DATE(published_at)
		ORDER BY date ASC
	`

	rows, err := r.db.QueryContext(ctx, fmt.Sprintf(viewsQuery, days), userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var date time.Time
		var views int
		if err := rows.Scan(&date, &views); err != nil {
			continue
		}
		performance.ViewsOverTime = append(performance.ViewsOverTime, ChartDataPoint{
			Date:  date.Format("2006-01-02"),
			Value: float64(views),
		})
	}

	// Content distribution by type and date
	contentQuery := `
		SELECT 
			DATE(published_at) as date,
			video_type,
			COUNT(*) as count
		FROM video_analytics 
		WHERE user_id = $1 
		AND published_at >= CURRENT_DATE - INTERVAL '%d days'
		GROUP BY DATE(published_at), video_type
		ORDER BY date ASC
	`

	rows, err = r.db.QueryContext(ctx, fmt.Sprintf(contentQuery, days), userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Group by date
	contentMap := make(map[string]*ContentTypeData)
	for rows.Next() {
		var date time.Time
		var videoType string
		var count int
		if err := rows.Scan(&date, &videoType, &count); err != nil {
			continue
		}

		dateStr := date.Format("2006-01-02")
		if contentMap[dateStr] == nil {
			contentMap[dateStr] = &ContentTypeData{
				Date: dateStr,
			}
		}

		switch videoType {
		case "archive", "vod":
			contentMap[dateStr].Broadcasts += count
		case "clip":
			contentMap[dateStr].Clips += count
		case "upload":
			contentMap[dateStr].Uploads += count
		}
	}

	// Convert map to slice
	for _, data := range contentMap {
		performance.ContentDistribution = append(performance.ContentDistribution, *data)
	}

	return performance, nil
}

// Analytics Jobs Methods

func (r *repository) CreateAnalyticsJob(ctx context.Context, job *AnalyticsJob) error {
	query := `
		INSERT INTO analytics_jobs (user_id, job_type, status, data_date)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	return r.db.GetContext(ctx, &job.ID, query, job.UserID, job.JobType, job.Status, job.DataDate)
}

func (r *repository) UpdateAnalyticsJob(ctx context.Context, jobID int, status string, errorMsg *string) error {
	query := `
		UPDATE analytics_jobs 
		SET status = $2, completed_at = NOW(), error_message = $3
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, jobID, status, errorMsg)
	return err
}

func (r *repository) GetAnalyticsJobs(ctx context.Context, userID string, limit int) ([]AnalyticsJob, error) {
	query := `
		SELECT id, user_id, job_type, status, started_at, completed_at, 
			   error_message, data_date, created_at
		FROM analytics_jobs 
		WHERE user_id = $1 
		ORDER BY created_at DESC 
		LIMIT $2
	`

	var jobs []AnalyticsJob
	err := r.db.SelectContext(ctx, &jobs, query, userID, limit)
	return jobs, err
}

func (r *repository) GetSystemStats(ctx context.Context) (*SystemStats, error) {
	query := `
		SELECT 
			COUNT(DISTINCT user_id) as total_users,
			COUNT(DISTINCT CASE WHEN created_at >= CURRENT_DATE - INTERVAL '7 days' THEN user_id END) as active_users,
			COUNT(*) as total_jobs,
			COUNT(CASE WHEN status = 'completed' THEN 1 END) as successful_jobs,
			COUNT(CASE WHEN status = 'failed' THEN 1 END) as failed_jobs,
			COALESCE(MAX(created_at), NOW()) as last_collection_run
		FROM analytics_jobs
		WHERE created_at >= CURRENT_DATE - INTERVAL '30 days'
	`

	var stats SystemStats
	row := r.db.QueryRowContext(ctx, query)

	err := row.Scan(
		&stats.TotalUsers, &stats.ActiveUsers, &stats.TotalJobs,
		&stats.SuccessfulJobs, &stats.FailedJobs, &stats.LastCollectionRun,
	)

	if err != nil {
		return nil, err
	}

	if stats.TotalJobs > 0 {
		stats.SuccessRate = float64(stats.SuccessfulJobs) / float64(stats.TotalJobs) * 100
	}

	stats.AverageCollectionTime = "~30s"

	return &stats, nil
}

// CheckUserAnalyticsData checks if a user has analytics data and when it was last updated
func (r *repository) CheckUserAnalyticsData(ctx context.Context, userID string) (bool, *time.Time, error) {
	query := `
		SELECT 
			CASE 
				WHEN COUNT(*) > 0 THEN true 
				ELSE false 
			END as has_data,
			MAX(created_at) as last_update
		FROM (
			SELECT created_at FROM channel_analytics WHERE user_id = $1
			UNION ALL
			SELECT created_at FROM video_analytics WHERE user_id = $1
			UNION ALL
			SELECT created_at FROM stream_sessions WHERE user_id = $1
		) all_data
	`

	var hasData bool
	var lastUpdate sql.NullTime

	err := r.db.QueryRowContext(ctx, query, userID).Scan(&hasData, &lastUpdate)
	if err != nil {
		return false, nil, err
	}

	var lastUpdatePtr *time.Time
	if lastUpdate.Valid {
		lastUpdatePtr = &lastUpdate.Time
	}

	return hasData, lastUpdatePtr, nil
}
