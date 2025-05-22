package handlers

import (
	"fmt"

	"github.com/baldybuilds/creatorsync/internal/server/helpers"
	"github.com/gofiber/fiber/v2"
)

func GetTwitchVideosHandler(c *fiber.Ctx) error {
	twitchContext, err := helpers.GetTwitchRequestContext(c)
	if err != nil {
		return helpers.HandleTwitchError(c, err)
	}
	
	twitchUserID := twitchContext.UserID
	twitchToken := twitchContext.AccessToken
	twitchClient := twitchContext.Client

	// TODO: Consider adding a 'limit' query parameter from the request
	// For now, using the previous default. This could be parsed from c.Query() before calling GetUserVideos.
	limit := 20 // Default limit
	videos, _, err := twitchClient.GetUserVideos(c.Context(), twitchToken, twitchUserID, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to fetch Twitch videos: %v", err),
		})
	}

	return c.JSON(fiber.Map{
		"videos": videos,
	})
}
