package handlers

import (
	"fmt"

	"github.com/baldybuilds/creatorsync/internal/server/helpers"
	"github.com/gofiber/fiber/v2"
)

// GetTwitchSubscribersHandler fetches the list of Twitch subscribers for the broadcaster
func GetTwitchSubscribersHandler(c *fiber.Ctx) error {
	twitchContext, err := helpers.GetTwitchRequestContext(c)
	if err != nil {
		return helpers.HandleTwitchError(c, err)
	}

	twitchUserID := twitchContext.UserID
	twitchToken := twitchContext.AccessToken
	twitchClient := twitchContext.Client

	// The twitchUserID from getTwitchRequestContext is the broadcaster's Twitch ID.

	// Fetch subscribers
	// TODO: Add support for 'limit' and 'afterCursor' query parameters from the request
	// For now, using default values. These could be parsed from c.Query().
	limit := 20       // Default limit
	afterCursor := "" // Default: no cursor

	subscriptionsResponse, err := twitchClient.GetBroadcasterSubscribers(c.Context(), twitchToken, twitchUserID, limit, afterCursor)
	if err != nil {
		// Consider if more specific error handling from Twitch API is needed here
		// For example, distinguishing between a 401 (bad token, though helper should catch most) vs 403 (no permission) vs 500
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to fetch Twitch subscribers: %v", err),
		})
	}

	return c.JSON(subscriptionsResponse)
}
