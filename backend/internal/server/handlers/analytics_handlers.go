package handlers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/baldybuilds/creatorsync/internal/server/helpers"
	"github.com/baldybuilds/creatorsync/internal/server/models"
	"github.com/baldybuilds/creatorsync/internal/twitch"
	"github.com/gofiber/fiber/v2"
)

func GetTwitchVideoAnalyticsSummaryHandler(c *fiber.Ctx) error {
	twitchContext, err := helpers.GetTwitchRequestContext(c)
	if err != nil {
		return helpers.HandleTwitchError(c, err)
	}

	twitchUserID := twitchContext.UserID
	twitchToken := twitchContext.AccessToken
	twitchClient := twitchContext.Client

	// Parse query parameters
	periodDaysQuery := c.Query("period_days", "0") // Default to 0 (all time / up to video_limit)
	periodDays, convErr := strconv.Atoi(periodDaysQuery)
	if convErr != nil {
		periodDays = 0 // Default on parse error
	}

	videoLimitQuery := c.Query("video_limit", "20") // Default to 20 videos
	videoLimit, convErr := strconv.Atoi(videoLimitQuery)
	if convErr != nil || videoLimit <= 0 {
		videoLimit = 20 // Default on parse error or invalid value
	}
	if videoLimit > 100 {
		videoLimit = 100 // Max limit for Twitch API GetUserVideos
	}

	// Fetch videos - GetUserVideos fetches most recent 'videoLimit' videos
	fetchedVideos, _, err := twitchClient.GetUserVideos(c.Context(), twitchToken, twitchUserID, videoLimit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to fetch Twitch videos: %v", err)})
	}

	var consideredVideos []twitch.VideoInfo
	now := time.Now()
	var dateRangeStart *time.Time

	if periodDays > 0 {
		cutoffDate := now.AddDate(0, 0, -periodDays)
		for _, v := range fetchedVideos {
			if v.PublishedAt.After(cutoffDate) || v.PublishedAt.Equal(cutoffDate) {
				consideredVideos = append(consideredVideos, v)
				if dateRangeStart == nil || v.PublishedAt.Before(*dateRangeStart) {
					tempDate := v.PublishedAt // Avoid taking address of loop variable
					dateRangeStart = &tempDate
				}
			}
		}
	} else {
		consideredVideos = fetchedVideos
		if len(consideredVideos) > 0 {
			oldestVideo := consideredVideos[len(consideredVideos)-1] // Assuming fetchedVideos is sorted newest to oldest
			dateRangeStart = &oldestVideo.PublishedAt
		}
	}

	analyticsSummary := models.VideoAnalyticsSummary{
		RequestedPeriodDays: periodDays,
		ActualDateRangeEnd:  now, // Or PublishedAt of the newest video if available and preferred
	}

	if dateRangeStart != nil {
		analyticsSummary.ActualDateRangeStart = dateRangeStart
	}

	if len(consideredVideos) == 0 {
		analyticsSummary.ContentDistribution = make(map[string]int)
		return c.JSON(analyticsSummary) // Return early with zeroed/empty analytics
	}

	analyticsSummary.TotalVideosConsidered = len(consideredVideos)
	analyticsSummary.ContentDistribution = make(map[string]int)

	for _, v := range consideredVideos {
		analyticsSummary.TotalViews += v.ViewCount
		analyticsSummary.ContentDistribution[v.Type]++
	}

	if analyticsSummary.TotalVideosConsidered > 0 {
		analyticsSummary.AverageViewsPerVideo = float64(analyticsSummary.TotalViews) / float64(analyticsSummary.TotalVideosConsidered)
	}

	return c.JSON(analyticsSummary)
}
