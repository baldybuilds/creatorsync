package handlers

import (
	"fmt"

	"github.com/baldybuilds/creatorsync/internal/server/helpers"
	"github.com/gofiber/fiber/v2"
)

func GetTwitchClipsHandler(c *fiber.Ctx) error {
	twitchContext, err := helpers.GetTwitchRequestContext(c)
	if err != nil {
		return helpers.HandleTwitchError(c, err)
	}

	twitchUserID := twitchContext.UserID
	twitchToken := twitchContext.AccessToken
	twitchClient := twitchContext.Client

	// TODO: Add query parameters for time range and pagination
	clips, err := twitchClient.GetClips(c.Context(), twitchToken, twitchUserID, 20)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to fetch Twitch clips: %v", err),
		})
	}

	return c.JSON(fiber.Map{
		"clips": clips,
	})
}
