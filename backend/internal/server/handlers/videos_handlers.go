package handlers

import (
	"fmt"
	"strconv"

	"github.com/baldybuilds/creatorsync/internal/database"
	"github.com/baldybuilds/creatorsync/internal/twitch"
	"github.com/gofiber/fiber/v2"
	"github.com/nicklaw5/helix/v2"
)

func GetTwitchVideosHandler(c *fiber.Ctx) error {
	db, ok := c.Locals("db").(database.Service)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database service not available",
		})
	}

	tokenHelper, err := twitch.NewTwitchTokenHelper(db)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to initialize token helper",
		})
	}

	twitchContext, err := tokenHelper.GetTwitchRequestContext(c)
	if err != nil {
		return twitch.HandleTwitchError(c, err)
	}

	videoType := c.Query("type", "archive")
	limitStr := c.Query("limit", "20")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > 100 {
		limit = 20
	}

	videosResp, err := twitchContext.HelixClient.GetVideos(&helix.VideosParams{
		UserID: twitchContext.TwitchUserID,
		Type:   videoType,
		First:  limit,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to fetch Twitch videos: %v", err),
		})
	}

	return c.JSON(fiber.Map{
		"videos": videosResp.Data.Videos,
		"pagination": fiber.Map{
			"cursor": videosResp.Data.Pagination.Cursor,
		},
	})
}
