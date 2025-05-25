package handlers

import (
	"fmt"

	"github.com/baldybuilds/creatorsync/internal/database"
	"github.com/baldybuilds/creatorsync/internal/twitch"
	"github.com/gofiber/fiber/v2"
	"github.com/nicklaw5/helix/v2"
)

func GetTwitchChannelHandler(c *fiber.Ctx) error {
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

	channelsResp, err := twitchContext.HelixClient.GetChannelInformation(&helix.GetChannelInformationParams{
		BroadcasterIDs: []string{twitchContext.TwitchUserID},
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to fetch Twitch channel info: %v", err),
		})
	}

	if len(channelsResp.Data.Channels) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Channel not found",
		})
	}

	return c.JSON(fiber.Map{
		"channel": channelsResp.Data.Channels[0],
	})
}
