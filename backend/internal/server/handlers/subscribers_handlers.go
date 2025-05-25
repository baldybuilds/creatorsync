package handlers

import (
	"fmt"
	"os"

	"github.com/baldybuilds/creatorsync/internal/database"
	"github.com/baldybuilds/creatorsync/internal/twitch"
	"github.com/gofiber/fiber/v2"
)

// GetTwitchSubscribersHandler fetches the list of Twitch subscribers for the broadcaster
func GetTwitchSubscribersHandler(c *fiber.Ctx) error {
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

	// Get basic Twitch client for subscriber operations
	twitchClient, err := twitch.NewClient(os.Getenv("TWITCH_CLIENT_ID"), os.Getenv("TWITCH_CLIENT_SECRET"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to initialize Twitch client",
		})
	}

	// Fetch subscribers
	// TODO: Add support for 'limit' and 'afterCursor' query parameters from the request
	// For now, using default values. These could be parsed from c.Query().
	limit := 20       // Default limit
	afterCursor := "" // Default: no cursor

	subscriptionsResponse, err := twitchClient.GetBroadcasterSubscribers(c.Context(), twitchContext.Token.AccessToken, twitchContext.TwitchUserID, limit, afterCursor)
	if err != nil {
		// Consider if more specific error handling from Twitch API is needed here
		// For example, distinguishing between a 401 (bad token, though helper should catch most) vs 403 (no permission) vs 500
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to fetch Twitch subscribers: %v", err),
		})
	}

	return c.JSON(subscriptionsResponse)
}
