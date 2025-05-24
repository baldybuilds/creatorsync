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

	// First, get the Twitch user ID from user info
	userInfo, err := dc.twitchClient.GetUserInfo(twitchToken)
	if err != nil {
		job.ErrorMessage = fmt.Sprintf("Failed to get user info: %v", err)
		return err
	}

	twitchUserID := userInfo.ID

	// Collect videos - use pagination to get more videos
	log.Printf("Fetching videos for user %s (Twitch ID: %s)", userID, twitchUserID)

	var allVideos []twitch.VideoInfo
	limit := 100 // Maximum per request
	totalVideosFetched := 0
	maxVideos := 500 // Reasonable limit to avoid infinite loops

	// Fetch videos with pagination
	cursor := ""
	for totalVideosFetched < maxVideos {
		var videos []twitch.VideoInfo
		var nextCursor string
		var err error

		if cursor == "" {
			videos, nextCursor, err = dc.twitchClient.GetUserVideos(ctx, twitchToken, twitchUserID, limit)
		} else {
			// Note: Currently the GetUserVideos doesn't support cursor parameter
			// This is a limitation we should address later
			videos, nextCursor, err = dc.twitchClient.GetUserVideos(ctx, twitchToken, twitchUserID, limit)
		}

		if err != nil {
			job.ErrorMessage = fmt.Sprintf("Failed to get videos: %v", err)
			log.Printf("Failed to get videos for user %s: %v", userID, err)
			return err
		}

		if len(videos) == 0 {
			break // No more videos
		}

		allVideos = append(allVideos, videos...)
		totalVideosFetched += len(videos)

		log.Printf("Fetched %d videos (total: %d) for user %s", len(videos), totalVideosFetched, userID)

		if nextCursor == "" || len(videos) < limit {
			break // No more pages or last page had fewer than limit
		}
		cursor = nextCursor
	}

	log.Printf("Total videos fetched: %d for user %s", len(allVideos), userID)
	videosSaved := 0

	// Save videos to database
	for _, video := range allVideos {
		durationSeconds := 0
		// TODO: Parse duration string properly (e.g., "1h23m45s" -> seconds)

		videoAnalytics := &VideoAnalytics{
			UserID:       userID,
			VideoID:      video.ID,
			Title:        video.Title,
			VideoType:    video.Type,
			Duration:     durationSeconds,
			ViewCount:    video.ViewCount,
			ThumbnailURL: video.ThumbnailURL,
			PublishedAt:  &video.PublishedAt,
		}

		if err := dc.repo.SaveVideoAnalytics(ctx, videoAnalytics); err != nil {
			log.Printf("Failed to save video analytics for video %s (%s): %v", video.ID, video.Title, err)
		} else {
			videosSaved++
		}
	}

	// Now collect clips
	log.Printf("Fetching clips for user %s (Twitch ID: %s)", userID, twitchUserID)

	clipLimit := 100 // Start with 100 clips
	clips, err := dc.twitchClient.GetClips(ctx, twitchToken, twitchUserID, clipLimit)
	if err != nil {
		log.Printf("Failed to get clips for user %s: %v", userID, err)
		// Don't fail the job for clips, continue
	} else {
		log.Printf("Found %d clips for user %s", len(clips), userID)
		clipsSaved := 0

		for _, clip := range clips {
			// Convert clip to video analytics format
			clipAnalytics := &VideoAnalytics{
				UserID:       userID,
				VideoID:      clip.ID,
				Title:        clip.Title,
				VideoType:    "clip",
				Duration:     int(clip.Duration), // Duration is in seconds for clips
				ViewCount:    clip.ViewCount,
				ThumbnailURL: clip.ThumbnailURL,
				PublishedAt:  &clip.CreatedAt,
			}

			if err := dc.repo.SaveVideoAnalytics(ctx, clipAnalytics); err != nil {
				log.Printf("Failed to save clip analytics for clip %s (%s): %v", clip.ID, clip.Title, err)
			} else {
				clipsSaved++
			}
		}

		log.Printf("Successfully saved %d out of %d clips for user %s", clipsSaved, len(clips), userID)
	}

	log.Printf("Data collection complete for user %s: %d videos, %d clips saved", userID, videosSaved, len(clips))
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
		return fmt.Errorf("failed to get Twitch token for user %s: %w", userID, err)
	}

	// Fetch user info from Twitch
	userInfo, err := dc.twitchClient.GetUserInfo(twitchToken)
	if err != nil {
		return fmt.Errorf("failed to get user info from Twitch for %s: %w", userID, err)
	}

	// Create user record with Twitch info
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

	log.Printf("✅ Created user record for %s (%s)", user.DisplayName, userID)
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
