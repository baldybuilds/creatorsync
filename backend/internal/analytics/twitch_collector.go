package analytics

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/baldybuilds/creatorsync/internal/database"
	"github.com/baldybuilds/creatorsync/internal/twitch"
)

type TwitchCollector struct {
	tokenHelper *twitch.TwitchTokenHelper
	client      *twitch.Client
	repo        Repository
}

func NewTwitchCollector(db database.Service, repo Repository) (*TwitchCollector, error) {
	tokenHelper, err := twitch.NewTwitchTokenHelper(db)
	if err != nil {
		return nil, fmt.Errorf("failed to create token helper: %w", err)
	}

	client, err := twitch.NewClient(
		os.Getenv("TWITCH_CLIENT_ID"),
		os.Getenv("TWITCH_CLIENT_SECRET"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create Twitch client: %w", err)
	}

	return &TwitchCollector{
		tokenHelper: tokenHelper,
		client:      client,
		repo:        repo,
	}, nil
}

func (tc *TwitchCollector) GetPlatform() Platform {
	return PlatformTwitch
}

func (tc *TwitchCollector) IsConnected(ctx context.Context, userID string) (bool, error) {
	_, err := tc.tokenHelper.GetValidTokenForUser(ctx, userID)
	if err != nil {
		if err.Error() == "twitch not connected" {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (tc *TwitchCollector) ValidateConnection(ctx context.Context, userID string) error {
	token, err := tc.tokenHelper.GetValidTokenForUser(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get valid token: %w", err)
	}

	// Validate token by making a simple API call
	valid, err := tc.client.ValidateToken(ctx, token.AccessToken)
	if err != nil {
		return fmt.Errorf("failed to validate token: %w", err)
	}

	if !valid {
		return fmt.Errorf("token is invalid")
	}

	return nil
}

func (tc *TwitchCollector) CollectChannelMetrics(ctx context.Context, userID string) (*PlatformMetrics, error) {
	token, err := tc.tokenHelper.GetValidTokenForUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get valid token: %w", err)
	}

	twitchUserID, err := tc.tokenHelper.GetTwitchUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get Twitch user ID: %w", err)
	}

	metrics := &PlatformMetrics{
		UserID:      userID,
		Platform:    PlatformTwitch,
		Date:        time.Now().Truncate(24 * time.Hour),
		MetricsData: make(map[string]interface{}),
	}

	// Get user info for total views
	userInfo, err := tc.client.GetUserInfo(token.AccessToken)
	if err == nil {
		metrics.TotalViews = userInfo.ViewCount
		metrics.MetricsData["display_name"] = userInfo.DisplayName
		metrics.MetricsData["login"] = userInfo.Login
		metrics.MetricsData["profile_image_url"] = userInfo.ProfileImageURL
	}

	// Get channel info
	channelInfo, err := tc.client.GetChannelInfo(ctx, token.AccessToken, twitchUserID)
	if err == nil {
		metrics.MetricsData["game_name"] = channelInfo.GameName
		metrics.MetricsData["title"] = channelInfo.Title
		metrics.MetricsData["language"] = channelInfo.Language
	}

	// Get follower count
	followerCount, err := tc.client.GetFollowerCount(token.AccessToken)
	if err != nil {
		metrics.FollowersCount = 0
	} else {
		metrics.FollowersCount = followerCount
	}

	// Get subscriber count (requires broadcaster to be affiliate/partner)
	subscriberCount, err := tc.client.GetSubscriberCount(token.AccessToken)
	if err != nil {
		metrics.SubscriberCount = 0
	} else {
		metrics.SubscriberCount = subscriberCount
	}

	// Get video count
	videos, _, err := tc.client.GetUserVideos(ctx, token.AccessToken, twitchUserID, 1)
	if err == nil && len(videos) > 0 {
		metrics.VideoCount = len(videos)
		metrics.MetricsData["has_videos"] = true
	}

	return metrics, nil
}

func (tc *TwitchCollector) CollectVideoMetrics(ctx context.Context, userID string, limit int) ([]VideoMetrics, error) {
	token, err := tc.tokenHelper.GetValidTokenForUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get valid token: %w", err)
	}

	twitchUserID, err := tc.tokenHelper.GetTwitchUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get Twitch user ID: %w", err)
	}

	// Get videos
	videos, _, err := tc.client.GetUserVideos(ctx, token.AccessToken, twitchUserID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get videos: %w", err)
	}

	var videoMetrics []VideoMetrics
	for _, video := range videos {
		videoData := map[string]interface{}{
			"type":          video.Type,
			"language":      video.Language,
			"thumbnail_url": video.ThumbnailURL,
			"url":           video.URL,
			"description":   video.Description,
			"created_at":    video.CreatedAt,
		}

		metrics := VideoMetrics{
			UserID:      userID,
			Platform:    PlatformTwitch,
			VideoID:     video.ID,
			Title:       video.Title,
			ViewCount:   video.ViewCount,
			Duration:    video.Duration,
			PublishedAt: video.PublishedAt,
			VideoData:   videoData,
		}

		videoMetrics = append(videoMetrics, metrics)
	}

	return videoMetrics, nil
}
