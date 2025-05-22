package analytics

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/baldybuilds/creatorsync/internal/clerk"
	"github.com/baldybuilds/creatorsync/internal/twitch"
)

type DataCollector interface {
	CollectDailyChannelData(ctx context.Context, userID string) error
	CollectStreamData(ctx context.Context, userID string) error
	CollectVideoData(ctx context.Context, userID string) error
	CollectAllUserData(ctx context.Context, userID string) error
}

type dataCollector struct {
	repo         Repository
	twitchClient *twitch.Client
}

func NewDataCollector(repo Repository, twitchClient *twitch.Client) DataCollector {
	return &dataCollector{
		repo:         repo,
		twitchClient: twitchClient,
	}
}

// CollectDailyChannelData collects channel metrics for a given day
func (dc *dataCollector) CollectDailyChannelData(ctx context.Context, userID string) error {
	job := &AnalyticsJob{
		UserID:   userID,
		JobType:  "daily_channel",
		Status:   "running",
		DataDate: &[]time.Time{time.Now()}[0],
	}

	if err := dc.repo.CreateAnalyticsJob(ctx, job); err != nil {
		log.Printf("Failed to create analytics job: %v", err)
	}

	defer func() {
		if job.ID > 0 {
			status := "completed"
			var errorMsg *string
			if job.ErrorMessage != "" {
				status = "failed"
				errorMsg = &job.ErrorMessage
			}
			dc.repo.UpdateAnalyticsJob(ctx, job.ID, status, errorMsg)
		}
	}()

	// Get user's Twitch OAuth token
	twitchToken, err := clerk.GetOAuthToken(ctx, userID, "oauth_twitch")
	if err != nil {
		job.ErrorMessage = fmt.Sprintf("Failed to get Twitch token: %v", err)
		return err
	}

	// Initialize analytics record with defaults
	analytics := &ChannelAnalytics{
		UserID:          userID,
		Date:            time.Now().Truncate(24 * time.Hour),
		FollowersCount:  0,
		FollowingCount:  0,
		TotalViews:      0,
		SubscriberCount: 0,
	}

	// Try to get user info first to get total view count
	log.Printf("Fetching user info for user %s", userID)
	userInfo, err := dc.twitchClient.GetUserInfo(twitchToken)
	if err != nil {
		log.Printf("Failed to get user info: %v", err)
	} else {
		log.Printf("Successfully got user info for %s (ID: %s, Login: %s, ViewCount: %d)",
			userInfo.DisplayName, userInfo.ID, userInfo.Login, userInfo.ViewCount)
		analytics.TotalViews = userInfo.ViewCount
	}

	// Try to get channel info
	log.Printf("Fetching channel info for user %s", userID)
	_, err = dc.twitchClient.GetChannelInfoWithToken(twitchToken)
	if err != nil {
		log.Printf("Failed to get channel info: %v", err)
	} else {
		log.Printf("Successfully got channel info for user %s", userID)
	}

	// Try to get follower count
	log.Printf("Fetching follower count for user %s", userID)
	followers, err := dc.twitchClient.GetFollowerCount(twitchToken)
	if err != nil {
		log.Printf("Failed to get follower count: %v", err)
	} else {
		log.Printf("Successfully got follower count for user %s: %d", userID, followers)
		analytics.FollowersCount = followers
	}

	// Try to get subscriber count
	log.Printf("Fetching subscriber count for user %s", userID)
	subscribers, err := dc.twitchClient.GetSubscriberCount(twitchToken)
	if err != nil {
		log.Printf("Failed to get subscriber count (may be normal for non-partners): %v", err)
	} else {
		log.Printf("Successfully got subscriber count for user %s: %d", userID, subscribers)
		analytics.SubscriberCount = subscribers
	}

	// Save to database (always save what we have, even if some calls failed)
	log.Printf("Saving channel analytics for user %s", userID)
	if err := dc.repo.SaveChannelAnalytics(ctx, analytics); err != nil {
		job.ErrorMessage = fmt.Sprintf("Failed to save channel analytics: %v", err)
		return err
	}

	log.Printf("Successfully collected and saved channel data for user %s (followers: %d, views: %d, subscribers: %d)",
		userID, analytics.FollowersCount, analytics.TotalViews, analytics.SubscriberCount)
	return nil
}

