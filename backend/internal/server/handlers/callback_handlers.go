package handlers

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

// TwitchCallbackHandler handles OAuth callback from Twitch
func TwitchCallbackHandler(c *fiber.Ctx) error {
	code := c.Query("code")
	state := c.Query("state")

	if code == "" {
		log.Printf("Error: No code provided in Twitch callback")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No authorization code provided",
		})
	}

	// Validate state parameter to prevent CSRF attacks
	// TODO: Implement proper state validation

	log.Printf("Received Twitch callback with code: %s and state: %s", code, state)

	// Here you would exchange the code for an access token
	// and associate it with the user's account

	// For now, just return success
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Twitch authentication successful",
	})
}
