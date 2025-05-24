package handlers

import (
	"fmt"

	"github.com/baldybuilds/creatorsync/internal/server/helpers"
	"github.com/gofiber/fiber/v2"
)

func GetTwitchChannelHandler(c *fiber.Ctx) error {
	twitchContext, err := helpers.GetTwitchRequestContext(c)
	if err != nil {
		return helpers.HandleTwitchError(c, err)
	}

	twitchUserID := twitchContext.UserID
	twitchToken := twitchContext.AccessToken
	twitchClient := twitchContext.Client

	channelInfo, err := twitchClient.GetChannelInfo(c.Context(), twitchToken, twitchUserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to fetch Twitch channel info: %v", err),
		})
	}

	return c.JSON(fiber.Map{
		"channel": channelInfo,
	})
}