// CollectVideoData collects video analytics (VODs, clips, highlights)
func (dc *dataCollector) CollectVideoData(ctx context.Context, userID string) error {
	job := &AnalyticsJob{
		UserID:  userID,
		JobType: "video_data",
		Status:  "running",
	}

	if err := dc.repo.CreateAnalyticsJob(ctx, job); err != nil {
		log.Printf("Failed to create analytics job: %v", err)
	}

	defer func() {
		if job.ID > 0 {
			status := "completed"
			var errorMsg *string
			if job.ErrorMessage != "" {
				status = "failed"
				errorMsg = &job.ErrorMessage
			}
			dc.repo.UpdateAnalyticsJob(ctx, job.ID, status, errorMsg)
		}
	}()

	// Get user's Twitch OAuth token
	twitchToken, err := clerk.GetOAuthToken(ctx, userID, "oauth_twitch")
	if err != nil {
		job.ErrorMessage = fmt.Sprintf("Failed to get Twitch token: %v", err)
		return err
	}

	// Collect VODs
	log.Printf("Fetching VODs for user %s", userID)
	vods, err := dc.twitchClient.GetVideos(twitchToken, "archive", 50)
	if err != nil {
		log.Printf("Failed to get VODs: %v", err)
	} else {
		log.Printf("Found %d VODs for user %s", len(vods), userID)
		videosSaved := 0
		for _, vod := range vods {
			// Convert duration string to seconds (simplified)
			durationSeconds := 0
			// TODO: Parse duration string properly (e.g., "1h23m45s" -> seconds)

			video := &VideoAnalytics{
				UserID:       userID,
				VideoID:      vod.ID,
				Title:        vod.Title,
				VideoType:    "vod",
				Duration:     durationSeconds,
				ViewCount:    vod.ViewCount,
				ThumbnailURL: vod.ThumbnailURL,
				PublishedAt:  &vod.PublishedAt,
			}

			if err := dc.repo.SaveVideoAnalytics(ctx, video); err != nil {
				log.Printf("Failed to save video analytics for VOD %s (%s): %v", vod.ID, vod.Title, err)
			} else {
				videosSaved++
				log.Printf("Saved video: %s (ID: %s, Views: %d)", vod.Title, vod.ID, vod.ViewCount)
			}
		}
		log.Printf("Successfully saved %d out of %d VODs for user %s", videosSaved, len(vods), userID)
	}

	log.Printf("Successfully completed video data collection for user %s", userID)
	return nil
}

// CollectStreamData collects basic stream data (simplified version)
func (dc *dataCollector) CollectStreamData(ctx context.Context, userID string) error {
	log.Printf("Stream data collection not yet implemented for user %s", userID)
	return nil
}

// ensureUserExists creates or updates a user record before collecting data
func (dc *dataCollector) ensureUserExists(ctx context.Context, userID string) error {
	// Check if user already exists
	existingUser, err := dc.repo.GetUserByClerkID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to check existing user: %w", err)
	}

	if existingUser != nil {
		return nil // User already exists
	}

	// Get user's Twitch OAuth token to fetch profile info
	twitchToken, err := clerk.GetOAuthToken(ctx, userID, "oauth_twitch")
	if err != nil {
		return fmt.Errorf("failed to get Twitch token: %w", err)
	}

	// Fetch user info from Twitch
	userInfo, err := dc.twitchClient.GetUserInfo(twitchToken)
	if err != nil {
		return fmt.Errorf("failed to get user info from Twitch: %w", err)
	}

	// Create user record
	user := &User{
		ID:              userID, // Use Clerk user ID as primary key
		ClerkUserID:     userID,
		TwitchUserID:    userInfo.ID,
		Username:        userInfo.Login,
		DisplayName:     userInfo.DisplayName,
		Email:           userInfo.Email,
		ProfileImageURL: userInfo.ProfileImageURL,
	}

	if err := dc.repo.CreateOrUpdateUser(ctx, user); err != nil {
		return fmt.Errorf("failed to create user record: %w", err)
	}

	log.Printf("Created user record for %s (%s)", user.DisplayName, userID)
	return nil
}

// CollectAllUserData runs all data collection for a user
func (dc *dataCollector) CollectAllUserData(ctx context.Context, userID string) error {
	log.Printf("Starting complete data collection for user %s", userID)

	// Ensure user record exists before collecting analytics
	if err := dc.ensureUserExists(ctx, userID); err != nil {
		log.Printf("Failed to ensure user exists for %s: %v", userID, err)
		return err
	}

	// Collect channel data
	if err := dc.CollectDailyChannelData(ctx, userID); err != nil {
		log.Printf("Channel data collection failed for user %s: %v", userID, err)
	}

	// Collect video data
	if err := dc.CollectVideoData(ctx, userID); err != nil {
		log.Printf("Video data collection failed for user %s: %v", userID, err)
	}

	// Collect stream data
	if err := dc.CollectStreamData(ctx, userID); err != nil {
		log.Printf("Stream data collection failed for user %s: %v", userID, err)
	}

	log.Printf("Completed data collection for user %s", userID)
	return nil
}
