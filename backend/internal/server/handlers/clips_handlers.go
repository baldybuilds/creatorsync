package handlers

import (
	"fmt"
	"os"

	"github.com/baldybuilds/creatorsync/internal/database"
	"github.com/baldybuilds/creatorsync/internal/twitch"
	"github.com/gofiber/fiber/v2"
)

func GetTwitchClipsHandler(c *fiber.Ctx) error {
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

	// Get basic Twitch client for clip operations
	twitchClient, err := twitch.NewClient(os.Getenv("TWITCH_CLIENT_ID"), os.Getenv("TWITCH_CLIENT_SECRET"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to initialize Twitch client",
		})
	}

	// TODO: Add query parameters for time range and pagination
	clips, err := twitchClient.GetClips(c.Context(), twitchContext.Token.AccessToken, twitchContext.TwitchUserID, 20)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to fetch Twitch clips: %v", err),
		})
	}

	return c.JSON(fiber.Map{
		"clips": clips,
	})
}
