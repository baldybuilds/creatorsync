package analytics

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/baldybuilds/creatorsync/internal/database"
	"github.com/baldybuilds/creatorsync/internal/twitch"
)

type Collector interface {
	CollectUserData(ctx context.Context, userID string) error
	CollectChannelData(ctx context.Context, userID string) error
	CollectVideoData(ctx context.Context, userID string) error
	ScheduleCollection(ctx context.Context, userID string, interval time.Duration) error
	IsHealthy() bool
}

type collector struct {
	dbService    database.DatabaseInterface
	standardDB   database.Service
	twitchClient *twitch.Client
	repository   Repository
}

func NewCollector(dbService database.DatabaseInterface, twitchClient *twitch.Client) Collector {
	standardDB := dbService.GetStandardDB()

	return &collector{
		dbService:    dbService,
		standardDB:   standardDB,
		twitchClient: twitchClient,
		repository:   NewRepository(standardDB),
	}
}

func (c *collector) CollectUserData(ctx context.Context, userID string) error {
	log.Printf("🚀 Starting data collection for user %s", userID)

	conn, err := c.dbService.GetConnection(ctx)
	if err != nil {
		log.Printf("❌ Failed to get database connection: %v", err)
		return fmt.Errorf("database connection failed: %w", err)
	}
	defer conn.Close()

	err = conn.WithTransaction(ctx, &database.TransactionOptions{
		Timeout: 5 * time.Minute,
	}, func(txConn *database.RequestConnection) error {
		if err := c.collectVideoDataInternal(ctx, txConn, userID); err != nil {
			return fmt.Errorf("video collection failed: %w", err)
		}

		if err := c.collectChannelDataInternal(ctx, txConn, userID); err != nil {
			log.Printf("⚠️ Channel data collection failed (non-fatal): %v", err)
		}

		return nil
	})

	if err != nil {
		log.Printf("❌ Data collection failed for user %s: %v", userID, err)
		return err
	}

	log.Printf("✅ Data collection completed for user %s", userID)
	return nil
}

func (c *collector) CollectChannelData(ctx context.Context, userID string) error {
	conn, err := c.dbService.GetConnection(ctx)
	if err != nil {
		return fmt.Errorf("database connection failed: %w", err)
	}
	defer conn.Close()

	return c.collectChannelDataInternal(ctx, conn, userID)
}

func (c *collector) CollectVideoData(ctx context.Context, userID string) error {
	conn, err := c.dbService.GetConnection(ctx)
	if err != nil {
		return fmt.Errorf("database connection failed: %w", err)
	}
	defer conn.Close()

	return c.collectVideoDataInternal(ctx, conn, userID)
}

func (c *collector) ScheduleCollection(ctx context.Context, userID string, interval time.Duration) error {
	log.Printf("📅 Scheduling collection for user %s every %v", userID, interval)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Printf("⏹️ Collection schedule cancelled for user %s", userID)
			return ctx.Err()
		case <-ticker.C:
			if err := c.CollectUserData(ctx, userID); err != nil {
				log.Printf("❌ Scheduled collection failed for user %s: %v", userID, err)
			}
		}
	}
}

func (c *collector) IsHealthy() bool {
	return c.dbService.IsHealthy()
}

func (c *collector) collectVideoDataInternal(ctx context.Context, conn *database.RequestConnection, userID string) error {
	oauthConfig, err := twitch.NewOAuthConfig()
	if err != nil {
		return fmt.Errorf("failed to create OAuth config: %w", err)
	}

	token, err := oauthConfig.GetStoredTokens(ctx, c.standardDB, userID)
	if err != nil {
		return fmt.Errorf("failed to get Twitch token: %w", err)
	}

	userInfo, err := c.twitchClient.GetUserInfo(token.AccessToken)
	if err != nil {
		return fmt.Errorf("failed to get user info: %w", err)
	}

	videos, err := c.fetchUserVideos(ctx, token.AccessToken, userInfo.ID)
	if err != nil {
		return fmt.Errorf("failed to fetch videos: %w", err)
	}

	savedCount, err := c.repository.SaveVideos(ctx, conn, userID, videos)
	if err != nil {
		return fmt.Errorf("failed to save videos: %w", err)
	}

	log.Printf("🎥 Saved %d videos for user %s", savedCount, userID)
	return nil
}

func (c *collector) collectChannelDataInternal(ctx context.Context, conn *database.RequestConnection, userID string) error {
	log.Printf("📊 Collecting channel data for user %s", userID)

	oauthConfig, err := twitch.NewOAuthConfig()
	if err != nil {
		return fmt.Errorf("failed to create OAuth config: %w", err)
	}

	token, err := oauthConfig.GetStoredTokens(ctx, c.standardDB, userID)
	if err != nil {
		return fmt.Errorf("failed to get Twitch token: %w", err)
	}

	userInfo, err := c.twitchClient.GetUserInfo(token.AccessToken)
	if err != nil {
		return fmt.Errorf("failed to get user info: %w", err)
	}

	followerCount, err := c.twitchClient.GetFollowerCount(token.AccessToken)
	if err != nil {
		log.Printf("⚠️ Failed to get follower count: %v", err)
		followerCount = 0
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
			subscriber_count = EXCLUDED.subscriber_count,
			created_at = CURRENT_TIMESTAMP
	`

	_, err = conn.Exec(ctx, query, userID,
		followerCount, 0, userInfo.ViewCount, 0)
	if err != nil {
		return fmt.Errorf("failed to save channel data: %w", err)
	}

	log.Printf("📈 Saved channel data for user %s (followers: %d, views: %d)", userID, followerCount, userInfo.ViewCount)
	return nil
}

func (c *collector) fetchUserVideos(ctx context.Context, token, twitchUserID string) ([]twitch.VideoInfo, error) {
	var allVideos []twitch.VideoInfo
	limit := 100
	maxVideos := 500

	for len(allVideos) < maxVideos {
		videos, _, err := c.twitchClient.GetUserVideos(ctx, token, twitchUserID, limit)
		if err != nil {
			return allVideos, err
		}

		if len(videos) == 0 || len(videos) < limit {
			allVideos = append(allVideos, videos...)
			break
		}

		allVideos = append(allVideos, videos...)
	}

	log.Printf("🎬 Fetched %d videos for Twitch user %s", len(allVideos), twitchUserID)
	return allVideos, nil
}
