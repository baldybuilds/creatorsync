package server

import (
	"fmt"
	"log"

	"github.com/baldybuilds/creatorsync/internal/clerk"
	"github.com/baldybuilds/creatorsync/internal/twitch"

	"github.com/gofiber/fiber/v2"
)

func (s *FiberServer) registerTwitchRoutes(api fiber.Router) {
	twitchGroup := api.Group("/twitch")
	twitchGroup.Get("/channel", s.getTwitchChannelHandler)
	twitchGroup.Get("/streams", s.getTwitchStreamsHandler)
	twitchGroup.Get("/videos", s.getTwitchVideosHandler)
	twitchGroup.Get("/callback", s.twitchCallbackHandler)
	twitchGroup.Get("/subscribe", s.twitchSubscribeHandler)
	twitchGroup.Get("/unsubscribe", s.twitchUnsubscribeHandler)
}

func (s *FiberServer) getTwitchChannelHandler(c *fiber.Ctx) error {
	user, err := clerk.GetUserFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
		})
	}

	clerkUser, err := clerk.GetUserByID(c.Context(), user.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to get user profile: %v", err),
		})
	}

	var twitchUserID string
	for _, account := range clerkUser.ExternalAccounts {
		if account.Provider == "oauth_twitch" {
			twitchUserID = account.ProviderUserID
			break
		}
	}

	if twitchUserID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Twitch account not connected",
		})
	}

	twitchToken, err := clerk.GetOAuthToken(c.Context(), user.ID, "oauth_twitch")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to get Twitch token: %v", err),
		})
	}

	twitchClient, err := twitch.NewClient()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to initialize Twitch client: %v", err),
		})
	}

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

func (s *FiberServer) getTwitchStreamsHandler(c *fiber.Ctx) error {
	// TO DO: implement getTwitchStreamsHandler
	return c.JSON(fiber.Map{
		"message": "getTwitchStreamsHandler not implemented",
	})
}

func (s *FiberServer) getTwitchVideosHandler(c *fiber.Ctx) error {
	// TO DO: implement getTwitchVideosHandler
	return c.JSON(fiber.Map{
		"message": "getTwitchVideosHandler not implemented",
	})
}

// twitchCallbackHandler handles OAuth callback from Twitch
func (s *FiberServer) twitchCallbackHandler(c *fiber.Ctx) error {
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

// twitchSubscribeHandler initiates the subscription to Twitch events
func (s *FiberServer) twitchSubscribeHandler(c *fiber.Ctx) error {
	user, err := clerk.GetUserFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
		})
	}

	// Get the user's Twitch token
	twitchToken, err := clerk.GetOAuthToken(c.Context(), user.ID, "oauth_twitch")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to get Twitch token: %v", err),
		})
	}

	// Initialize Twitch client
	twitchClient, err := twitch.NewClient()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to initialize Twitch client: %v", err),
		})
	}

	// Validate the token
	valid, err := twitchClient.ValidateToken(c.Context(), twitchToken)
	if err != nil || !valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid or expired Twitch token",
		})
	}

	// TODO: Implement actual subscription to Twitch events

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Subscribed to Twitch events",
	})
}

// twitchUnsubscribeHandler removes subscriptions to Twitch events
func (s *FiberServer) twitchUnsubscribeHandler(c *fiber.Ctx) error {
	user, err := clerk.GetUserFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
		})
	}

	// TODO: Implement unsubscription from Twitch events

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Unsubscribed from Twitch events",
		"user_id": user.ID,
	})
}
